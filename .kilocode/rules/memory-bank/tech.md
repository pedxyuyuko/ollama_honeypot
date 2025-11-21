# Tech

## Technologies Used
-   **Language:** Go (v1.25.4)
-   **Web Framework:** Gin (v1.11.0)
-   **CLI Library:** Cobra (v1.10.1)
-   **Logging Library:** Logrus (v1.9.3)

## Development Setup
-   **Build:** Standard Go build system (`go build`).
-   **Dependency Management:** Go Modules (`go.mod`).

## Technical Constraints
-   **Port:** Defaults to 11434 (Ollama standard), but must be configurable.
-   **Performance:** Must be lightweight to run alongside other services if needed.

## Dependencies
-   `github.com/gin-gonic/gin`: HTTP web framework.
-   `github.com/spf13/cobra`: CLI application structure.
-   `github.com/sirupsen/logrus`: Structured logging library.