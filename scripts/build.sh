#!/bin/bash

# Build script for gohabit with version information injection
# Usage: ./scripts/build.sh

set -e

# Get build information
VERSION=${VERSION:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_REF=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Output directory
OUTPUT_DIR="bin"
BINARY_NAME="gohabit"

# Create output directory if it doesn't exist
mkdir -p "${OUTPUT_DIR}"

echo "Building ${BINARY_NAME}..."
echo "Version: ${VERSION}"
echo "Git Commit: ${GIT_COMMIT}"
echo "Git Ref: ${GIT_REF}"
echo "Git Tag: ${GIT_TAG}"
echo "Build Date: ${BUILD_DATE}"

# Build with ldflags
go build -ldflags "\
    -X 'github.com/hasnpr/gohabit/internal/app.Version=${VERSION}' \
    -X 'github.com/hasnpr/gohabit/internal/app.GitCommit=${GIT_COMMIT}' \
    -X 'github.com/hasnpr/gohabit/internal/app.GitRef=${GIT_REF}' \
    -X 'github.com/hasnpr/gohabit/internal/app.GitTag=${GIT_TAG}' \
    -X 'github.com/hasnpr/gohabit/internal/app.BuildDate=${BUILD_DATE}'" \
    -o "${OUTPUT_DIR}/${BINARY_NAME}" .

echo "Build completed: ${OUTPUT_DIR}/${BINARY_NAME}"