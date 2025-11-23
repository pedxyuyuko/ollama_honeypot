package api

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Role      string   `json:"role"`
	Content   string   `json:"content"`
	ToolCalls []string `json:"tool_calls"`
}

type ChatRequest struct {
	Model    string         `json:"model"`
	Messages []Message      `json:"messages"`
	Stream   bool           `json:"stream"`
	Tools    string         `json:"tools"`
	Options  map[string]any `json:"options"`
	Think    string         `json:"think"`
}

func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": fmt.Sprintf("invalid request: %s", err)})
		return
	}

	// Normalize model name
	fullModel := req.Model
	if !strings.Contains(fullModel, ":") {
		fullModel += ":latest"
	}

	// Filter messages to only user role
	var userMessages []Message
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			userMessages = append(userMessages, msg)
		}
	}

	AuditLogger.WithFields(logrus.Fields{
		"ip":       c.ClientIP(),
		"model":    fullModel,
		"messages": userMessages,
		"options":  req.Options,
		"think":    req.Think,
	}).Info("chat")

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
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": chunk,
					},
					"done": false,
				}
				time.Sleep(50 * time.Millisecond) // Simulate typing delay
			}

			// Final message
			loadDuration := int64(rand.Intn(9000000) + 1000000)
			totalDuration := time.Since(startTime).Nanoseconds()
			ch <- map[string]interface{}{
				"model":                fullModel,
				"created_at":           time.Now().Format(time.RFC3339),
				"done":                 true,
				"total_duration":       totalDuration,
				"load_duration":        loadDuration, // fake
				"prompt_eval_count":    len(lastUserMessage),
				"prompt_eval_duration": 200000000, // fake
				"eval_count":           len(fakeText),
				"eval_duration":        totalDuration - loadDuration - 200000000,
			}
		}()
		return ch
	})
}
