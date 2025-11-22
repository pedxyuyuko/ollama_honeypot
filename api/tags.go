package api

import (
	"github.com/gin-gonic/gin"
)

func TagsHandler(c *gin.Context) {
	models := make([]Model, 0, len(Models))
	for _, model := range Models {
		models = append(models, model)
	}
	c.JSON(200, gin.H{"models": models})
}
