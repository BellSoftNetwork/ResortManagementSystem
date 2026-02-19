# Multi-stage Dockerfile for api-core
# Build arguments for multi-platform support
ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

FROM golang:1.23-alpine AS base

# Install base packages
RUN apk add --no-cache git make bash curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# Development stage
FROM base AS development

# Install development tools
RUN go install github.com/air-verse/air@v1.52.3
RUN go install github.com/go-delve/delve/cmd/dlv@v1.24.0

# Copy air config
COPY .air.toml .

# Copy source code
COPY . .

# Expose ports
EXPOSE 8080 2345

# Environment variables
ENV CGO_ENABLED=0
ENV PROFILE=local

# Run air for hot reload
CMD ["air", "-c", ".air.toml"]

# Builder stage
FROM base AS builder

# Copy source code
COPY . .

# Copy migration files
RUN mkdir -p /app/internal/migrations
COPY internal/migrations/*.sql /usr/local/share/migrations/

# Build arguments for target platform
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# Build the application
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o main ./cmd/server/main.go

# Build the migration tool
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o migrate ./cmd/migrate/main.go

# Production stage
FROM alpine:latest AS production

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/migrate /usr/local/bin/migrate

# Copy configuration files
COPY --from=builder /app/config ./config

# Create directory for static files
RUN mkdir -p public

# Change ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Environment variables will be set by Kubernetes deployment

# Run the application
CMD ["./main"]
