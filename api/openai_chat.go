package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	Tools       string          `json:"tools"`
	Thinking    string          `json:"reasoning_effort"`
	TopK        float32         `json:"top_k"`
	TopP        float32         `json:"top_p"`
	MinP        float32         `json:"min_p"`
	PP          float32         `json:"presence_penalty"`
	RP          float32         `json:"repeat_penalty"`
	FP          float32         `json:"frequency_penalty"`
	Temperature float32         `json:"temperature"`
	MaxTokens   int             `json:"max_tokens"`
}

type OpenAIDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type OpenAIChoice struct {
	Index        int            `json:"index"`
	Message      *OpenAIMessage `json:"message,omitempty"`
	Delta        *OpenAIDelta   `json:"delta,omitempty"`
	FinishReason *string        `json:"finish_reason"`
}

type OpenAIChatResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []OpenAIChoice `json:"choices"`
	Usage   map[string]int `json:"usage,omitempty"`
}

func stringPtr(s string) *string {
	return &s
}

func getContents(msgs []OpenAIMessage) []string {
	var res []string
	for _, m := range msgs {
		res = append(res, m.Content)
	}
	return res
}

func streamSSE(c *gin.Context, generator func() <-chan OpenAIChatResponse) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Writer.WriteHeader(200)
	ch := generator()
	for resp := range ch {
		data, _ := json.Marshal(resp)
		_, _ = c.Writer.Write([]byte("data: "))
		_, _ = c.Writer.Write(data)
		_, _ = c.Writer.Write([]byte("\n\n"))
		if flusher, ok := c.Writer.(interface{ Flush() }); ok {
			flusher.Flush()
		}
	}
	// Send [DONE]
	_, _ = c.Writer.Write([]byte("data: [DONE]\n\n"))
	if flusher, ok := c.Writer.(interface{ Flush() }); ok {
		flusher.Flush()
	}
}

func OpenAIChatHandler(c *gin.Context) {
	var req OpenAIChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": gin.H{"message": "invalid request", "type": "invalid_request_error"}})
		return
	}

	// Filter messages to only user role
	var userMessages []OpenAIMessage
	for _, msg := range req.Messages {
		if msg.Role == "user" {
			userMessages = append(userMessages, msg)
		}
	}

	AuditLogger.WithFields(logrus.Fields{
		"ip":                c.ClientIP(),
		"model":             req.Model,
		"messages":          userMessages,
		"reasoning_effort":  req.Thinking,
		"top_k":             req.TopK,
		"temperature":       req.Temperature,
		"max_tokens":        req.MaxTokens,
		"top_p":             req.TopP,
		"frequency_penalty": req.FP,
		"presence_penalty":  req.PP,
	}).Info("openai_chat")

	// Generate fake response text
	var fakeText string
	if len(Responses) == 0 {
		fakeText = "No responses loaded."
	} else {
		idx := rand.Intn(len(Responses))
		resp := Responses[idx]
		repeat := resp.RepeatMin + rand.Intn(resp.RepeatMax-resp.RepeatMin+1)
		fakeText = strings.Repeat(resp.Text, repeat)
	}

	if req.Stream {
		streamSSE(c, func() <-chan OpenAIChatResponse {
			ch := make(chan OpenAIChatResponse)
			go func() {
				defer close(ch)
				id := fmt.Sprintf("chatcmpl-%d", rand.Int63())
				created := time.Now().Unix()

				// First chunk: role
				ch <- OpenAIChatResponse{
					ID:      id,
					Object:  "chat.completion.chunk",
					Created: created,
					Model:   req.Model,
					Choices: []OpenAIChoice{{
						Index: 0,
						Delta: &OpenAIDelta{Role: "assistant"},
					}},
				}

				// Content chunks
				chunks := splitIntoChunks(fakeText, 20) // 20 char chunks
				for _, chunk := range chunks {
					ch <- OpenAIChatResponse{
						ID:      id,
						Object:  "chat.completion.chunk",
						Created: created,
						Model:   req.Model,
						Choices: []OpenAIChoice{{
							Index: 0,
							Delta: &OpenAIDelta{Content: chunk},
						}},
					}
					time.Sleep(50 * time.Millisecond) // Simulate typing delay
				}

				// Final chunk
				finish := "stop"
				ch <- OpenAIChatResponse{
					ID:      id,
					Object:  "chat.completion.chunk",
					Created: created,
					Model:   req.Model,
					Choices: []OpenAIChoice{{
						Index:        0,
						Delta:        &OpenAIDelta{},
						FinishReason: &finish,
					}},
				}
			}()
			return ch
		})
	} else {
		// Non-streaming response
		id := fmt.Sprintf("chatcmpl-%d", rand.Int63())
		promptTokens := len(strings.Join(getContents(userMessages), " "))
		completionTokens := len(fakeText)
		response := OpenAIChatResponse{
			ID:      id,
			Object:  "chat.completion",
			Created: time.Now().Unix(),
			Model:   req.Model,
			Choices: []OpenAIChoice{{
				Index: 0,
				Message: &OpenAIMessage{
					Role:    "assistant",
					Content: fakeText,
				},
				FinishReason: stringPtr("stop"),
			}},
			Usage: map[string]int{
				"prompt_tokens":     promptTokens,
				"completion_tokens": completionTokens,
				"total_tokens":      promptTokens + completionTokens,
			},
		}
		c.JSON(200, response)
	}
}
