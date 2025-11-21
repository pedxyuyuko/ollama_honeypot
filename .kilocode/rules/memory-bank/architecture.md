# Architecture

## System Architecture
The application is a single binary written in Go. It functions as a CLI tool that starts an HTTP server.

### Components
1.  **CLI Interface (Cobra):** Handles command-line arguments, flags, and subcommands (e.g., `serve`).
2.  **HTTP Server (Gin):** Handles incoming API requests, routing them to appropriate handlers.
3.  **Honeypot Logic:**
    -   **Mock Handlers:** Generate fake responses for Ollama API endpoints.
    -   **Logger:** Intercepts requests and logs details (IP, payload) for analysis.

## Source Code Paths
-   `main.go`: Application entry point. Initializes the root command.
-   `cmd/`: Contains Cobra command definitions.
    -   `cmd/serve.go`: Defines the `serve` command which starts the Gin server and defines routes.

## Key Technical Decisions
-   **Go:** Chosen for performance, single-binary deployment, and strong concurrency support.
-   **Gin:** A high-performance HTTP web framework for Go. Used for routing and handling API requests.
-   **Cobra:** A library for creating powerful modern CLI applications. Used for managing the application's commands and flags.

## Design Patterns
-   **Command Pattern:** Used by Cobra to structure the CLI application.
-   **Middleware Pattern:** Gin middleware will be used for logging and potentially for other cross-cutting concerns like authentication (if added later) or rate limiting.