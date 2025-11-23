package api

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VersionHandler(c *gin.Context) {
	AuditLogger.WithFields(logrus.Fields{
		"ip": c.ClientIP(),
	}).Info("version")

	data, err := os.ReadFile(filepath.Join(MockPath, "version.json"))
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.Data(500, "application/json", []byte(`{"error": "Internal server error"}`))
		return
	}
	c.Data(200, "application/json", data)
}
