package main

import (
	"log"

	"github.com/pedxyuyuko/ollama_honeypot/v2/cmd"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ollama-honeypot",
	Short: "A Ollama Honeypot",
	Long:  `A honeypot for Ollama API requests.`,
}

func init() {
	rootCmd.AddCommand(cmd.Serve)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
