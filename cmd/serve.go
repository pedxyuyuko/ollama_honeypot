package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/pedxyuyuko/ollama_honeypot/v2/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// requestLogger is a Gin middleware for logging requests
func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read request body for logging
		var body string
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			if c.Request.Body != nil {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err == nil {
					body = string(bodyBytes)
					// Restore the body for further processing
					c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				}
			}
		}

		// Log using logrus
		logrus.WithFields(logrus.Fields{
			"timestamp": time.Now().Format(time.RFC3339),
			"ip":        c.ClientIP(),
			"method":    c.Request.Method,
			"path":      c.Request.URL.Path,
			"body":      body,
		}).Info("Request received")

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

		if err := api.LoadModels(); err != nil {
			log.Printf("Failed to load models: %v", err)
		}

		// Set up logrus for JSON structured logging
		logrus.SetFormatter(&logrus.JSONFormatter{})

		// Create Gin router without default middleware
		r := gin.New()

		// Add custom logging middleware
		r.Use(requestLogger())

		// Add recovery middleware
		r.Use(gin.Recovery())

		r.GET("/", api.HealthHandler)
		r.HEAD("/", api.HealthHandler)

		r.GET("/api/version", api.VersionHandler)

		r.GET("/api/tags", api.TagsHandler)

		r.POST("/api/pull", api.PullHandler)

		r.DELETE("/api/delete", api.DeleteHandler)

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
}
