# Product

## Why this project exists
Ollama is a popular tool for running large language models locally. However, exposing Ollama APIs can pose security risks if not properly secured. This honeypot serves as a security monitoring tool to detect attempts to access or exploit Ollama instances.

## Problems it solves
- Detect unauthorized access attempts to Ollama APIs
- Log suspicious activity for security analysis
- Provide a decoy to divert attackers from real Ollama instances
- Help in understanding attack patterns on AI/ML services

## How it should work
The honeypot runs an HTTP server on port 11434 that responds to API requests similar to Ollama. All requests are logged with details like IP, headers, body, etc. The responses mimic legitimate Ollama responses to avoid detection.

## User experience goals
- Easy to deploy and run
- Minimal resource usage
- Comprehensive logging
- Configurable responses
- Integration with security monitoring tools