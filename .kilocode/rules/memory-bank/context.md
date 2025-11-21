# Context

## Current Work Focus
The project is in the initial setup phase. The basic structure is in place with a `main.go` entry point and a `serve` command using Cobra and Gin. The immediate focus is on expanding the API endpoints to mimic Ollama's behavior and implementing the logging mechanism.

## Recent Changes
- Initialized project structure.
- Added `main.go` and `cmd/serve.go`.
- Set up `go.mod` with dependencies (`gin`, `cobra`).
- Initialized Memory Bank.
- Added logrus logging framework with structured JSON logging for request details (IP, timestamp, method, path, body).
- Refactored API endpoint handlers into separate files under ./api directory for better organization.

## Next Steps
1.  Implement `POST /api/generate` endpoint to mimic text generation.
2.  Implement `POST /api/chat` endpoint to mimic chat completion.
3.  Implement `GET /api/tags` endpoint to list fake models.
4.  Add configuration support (env vars/flags) for port and logging options.