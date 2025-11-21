package cmd

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the honeypot server",
	Long:  `Start the HTTP server to serve as a honeypot for Ollama API.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/", func(c *gin.Context) {
			c.String(200, "Ollama Honeypot")
		})
		fmt.Println("Starting honeypot server on :11434")
		log.Fatal(r.Run(":11434"))
	},
}
