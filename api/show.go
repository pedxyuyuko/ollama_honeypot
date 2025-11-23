package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ShowRequest struct {
	Name string `json:"name"`
}

type ShowResponse struct {
	License    string                 `json:"license"`
	Modelfile  string                 `json:"modelfile"`
	Parameters string                 `json:"parameters"`
	Template   string                 `json:"template"`
	Details    ModelDetails           `json:"details"`
	ModelInfo  map[string]interface{} `json:"model_info"`
	ModifiedAt string                 `json:"modified_at"`
}

func ShowHandler(c *gin.Context) {
	var req ShowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	AuditLogger.WithFields(logrus.Fields{
		"ip": c.ClientIP(),
	}).Info("show")

	model, exists := Models[req.Name]
	if !exists {
		c.JSON(404, gin.H{"error": "Model not found"})
		return
	}

	response := ShowResponse{
		License:    "MIT",
		Parameters: "temperature 0.7\nnum_ctx 2048",
		Template:   "{{ .Prompt }}",
		Details:    model.Details,
		ModelInfo: map[string]interface{}{
			"general.architecture":    model.Details.Architecture,
			"general.file_type":       model.Details.FileType,
			"general.parameter_count": "7B", // fake
			"general.quantization":    model.Details.FileType,
		},
		ModifiedAt: model.ModifiedAt,
	}

	c.JSON(200, response)
}
