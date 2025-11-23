# Use the official Go image as the base image for building
FROM golang:1.25.4-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o ollama_honeypot .

# Use a minimal base image for the runtime
FROM alpine:latest

# Install ca-certificates for HTTPS requests if needed
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/ollama_honeypot .

# Expose the default port
EXPOSE 11434

# Set the entrypoint to run the application
ENTRYPOINT ["./ollama_honeypot"]

# Default command to serve
CMD ["serve"]