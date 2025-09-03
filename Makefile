# Makefile for Proxy Benchmark Utility

# Build variables
BINARY=proxy-benchmark
MAIN_FILE=main.go

# Build the application
build:
	go build -o ${BINARY} .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Clean build artifacts
clean:
	rm -f ${BINARY}

# Install dependencies
deps:
	go mod tidy

# Build for multiple platforms
build-all: build-linux build-windows build-mac

build-linux:
	GOOS=linux GOARCH=amd64 go build -o ${BINARY}-linux .

build-windows:
	GOOS=windows GOARCH=amd64 go build -o ${BINARY}-windows.exe .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY}-mac .

# Default target
default: build

.PHONY: build test test-cover fmt vet clean deps build-all build-linux build-windows build-mac default