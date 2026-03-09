# Baboon - Typing Practice Application
# Multi-stage Docker build for minimal production image

# =============================================================================
# Build stage - Compile Go binary
# =============================================================================
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -a -installsuffix cgo \
    -o baboon .

# =============================================================================
# Production stage - Minimal runtime image
# =============================================================================
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata wget

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/baboon .

# Copy web frontend (must be pre-built: cd web && npm run build)
COPY web/dist ./web/dist

# Create non-root user for security
RUN addgroup -g 1000 baboon && \
    adduser -D -u 1000 -G baboon baboon && \
    mkdir -p /home/baboon/.config/baboon && \
    chown -R baboon:baboon /home/baboon

# Switch to non-root user
USER baboon

# Expose the default port
EXPOSE 8787

# Health check endpoint
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8787/api/health || exit 1

# Default entrypoint runs web server mode
ENTRYPOINT ["./baboon"]
CMD ["web", "-port", "8787"]
