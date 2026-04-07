# Use Golang image for building
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build API
RUN go build -o api ./cmd/api/main.go

# Build Worker
RUN go build -o worker ./cmd/worker/main.go

# Use a small image for running
FROM alpine:latest

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/api /app/worker ./
# Copy UI files
COPY --from=builder /app/ui ./ui

# Create directories for storage
RUN mkdir -p uploads processed logs

# Default command (can be overridden in docker-compose)
CMD ["./api"]
