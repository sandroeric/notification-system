# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum (if exists) and download dependencies
COPY go.mod ./
# COPY go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o notification-server ./cmd/server

# Final stage
FROM alpine:latest

WORKDIR /app

# Ensure tzdata is available and certificates
RUN apk --no-cache add ca-certificates tzdata

# Create data directory for SQLite
RUN mkdir -p /app/data

# Copy built binary from builder
COPY --from=builder /app/notification-server .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./notification-server"]
