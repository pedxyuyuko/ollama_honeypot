package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type OpenAIModel struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

type OpenAIModelsResponse struct {
	Object string        `json:"object"`
	Data   []OpenAIModel `json:"data"`
}

func OpenAIModelsHandler(c *gin.Context) {
	AuditLogger.WithFields(logrus.Fields{
		"ip": c.ClientIP(),
	}).Info("openai_models")

	data := make([]OpenAIModel, 0)
	for _, model := range Models {
		openAIModel := OpenAIModel{
			ID:      model.Name,
			Object:  "model",
			Created: 1677610602, // Fixed timestamp for honeypot
			OwnedBy: "openai",
		}
		data = append(data, openAIModel)
	}
	response := OpenAIModelsResponse{
		Object: "list",
		Data:   data,
	}
	c.JSON(200, response)
}
