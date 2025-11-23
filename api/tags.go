package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ModelDetailsOutput struct {
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

type ModelInfo struct {
	Name       string             `json:"name"`
	Model      string             `json:"model"`
	ModifiedAt string             `json:"modified_at"`
	Size       int64              `json:"size"`
	Digest     string             `json:"digest"`
	Details    ModelDetailsOutput `json:"details"`
}

type TagsResponse struct {
	Models []ModelInfo `json:"models"`
}

func TagsHandler(c *gin.Context) {
	AuditLogger.WithFields(logrus.Fields{
		"ip": c.ClientIP(),
	}).Info("tags")

	models := make([]ModelInfo, 0)
	for _, model := range Models {
		modelInfo := ModelInfo{
			Name:       model.Name,
			Model:      model.Name,
			ModifiedAt: model.ModifiedAt,
			Size:       model.Size,
			Digest:     model.Digest,
			Details: ModelDetailsOutput{
				Format:            model.Details.ModelFormat,
				Family:            model.Details.ModelFamily,
				Families:          model.Details.ModelFamilies,
				ParameterSize:     model.Details.ModelType,
				QuantizationLevel: model.Details.FileType,
			},
		}
		models = append(models, modelInfo)
	}
	response := TagsResponse{Models: models}
	c.JSON(200, response)
}
