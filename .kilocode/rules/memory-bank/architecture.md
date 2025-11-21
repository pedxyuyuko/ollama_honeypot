# Architecture

## System Architecture
The system is a simple HTTP server built with the Gin framework that listens on port 11434. It uses Cobra for CLI commands to provide a command-line interface.

## Source Code Paths
- main.go: Entry point that sets up the CLI using Cobra
- cmd/serve.go: Implementation of the "serve" command that starts the HTTP server
- go.mod: Defines the Go module and dependencies

## Key Technical Decisions
- Go language chosen for its performance, simplicity, and strong concurrency support
- Gin web framework selected for HTTP handling due to its high performance and ease of use
- Cobra CLI library used for command-line interface due to its popularity and robustness in the Go ecosystem
- Port 11434 chosen to match Ollama's default port for realistic honeypot behavior

## Design Patterns
- Command pattern implemented via Cobra for CLI structure
- Handler pattern used for HTTP route definitions in Gin

## Component Relationships
- main.go imports and uses the cmd package
- cmd/serve.go imports and uses Gin for creating and running the HTTP server
- The serve command in cmd/serve.go sets up routes and starts the server

## Critical Implementation Paths
- Application startup: main.go -> cmd.Serve execution
- Server initialization: cmd/serve.go Run function creates Gin router and starts listening
- Request handling: Currently only handles GET requests to "/" with a simple response