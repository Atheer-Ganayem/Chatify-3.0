# Stage 1: Build the Go binary
ARG GO_VERSION=1.23.0
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the rest of the code
COPY . .

# Build the app binary
RUN go build -v -o /run-app .

# Stage 2: Run in a minimal Debian environment with CA certs
FROM debian:bookworm

# Install CA certificates to avoid TLS errors (MongoDB, S3, etc.)
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the built binary from the builder stage
COPY --from=builder /run-app /usr/local/bin/

# Expose the app port (optional; not required by Fly)
EXPOSE 8080

# Run the app
CMD ["run-app"]
