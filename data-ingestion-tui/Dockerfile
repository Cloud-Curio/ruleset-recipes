# Data Ingestion TUI - Docker Configuration
# 
# Multi-stage Docker build for the Data Ingestion TUI application.
# Optimized for production deployment with minimal image size.
# 
# Author: Codegen AI Assistant
# Created: 2024
# License: MIT

# Build stage
FROM golang:1.21-alpine AS builder

# Set build arguments
ARG VERSION=1.0.0
ARG BUILD_TIME
ARG GIT_COMMIT

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    sqlite \
    build-base

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags "-X main.version=${VERSION} -X main.buildTime=${BUILD_TIME} -X main.gitCommit=${GIT_COMMIT} -w -s" \
    -a -installsuffix cgo \
    -o dit \
    ./cmd/main.go

# Runtime stage
FROM alpine:3.18

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite \
    curl \
    bash \
    && rm -rf /var/cache/apk/*

# Create non-root user
RUN addgroup -g 1001 -S dit && \
    adduser -u 1001 -S dit -G dit

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/dit /usr/local/bin/dit

# Copy configuration files
COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/migrations /app/migrations

# Create data directory
RUN mkdir -p /app/data && \
    chown -R dit:dit /app

# Switch to non-root user
USER dit

# Set environment variables
ENV DIT_DATA_DIR=/app/data
ENV DIT_CONFIG_FILE=/app/configs/config.yaml
ENV DIT_LOG_LEVEL=info

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose ports
EXPOSE 2222 8080

# Set entrypoint
ENTRYPOINT ["dit"]

# Default command
CMD ["server", "--host", "0.0.0.0", "--port", "2222"]

# Labels
LABEL maintainer="Codegen AI Assistant" \
      version="${VERSION}" \
      description="Data Ingestion TUI - Comprehensive data management interface" \
      org.opencontainers.image.title="Data Ingestion TUI" \
      org.opencontainers.image.description="A comprehensive TUI for managing data ingestion from government and public APIs" \
      org.opencontainers.image.version="${VERSION}" \
      org.opencontainers.image.created="${BUILD_TIME}" \
      org.opencontainers.image.revision="${GIT_COMMIT}" \
      org.opencontainers.image.source="https://github.com/your-org/data-ingestion-tui" \
      org.opencontainers.image.licenses="MIT"
