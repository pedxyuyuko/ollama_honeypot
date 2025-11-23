package api

import (
	"math/rand"
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

	// Check if prompt is empty
	if strings.TrimSpace(req.Prompt) == "" {
		// Stream empty response
		streamNDJSON(c, func() <-chan map[string]interface{} {
			ch := make(chan map[string]interface{})
			go func() {
				defer close(ch)
				startTime := time.Now()
				loadDuration := int64(rand.Intn(9000000) + 1000000)
				totalDuration := time.Since(startTime).Nanoseconds()
				ch <- map[string]interface{}{
					"model":             fullModel,
					"created_at":        time.Now().Format(time.RFC3339),
					"done":              true,
					"total_duration":    totalDuration,
					"load_duration":     loadDuration,
					"prompt_eval_count": 0,
					"eval_count":        0,
					"eval_duration":     totalDuration - loadDuration,
				}
			}()
			return ch
		})
		return
	}

	// Stream fake response
	streamNDJSON(c, func() <-chan map[string]interface{} {
		ch := make(chan map[string]interface{})
		go func() {
			defer close(ch)
			startTime := time.Now()

			// Fake response text
			var fakeText string
			if len(Responses) == 0 {
				fakeText = "No responses loaded."
			} else {
				idx := rand.Intn(len(Responses))
				resp := Responses[idx]
				repeat := resp.RepeatMin + rand.Intn(resp.RepeatMax-resp.RepeatMin+1)
				fakeText = strings.Repeat(resp.Text, repeat)
			}

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
			loadDuration := int64(rand.Intn(9000000) + 1000000)
			totalDuration := time.Since(startTime).Nanoseconds()
			ch <- map[string]interface{}{
				"model":             fullModel,
				"created_at":        time.Now().Format(time.RFC3339),
				"done":              true,
				"total_duration":    totalDuration,
				"load_duration":     loadDuration, // fake
				"prompt_eval_count": len(req.Prompt),
				"eval_count":        len(fakeText),
				"eval_duration":     totalDuration - loadDuration,
			}
		}()
		return ch
	})
}
