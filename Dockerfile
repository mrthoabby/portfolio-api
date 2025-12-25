FROM golang:1.25-alpine AS builder

# Install ca-certificates package (needed for TLS connections to MongoDB Atlas)
# Using --update-cache ensures we get the latest certificates without full upgrade
RUN apk add --no-cache --update-cache ca-certificates

WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./

# Download dependencies using BuildKit cache mount for faster rebuilds
# The cache persists between builds, so dependencies are only downloaded once
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download && go mod verify

# Copy source code (order matters: copy directories that change less frequently first)
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# =============================================================================
# Build binary - Optimized with BuildKit cache mounts
# =============================================================================
# Using BuildKit cache mounts for both module cache and build cache
# This significantly speeds up subsequent builds by reusing compiled packages

# Build arguments (can be overridden from build command)
ARG VERSION=unknown
ARG BUILD_DATE

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    BUILD_DATE=${BUILD_DATE:-$(date -u +%Y-%m-%dT%H:%M:%SZ)} && \
    echo "Building Portfolio API (version: ${VERSION}, built: $BUILD_DATE)" && \
    \
    # Build the binary with cached dependencies and build artifacts
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -extldflags '-static' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.Version=${VERSION}' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.BuildDate=${BUILD_DATE}'" \
        -a -installsuffix cgo \
        -o main ./cmd/api

# =============================================================================
# STAGE 2: Final production image
# =============================================================================
# This starts a NEW, fresh image (distroless) - completely separate from stage 1
# It does NOT include anything from stage 1 automatically
FROM gcr.io/distroless/static-debian12:nonroot

# Copy CA certificates bundle from the "builder" stage (stage 1)
# SYNTAX: COPY --from=STAGE_NAME /path/in/stage1 /path/in/stage2
# 
# What happens:
#   1. In stage 1 (builder), we installed ca-certificates with: apk add ca-certificates
#   2. That created the file: /etc/ssl/certs/ca-certificates.crt in the builder image
#   3. Now we COPY that file from builder image to this final image
#   4. The "builder" image still exists temporarily during build, so we can copy from it
#
# WHY: The distroless image does NOT include CA certificates by default.
#      Our Go code (database.go) needs these certificates to verify TLS 
#      connections to MongoDB Atlas (which uses HTTPS/TLS).
#      We can't install them here because distroless has NO package manager (no apk, no apt).
#      So we install them in stage 1 (where we CAN use apk) and copy them here.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the "builder" stage (stage 1)
# The binary was compiled in stage 1 and saved at /app/main
# We copy it to the same location in this final image
COPY --from=builder /app/main /app/main

# Use non-root user (UID 65532, provided by distroless:nonroot)
USER nonroot:nonroot

# =============================================================================
# Required Environment Variables (passed at runtime):
# - DATABASE_URL: MongoDB connection string
# - DATABASE_NAME: MongoDB database name  
# - ALLOWED_ORIGINS: Comma-separated CORS origins
# - PORT: Server port (required, no default)
# =============================================================================

ENTRYPOINT ["/app/main"]
