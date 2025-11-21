package api

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func TagsHandler(c *gin.Context) {
	data, err := os.ReadFile(filepath.Join(MockPath, "tags.json"))
	if err != nil {
		c.Header("Content-Type", "application/json")
		c.Data(500, "application/json", []byte(`{"error": "Internal server error"}`))
		return
	}
	c.Data(200, "application/json", data)
}
