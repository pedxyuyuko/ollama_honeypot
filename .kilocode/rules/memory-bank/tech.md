# Tech

## Technologies Used
- Go 1.25.4
- Gin web framework
- Cobra CLI library

## Development Setup
- Clone the repository
- Ensure Go 1.25.4 or later is installed
- Run `go mod tidy` to download dependencies
- Build with `go build`
- Run with `./ollama_honeypot serve`

## Technical Constraints
- Must run on port 11434 to mimic Ollama
- Should be lightweight and resource-efficient
- No external database required (logs can be to stdout or files)

## Dependencies
- github.com/gin-gonic/gin v1.11.0
- github.com/spf13/cobra v1.10.1
- Numerous indirect dependencies for Gin and Go standard library enhancements

## Tool Usage Patterns
- VSCode with Go extension for development
- Launch configuration in .vscode/launch.json for debugging the serve command
- Standard Go build and run commands