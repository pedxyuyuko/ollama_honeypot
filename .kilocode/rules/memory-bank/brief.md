# Brief

This project is an Ollama API honeypot designed to detect and log unauthorized or malicious attempts to access Ollama services. The core goal is to provide a security tool that mimics the Ollama API to identify potential threats in environments where Ollama is deployed.

Key requirements:
- Mimic Ollama API endpoints
- Log all incoming requests
- Capture unauth prompt and model
- Run on default Ollama port (11434)
- Provide basic CLI interface
- Read config file .env or system env vars
