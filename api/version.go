package api

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func VersionHandler(c *gin.Context) {
	data, err := os.ReadFile(filepath.Join(MockPath, "version.json"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var version map[string]interface{}
	if err := json.Unmarshal(data, &version); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, version)
}
