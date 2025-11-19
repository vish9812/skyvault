# Stage 1: Build the web frontend
FROM node:20-alpine AS web-builder

WORKDIR /build/web

# Install pnpm
RUN npm install -g pnpm

# Copy web package files
COPY web/package.json web/pnpm-lock.yaml* ./

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy web source code
COPY web/ ./

# Build the web application
RUN pnpm run build

# Stage 2: Build the Go backend
FROM golang:1.23-alpine AS server-builder

WORKDIR /build/server

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY server/go.mod server/go.sum ./

# Download dependencies
RUN go mod download

# Copy server source code
COPY server/ ./

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o bin/server ./cmd/main.go

# Stage 3: Final runtime image
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Copy server binary from builder
COPY --from=server-builder /build/server/bin/server /app/server

# Copy migrations
COPY --from=server-builder /build/server/internal/infrastructure/internal/repository/internal/migrations /app/migrations

# Copy web build output
COPY --from=web-builder /build/web/dist /app/static

# Create data directory
RUN mkdir -p /app/data

# Set default environment variables for container-internal paths
ENV SERVER__PATH=/app \
    SERVER__DATA_DIR=/app/data \
    SERVER__MIGRATIONS_DIR=/app/migrations \
    SERVER__ADDR=0.0.0.0:8090

# Expose the application port
EXPOSE 8090

# Add healthcheck
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8090/api/v1/pub/health || exit 1

CMD ["/app/server"]
