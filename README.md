# Ollama Honeypot

A lightweight Go application that mimics the Ollama API to detect and log unauthorized attempts to access Ollama services.

## Problem Statement

Ollama is a popular tool for running LLMs locally. As its adoption grows, it becomes a target for attackers looking to exploit exposed API endpoints. Users need a way to detect if their exposed ports are being scanned or targeted by malicious actors seeking to use their compute resources or steal data.

## Solution

`ollama_honeypot` is a lightweight Go application that mimics the Ollama API. It listens on the standard Ollama port (11434) and responds to common API requests (like generating text or listing models) with fake data. Crucially, it logs every request, capturing details about the attacker's intent, such as the model requested and the prompt used.

## Features

- **Mock API Endpoints:**
  - `GET /`: Health check/status.
  - `POST /api/generate`: Mimic text generation with streaming NDJSON responses.
  - `POST /api/chat`: Mimic chat completion.
  - `GET /api/tags`: List fake available models.
  - `POST /api/pull`: Simulate model pulling with configurable download speeds.
  - `DELETE /api/delete`: Mimic model deletion.
  - `GET /api/ps`: Mimic running model status listing.
  - `GET /api/show`: Show model information.
  - `GET /api/version`: Get version information.

- **Request Logging:** Capture IP, timestamp, endpoint, method, and request body (prompts) with structured JSON logging.
- **Configuration:** Support for custom ports, logging preferences, and mock data paths via environment variables or CLI flags.
- **Low Resource Usage:** Designed to be lightweight and run alongside other services.

## Installation

1. Ensure you have Go 1.25.4 or later installed.

2. Clone the repository:
   ```bash
   git clone https://github.com/pedxyuyuko/ollama_honeypot.git
   cd ollama_honeypot
   ```

3. (Optional) Preseed mock data files:
   ```bash
   mkdir -p mock
   curl -o mock/tags.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/tags.json
   curl -o mock/response.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/response.json
   curl -o mock/version.json https://raw.githubusercontent.com/pedxyuyuko/ollama_honeypot/refs/heads/master/mock/version.json
   ```

4. Build the application:
   ```bash
   go build -o ollama_honeypot .
   ```

## Docker

You can run the honeypot using Docker for easy deployment and isolation.

### Using Docker Compose

1. Ensure Docker and Docker Compose are installed.

2. Copy `.example.env` to `.env` and configure environment variables as needed.

3. Run the honeypot:

   ```bash
   docker-compose up
   ```

   This will start the honeypot on port 11434, mounting the `mock` and `logs` directories for persistent data.

### Building and Running Manually

1. Build the Docker image:

   ```bash
   docker build -t ollama_honeypot .
   ```

2. Run the container:

   ```bash
   docker run -p 11434:11434 -v $(pwd)/mock:/app/mock -v $(pwd)/logs:/app/logs --env-file .env ollama_honeypot
   ```

   - `-p 11434:11434`: Maps the container's port 11434 to the host's port 11434.
   - `-v $(pwd)/mock:/app/mock`: Mounts the local `mock` directory to persist mock data.
   - `-v $(pwd)/logs:/app/logs`: Mounts the local `logs` directory to persist logs.
   - `--env-file .env`: Loads environment variables from the `.env` file.

## Usage

Run the honeypot server:

```bash
./ollama_honeypot serve
```

By default, it starts on port 11434 (Ollama's standard port).

### CLI Options

- `--port`: Specify the port to bind to (default: 11434)
- `--log-path`: Path to the audit log file (optional, logs to console if not set)
- `--mock-path`: Path to the directory containing mock data files (default: ./mock)
- `--help`: Show help information

## Configuration

Configuration can be done via environment variables or CLI flags. Environment variables can be set in a `.env` file.

Copy `.example.env` to `.env` and modify the values as needed:

```bash
cp .example.env .env
```

### Environment Variables

- `PORT`: Port to bind the server to (default: 11434)
- `MOCK_PATH`: Path to the directory containing mock data files (default: ./mock)
- `LOG_PATH`: Path to the audit log file (optional, if not set, logs only to console)
- `DEBUG`: Enable debug mode (set to 1 for debug logging, default: 0)
- `DOWNLOAD_SPEED`: Download speed in bytes per second for simulating model pulls (default: 52428800, 50MB/s)
- `DOWNLOAD_SPEED_VARIANCE`: Variance factor for random speed fluctuations (0.0 to 1.0, default: 0.2)
- `DOWNLOAD_SPEED_WAVE_PERIOD`: Period in seconds for sinusoidal speed variation (default: 1.0)
- `DOWNLOAD_SPEED_WAVE_AMPLITUDE`: Amplitude factor for sinusoidal speed variation (0.0 to 1.0, default: 0.5)

## Logging

The application uses Logrus for structured logging. Request details (IP, timestamp, method, path, body) are logged for analysis.

- Console output: Text format for readability.
- File output (if LOG_PATH is set): JSON format for parsing.

## Dependencies

- `github.com/gin-gonic/gin`: HTTP web framework.
- `github.com/spf13/cobra`: CLI application structure.
- `github.com/sirupsen/logrus`: Structured logging library.
- `github.com/joho/godotenv`: Environment variable loading.

## Contributing

Contributions are welcome! Please open issues or submit pull requests on GitHub.

## License

This project is licensed under the terms specified in the LICENSE file.