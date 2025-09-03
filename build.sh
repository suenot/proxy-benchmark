#!/bin/bash

# Build script for proxy-benchmark

set -e

echo "Building proxy-benchmark..."

# Clean previous builds
rm -rf build/
mkdir -p build/

# Get version info
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Build for multiple platforms
echo "Building for Linux amd64..."
GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/proxy-benchmark-linux-amd64 .

echo "Building for Linux arm64..."
GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o build/proxy-benchmark-linux-arm64 .

echo "Building for macOS amd64..."
GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/proxy-benchmark-darwin-amd64 .

echo "Building for macOS arm64..."
GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o build/proxy-benchmark-darwin-arm64 .

echo "Building for Windows amd64..."
GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o build/proxy-benchmark-windows-amd64.exe .

echo "Build completed! Binaries are in the build/ directory:"
ls -la build/