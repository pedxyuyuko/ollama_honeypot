package api

import (
	"time"

	"github.com/gin-gonic/gin"
)

type RunningModel struct {
	Name      string       `json:"name"`
	Size      int64        `json:"size"`
	SizeVRAM  int64        `json:"size_vram"`
	Digest    string       `json:"digest"`
	Details   ModelDetails `json:"details"`
	ExpiresAt string       `json:"expires_at"`
}

type PsResponse struct {
	Models []RunningModel `json:"models"`
}

func PsHandler(c *gin.Context) {
	running := make([]RunningModel, 0, len(Models))
	for _, model := range Models {
		expiresAt := time.Now().Add(time.Duration(2) * time.Hour).Add(time.Duration(30) * time.Minute).Format(time.RFC3339)
		rm := RunningModel{
			Name:      model.Name,
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
