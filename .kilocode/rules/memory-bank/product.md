# Product

## Problem Statement
Ollama is a popular tool for running LLMs locally. As its adoption grows, it becomes a target for attackers looking to exploit exposed API endpoints. Users need a way to detect if their exposed ports are being scanned or targeted by malicious actors seeking to use their compute resources or steal data.

## Solution
`ollama_honeypot` is a lightweight Go application that mimics the Ollama API. It listens on the standard Ollama port (11434) and responds to common API requests (like generating text or listing models) with fake data. Crucially, it logs every request, capturing details about the attacker's intent, such as the model requested and the prompt used.

## User Experience Goals
1.  **Zero Config Start:** Ideally, the user can just run the binary and it starts listening on the default port.
2.  **High Fidelity Deception:** The API responses should be realistic enough to keep automated scanners and simple scripts engaged.
3.  **Clear Visibility:** Logs should be easy to read and parse, providing immediate insight into potential threats.
4.  **Low Resource Usage:** As a honeypot, it should consume minimal system resources.

## Key Features
-   **Mock API Endpoints:**
    -   `GET /`: Health check/status.
    -   `POST /api/generate`: Mimic text generation.
    -   `POST /api/chat`: Mimic chat completion.
    -   `GET /api/tags`: List fake available models.
-   **Request Logging:** Capture IP, timestamp, endpoint, method, and request body (prompts).
-   **Configuration:** Support for custom ports and logging preferences via environment variables or CLI flags.