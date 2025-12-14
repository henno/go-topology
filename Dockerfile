# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server ./cmd/server

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server /app/server

# Copy web assets and config
COPY web /app/web
COPY config.json /app/config.json

# Expose port
EXPOSE 9090

# Run the server
CMD ["/app/server"]
