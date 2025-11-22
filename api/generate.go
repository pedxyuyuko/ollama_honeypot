package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

func splitIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func GenerateHandler(c *gin.Context) {
	var req GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	// Normalize model name
	fullModel := req.Model
	if !strings.Contains(fullModel, ":") {
		fullModel += ":latest"
	}

	// Check if model exists
	if _, exists := Models[fullModel]; !exists {
		c.JSON(404, gin.H{"error": "model not found"})
		return
	}

	// Stream fake response
	streamNDJSON(c, func() <-chan map[string]interface{} {
		ch := make(chan map[string]interface{})
		go func() {
			defer close(ch)
			startTime := time.Now()

			// Fake response text
			fakeText := "This is a simulated response from the Ollama honeypot. Your prompt was: \"" + req.Prompt + "\". This text generation is fake and intended for security monitoring purposes."

			// Split into chunks
			chunks := splitIntoChunks(fakeText, 20) // 20 char chunks

			for _, chunk := range chunks {
				ch <- map[string]interface{}{
					"model":      fullModel,
					"created_at": time.Now().Format(time.RFC3339),
					"response":   chunk,
					"done":       false,
				}
				time.Sleep(50 * time.Millisecond) // Simulate typing delay
			}

			// Final message
			totalDuration := time.Since(startTime).Nanoseconds()
			ch <- map[string]interface{}{
				"model":             fullModel,
				"created_at":        time.Now().Format(time.RFC3339),
				"done":              true,
				"total_duration":    totalDuration,
				"load_duration":     1000000, // fake
				"prompt_eval_count": len(req.Prompt),
				"eval_count":        len(fakeText),
				"eval_duration":     totalDuration - 1000000,
			}
		}()
		return ch
	})
}
