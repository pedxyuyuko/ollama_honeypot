package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedxyuyuko/ollama_honeypot/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// requestLogger is a Gin middleware for logging requests
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := logrus.Fields{
			"ip":        c.ClientIP(),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"userAgent": c.GetHeader("User-Agent"),
			"query":     c.Request.URL.RawQuery,
		}

		// Capture request body for methods that typically contain payloads
		if strings.Contains("POST PUT PATCH", c.Request.Method) {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// Restore the body for subsequent handlers
				c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				bodyStr := string(bodyBytes)
				// Truncate large bodies to prevent log bloat
				if len(bodyStr) > 1000 {
					bodyStr = bodyStr[:1000] + "..."
				}
				fields["body"] = bodyStr
			}
		}

		// Log incoming request details using structured logging
		logrus.WithFields(fields).Info("Incoming HTTP request")

		c.Next()
	}
}

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the honeypot server",
	Long:  `Start the HTTP server to serve as a honeypot for Ollama API.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load .env if exists
		_ = godotenv.Load()

		// Determine port
		portStr, _ := cmd.Flags().GetString("port")
		if portStr == "" {
			portStr = os.Getenv("PORT")
		}
		if portStr == "" {
			portStr = "11434"
		}

		// Determine mock path
		mockPath, _ := cmd.Flags().GetString("mock-path")
		if mockPath == "" {
			mockPath = os.Getenv("MOCK_PATH")
		}
		if mockPath == "" {
			mockPath = "./mock"
		}
		api.MockPath = mockPath

		// Determine audit log path
		auditLogPath, _ := cmd.Flags().GetString("log-path")
		if auditLogPath == "" {
			auditLogPath = os.Getenv("LOG_PATH")
		}
		// create path if not exist
		if _, err := os.Stat(auditLogPath); os.IsNotExist(err) {
			if err := os.MkdirAll(auditLogPath, 0755); err != nil {
				log.Printf("Failed to create mock path: %v", err)
			}
		}

		// Initialize audit logger
		api.InitAuditLogger(auditLogPath)

		if err := api.LoadModels(); err != nil {
			log.Printf("Failed to load models: %v", err)
		}

		if err := api.LoadResponses(); err != nil {
			log.Printf("Failed to load responses: %v", err)
		}

		// Set up logrus for text logging to console
		logrus.SetFormatter(&logrus.TextFormatter{})

		if os.Getenv("DEBUG") == "" || os.Getenv("DEBUG") == "0" {
			gin.SetMode(gin.ReleaseMode)
		}

		// Create Gin router without default middleware
		r := gin.New()

		_ = r.SetTrustedProxies(nil)

		// Add custom logging middleware
		r.Use(requestLogger())

		// Add recovery middleware
		r.Use(gin.Recovery())

		r.GET("/", api.HealthHandler)
		r.HEAD("/", api.HealthHandler)

		r.GET("/api/version", api.VersionHandler)

		r.GET("/api/tags", api.TagsHandler)

		r.GET("/api/ps", api.PsHandler)

		r.POST("/api/pull", api.PullHandler)

		r.DELETE("/api/delete", api.DeleteHandler)

		r.POST("/api/show", api.ShowHandler)

		r.POST("/api/generate", api.GenerateHandler)

		r.POST("/api/chat", api.ChatHandler)

		r.GET("/v1/models", api.OpenAIModelsHandler)
		r.POST("/v1/chat/completions", api.OpenAIChatHandler)

		r.GET("/models", api.OpenAIModelsHandler)
		r.POST("/chat/completions", api.OpenAIChatHandler)

		// Handle shutdown signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			log.Println("Shutting down...")
			if err := api.SaveModels(); err != nil {
				log.Printf("Error saving models: %v", err)
			}
			os.Exit(0)
		}()

		fmt.Printf("Starting honeypot server on :%s\n", portStr)
		log.Fatal(r.Run(":" + portStr))
	},
}

func init() {
	Serve.Flags().StringP("port", "p", "", "Port to bind to")
	Serve.Flags().StringP("mock-path", "m", "", "Path to mock data directory")
	Serve.Flags().StringP("log-path", "a", "", "Path to audit log file")
}
