package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RunningModel struct {
	Name      string       `json:"name"`
	Model     string       `json:"model"`
	Size      int64        `json:"size"`
	SizeVRAM  int64        `json:"size_vram"`
	Context   int          `json:"context_length"`
	Digest    string       `json:"digest"`
	Details   ModelDetails `json:"details"`
	ExpiresAt string       `json:"expires_at"`
}

type PsResponse struct {
	Models []RunningModel `json:"models"`
}

func PsHandler(c *gin.Context) {
	AuditLogger.WithFields(logrus.Fields{
		"ip": c.ClientIP(),
	}).Info("ps")

	running := make([]RunningModel, 0, len(Models))
	for _, model := range Models {
		expiresAt := time.Now().Add(time.Duration(2) * time.Hour).Add(time.Duration(30) * time.Minute).Format(time.RFC3339)
		rm := RunningModel{
			Name:      model.Name,
			Model:     model.Name,
			Context:   64 * 1024,
			Size:      model.Size,
			SizeVRAM:  model.Size,
			Digest:    model.Digest,
			Details:   model.Details,
			ExpiresAt: expiresAt,
		}
		running = append(running, rm)
	}
	c.JSON(200, PsResponse{Models: running})
}
