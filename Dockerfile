# =============================================================================
# Portfolio API - Production Dockerfile
# =============================================================================
# Build: docker build -t portfolio-api .
# Run:   docker run -p HOST_PORT:CONTAINER_PORT \
#          -e DATABASE_URL=mongodb://... \
#          -e DATABASE_NAME=portfolio_db \
#          -e ALLOWED_ORIGINS=https://yoursite.com \
#          -e PORT=CONTAINER_PORT \
#          portfolio-api
# =============================================================================

# Build stage
FROM golang:1.25-alpine AS builder

# Install security updates and ca-certificates
RUN apk update && apk upgrade && apk add --no-cache ca-certificates

WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# =============================================================================
# Build binary - Fast version (no git detection)
# =============================================================================

RUN BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    echo "Building Portfolio API (built: $BUILD_DATE)" && \
    \
    # Build the binary
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -extldflags '-static' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.Version=v0.0.0' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.GitCommit=unknown' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.BuildDate=${BUILD_DATE}'" \
    -a -installsuffix cgo \
    -o main ./cmd/api

# =============================================================================
# Final stage - minimal image for security
# =============================================================================
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates for HTTPS connections (e.g., to MongoDB Atlas)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary
COPY --from=builder /app/main /app/main

# Use non-root user (UID 65532, provided by distroless:nonroot)
USER nonroot:nonroot

# =============================================================================
# Required Environment Variables (passed at runtime):
# - DATABASE_URL: MongoDB connection string
# - DATABASE_NAME: MongoDB database name  
# - ALLOWED_ORIGINS: Comma-separated CORS origins
#
# Optional Environment Variables:
# - PORT: Server port (default: 8080)
# - ENV: Environment name (default: development)
# =============================================================================

ENTRYPOINT ["/app/main"]
