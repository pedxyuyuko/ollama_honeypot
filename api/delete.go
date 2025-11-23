package api

import (
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func DeleteHandler(c *gin.Context) {
	var req struct {
		Model string `json:"name"`
	}
	// print request body
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	logrus.Infof("Request body: %s", string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&req); err != nil {
		if req.Model == "" {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}
	}

	AuditLogger.WithFields(logrus.Fields{
		"ip":    c.ClientIP(),
		"model": req.Model,
	}).Info("delete")

	c.Status(200)
}
