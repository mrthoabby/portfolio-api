# =============================================================================
# Portfolio API - Production Dockerfile
# =============================================================================
# Build: docker build -t portfolio-api .
# Run:   docker run -p 8080:8080 \
#          -e DATABASE_URL=mongodb://... \
#          -e DATABASE_NAME=portfolio_db \
#          -e ALLOWED_ORIGINS=https://yoursite.com \
#          portfolio-api
# =============================================================================

# Build stage
FROM golang:1.25-alpine AS builder

# Install security updates and git (for go mod and version detection)
RUN apk update && apk upgrade && apk add --no-cache ca-certificates git

WORKDIR /app

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code (including .git for version detection)
COPY . .

# =============================================================================
# Semantic Versioning Script
# =============================================================================
# This script calculates the version based on Conventional Commits:
# - fix/fixed/hotfix/patch: → increments PATCH (x.x.1 → x.x.2)
# - feat/feature/refactor:  → increments MINOR (x.1.x → x.2.0)
# - MAJOR: only incremented manually via git tag
#
# Examples:
#   - On tag v1.2.0:                    → v1.2.0
#   - 3 fix commits after v1.2.0:       → v1.2.3
#   - 2 feat + 1 fix after v1.2.0:      → v1.4.1
# =============================================================================

RUN VERSION="" && \
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    BUILD_DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    \
    # Check if we're exactly on a tag
    EXACT_TAG=$(git describe --tags --exact-match HEAD 2>/dev/null || echo "") && \
    \
    if [ -n "$EXACT_TAG" ]; then \
        # We're exactly on a tag, use it as-is
        VERSION="$EXACT_TAG"; \
        echo ""; \
        echo "========================================"; \
        echo "  Building RELEASE version: $VERSION"; \
        echo "  Commit: $GIT_COMMIT"; \
        echo "  Date: $BUILD_DATE"; \
        echo "========================================"; \
        echo ""; \
    else \
        # We're not on a tag, calculate version from commits
        LAST_TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0") && \
        \
        # Parse current version components
        TAG_CLEAN=$(echo "$LAST_TAG" | sed 's/^v//') && \
        MAJOR=$(echo "$TAG_CLEAN" | cut -d. -f1) && \
        MINOR=$(echo "$TAG_CLEAN" | cut -d. -f2) && \
        PATCH=$(echo "$TAG_CLEAN" | cut -d. -f3) && \
        \
        # Ensure we have valid numbers
        MAJOR=${MAJOR:-0} && \
        MINOR=${MINOR:-0} && \
        PATCH=${PATCH:-0} && \
        \
        # Count commits by type since last tag
        # MINOR incrementing commits: feat, feature, refactor
        MINOR_COMMITS=$(git log ${LAST_TAG}..HEAD --oneline 2>/dev/null | \
            grep -iE "^[a-f0-9]+ (feat|feature|refactor)(\(.*\))?:" | wc -l | tr -d ' ') && \
        \
        # PATCH incrementing commits: fix, fixed, hotfix, patch, bugfix
        PATCH_COMMITS=$(git log ${LAST_TAG}..HEAD --oneline 2>/dev/null | \
            grep -iE "^[a-f0-9]+ (fix|fixed|hotfix|patch|bugfix)(\(.*\))?:" | wc -l | tr -d ' ') && \
        \
        # Calculate new version
        # If there are MINOR commits, increment MINOR and reset PATCH to count of fixes after
        if [ "$MINOR_COMMITS" -gt 0 ]; then \
            MINOR=$((MINOR + MINOR_COMMITS)); \
            PATCH=$PATCH_COMMITS; \
        else \
            # Only PATCH commits (or other commits)
            PATCH=$((PATCH + PATCH_COMMITS)); \
        fi && \
        \
        # Build version string with commit info for non-release builds
        VERSION="v${MAJOR}.${MINOR}.${PATCH}-dev.${GIT_COMMIT}"; \
        \
        echo ""; \
        echo "========================================"; \
        echo "  Building DEV version: $VERSION"; \
        echo "  Based on tag: $LAST_TAG"; \
        echo "  Minor commits (feat/refactor): $MINOR_COMMITS"; \
        echo "  Patch commits (fix/hotfix): $PATCH_COMMITS"; \
        echo "  Commit: $GIT_COMMIT"; \
        echo "  Date: $BUILD_DATE"; \
        echo "========================================"; \
        echo ""; \
    fi && \
    \
    # Build the binary with version info
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
        -ldflags="-w -s -extldflags '-static' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.Version=${VERSION}' \
            -X 'github.com/mrthoabby/portfolio-api/internal/version.GitCommit=${GIT_COMMIT}' \
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

# Expose port (configurable via PORT env var, default 8080)
EXPOSE 8080

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
