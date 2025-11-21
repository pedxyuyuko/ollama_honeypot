package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var Serve = &cobra.Command{
	Use:   "serve",
	Short: "Start the honeypot server",
	Long:  `Start the HTTP server to serve as a honeypot for Ollama API.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting honeypot server on :11434")
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Ollama Honeypot")
		})
		log.Fatal(http.ListenAndServe(":11434", nil))
	},
}
