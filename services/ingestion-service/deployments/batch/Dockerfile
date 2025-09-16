# Build stage
FROM golang:1.21-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build all batch commands
RUN mkdir -p /app/bin && \
    go build -o /app/bin/batch-trending ./cmd/batch/trending && \
    go build -o /app/bin/batch-schedule-snapshots ./cmd/batch/schedule-snapshots && \
    go build -o /app/bin/batch-websub-renewal ./cmd/batch/websub-renewal && \
    go build -o /app/bin/batch-rankings ./cmd/batch/rankings

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

# Copy binaries from builder
COPY --from=builder /app/bin /app/bin

# Set ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Set working directory
WORKDIR /app

# Default command (can be overridden)
CMD ["/app/bin/batch-trending"]