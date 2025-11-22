package api

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

func ChatHandler(c *gin.Context) {
	var req ChatRequest
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

			// Get the last user message for context
			var lastUserMessage string
			for _, msg := range req.Messages {
				if msg.Role == "user" {
					lastUserMessage = msg.Content
				}
			}

			// Fake response text
			fakeText := "This is a simulated chat response from the Ollama honeypot. Your last message was: \"" + lastUserMessage + "\". This conversation is fake and intended for security monitoring purposes."

			// Split into chunks
			chunks := splitIntoChunks(fakeText, 20) // 20 char chunks

			for _, chunk := range chunks {
				ch <- map[string]interface{}{
					"model":      fullModel,
					"created_at": time.Now().Format(time.RFC3339),
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": chunk,
					},
					"done": false,
				}
				time.Sleep(50 * time.Millisecond) // Simulate typing delay
			}

			// Final message
			totalDuration := time.Since(startTime).Nanoseconds()
			ch <- map[string]interface{}{
				"model":                fullModel,
				"created_at":           time.Now().Format(time.RFC3339),
				"done":                 true,
				"total_duration":       totalDuration,
				"load_duration":        1000000, // fake
				"prompt_eval_count":    len(lastUserMessage),
				"prompt_eval_duration": 200000000, // fake
				"eval_count":           len(fakeText),
				"eval_duration":        totalDuration - 1000000 - 200000000,
			}
		}()
		return ch
	})
}
