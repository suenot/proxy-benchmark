# Installation & Setup

<cite>
**Referenced Files in This Document**   
- [README.md](file://README.md)
- [Makefile](file://Makefile)
- [build.sh](file://build.sh)
- [go.mod](file://go.mod)
- [main.go](file://main.go)
- [.env.example](file://.env.example) - *Added in recent commit*
- [ENV_CONFIG.md](file://ENV_CONFIG.md) - *Added in recent commit*
- [env_config.go](file://env_config.go) - *Added in recent commit*
- [test_with_env.sh](file://test_with_env.sh) - *Added in recent commit*
</cite>

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Building from Source](#building-from-source)
3. [Binary Distribution Methods](#binary-distribution-methods)
4. [Environment Configuration](#environment-configuration)
5. [Build Process Integration](#build-process-integration)
6. [Common Setup Issues](#common-setup-issues)
7. [Troubleshooting Tips](#troubleshooting-tips)
8. [Performance Considerations](#performance-considerations)

## Prerequisites

Before installing the proxy-benchmark tool, ensure your system meets the following requirements:

- **Go version 1.19 or higher**: The project requires Go 1.23.0 as specified in go.mod, with toolchain compatibility for Go 1.24.6
- **Git**: Required for cloning the repository from GitHub
- **Make (optional)**: While not mandatory, Make simplifies the build process through the provided Makefile
- **External dependencies**: The tool relies on two key external packages:
  - `golang.org/x/net v0.43.0` for SOCKS5 proxy implementation
  - `github.com/montanaflynn/stats v0.7.1` for statistical calculations

These dependencies are automatically managed through Go modules and will be downloaded during the build process.

**Section sources**
- [go.mod](file://go.mod#L1-L10)
- [README.md](file://README.md#L15-L20)

## Building from Source

The proxy-benchmark tool can be compiled from source using multiple methods. The recommended approach utilizes the provided Makefile for simplicity and consistency.

### Using Makefile

```bash
# Clone the repository
git clone https://github.com/suenot/proxy-benchmark.git
cd proxy-benchmark

# Build the application using Make
make build
```

The Makefile defines several targets that streamline development and deployment:
- `build`: Compiles the main binary as `proxy-benchmark`
- `test`: Runs all unit tests with verbose output
- `test-cover`: Executes tests with coverage reporting
- `fmt`: Formats code according to Go standards
- `vet`: Analyzes code for potential errors
- `clean`: Removes compiled binaries
- `deps`: Ensures all dependencies are properly installed via `go mod tidy`
- `build-all`: Creates binaries for Linux, Windows, and macOS platforms

The default target is `build`, so running `make` without arguments will compile the application.

### Using build.sh Script

For comprehensive cross-platform builds with version information, use the build.sh script:

```bash
# Execute the build script
./build.sh
```

This script performs the following operations:
1. Cleans previous build artifacts by removing the `build/` directory
2. Creates a fresh `build/` directory
3. Retrieves version information from Git (tag, commit hash, dirty status)
4. Captures build timestamp
5. Compiles binaries for multiple platforms with embedded version metadata:
   - Linux AMD64 and ARM64
   - macOS AMD64 and ARM64 (Apple Silicon)
   - Windows AMD64
6. Outputs all binaries to the `build/` directory

The script uses ldflags to embed version information into the binary, making it easier to track which version is deployed.

### Direct Go Build

Alternatively, you can build directly using the Go command:

```bash
# Direct compilation
go build -o proxy-benchmark .

# Or simply
go build
```

This method compiles the application using the main package defined in main.go and produces a binary named after the current directory.

**Section sources**
- [Makefile](file://Makefile#L1-L50)
- [build.sh](file://build.sh#L1-L37)
- [main.go](file://main.go#L1-L81)

## Binary Distribution Methods

While building from source is recommended for customization and verification, the project supports multiple binary distribution approaches for convenience.

### Cross-Platform Compilation

The Makefile provides platform-specific build targets:
```bash
# Linux
make build-linux

# Windows
make build-windows

# macOS (Intel)
make build-mac
```

These targets use Go's cross-compilation capabilities to generate binaries for different operating systems and architectures without requiring native compilation environments.

### Comprehensive Build Script

The build.sh script offers a more complete solution for distributing binaries across platforms. It generates:
- `proxy-benchmark-linux-amd64`
- `proxy-benchmark-linux-arm64`
- `proxy-benchmark-darwin-amd64`
- `proxy-benchmark-darwin-arm64`
- `proxy-benchmark-windows-amd64.exe`

All binaries include embedded version information accessible at runtime, enabling better tracking and debugging in production environments.

### Version Information

Both build methods capture important metadata:
- **Version**: Git tag or "dev" if no tag exists
- **Build Time**: UTC timestamp of compilation
- **Git Commit**: Short hash of the latest commit

This information helps identify exactly which codebase was used for any given binary, crucial for troubleshooting and reproducibility.

**Section sources**
- [build.sh](file://build.sh#L1-L37)
- [Makefile](file://Makefile#L40-L50)

## Environment Configuration

Proper environment setup ensures smooth operation of the proxy-benchmark tool.

### PATH Configuration

After building, add the binary location to your system PATH:

```bash
# Add to PATH temporarily
export PATH=$PATH:/path/to/proxy-benchmark

# Add permanently (add to shell profile)
echo 'export PATH=$PATH:/path/to/proxy-benchmark' >> ~/.bashrc
# or for macOS users
echo 'export PATH=$PATH:/path/to/proxy-benchmark' >> ~/.zshrc
```

This allows execution of the tool from any directory using `proxy-benchmark`.

### Configuration File Setup

Create a configuration file based on the provided example:

```bash
# Copy example configuration
cp config.example.json config.json
```

Edit `config.json` to specify your proxy list and benchmark parameters. The tool expects this file in the working directory by default, but you can specify alternative locations using the `-config` flag.

### Secure Credential Management via Environment Variables

To securely manage proxy credentials without hardcoding them in configuration files, the tool now supports loading from environment variables.

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your real proxy credentials:**
   ```bash
   # Test proxy credentials (comma-separated list)
   TEST_PROXIES=http:your-proxy1.com:8080:user:pass:enabled,socks:your-proxy2.com:1080:user:pass:enabled
   ```

3. **Use environment-based testing:**
   ```bash
   ./test_with_env.sh
   ```

The `env_config.go` file provides helper functions like `GetTestProxies()` and `LoadGitHubTestConfig()` that load configurations using environment variables. Template configuration files (`*.template.json`) support `${TEST_PROXIES}` placeholders, which are substituted at runtime.

Security features:
- ✅ `.env` files are gitignored
- ✅ Template files contain no credentials
- ✅ Generated config files are gitignored
- ✅ Example files are safe to commit

**Section sources**
- [README.md](file://README.md#L35-L100)
- [config.example.json](file://config.example.json)
- [.env.example](file://.env.example)
- [ENV_CONFIG.md](file://ENV_CONFIG.md)
- [env_config.go](file://env_config.go#L9-L108)
- [test_with_env.sh](file://test_with_env.sh#L1-L105)

## Build Process Integration

The build infrastructure integrates seamlessly with the project structure and development workflow.

### Project Structure Alignment

The build system aligns with the modular code organization:
- **main.go**: Entry point compiled into executable
- **config.go**: Configuration loading functionality
- **benchmark.go**: Core benchmarking engine
- **proxy.go**: Proxy parsing and management
- **http_client.go**, **socks5_client.go**: Protocol-specific clients
- **metrics.go**, **statistics.go**: Data collection and analysis
- **reporter.go**: Report generation

The simple `go build` command compiles all these components into a single binary.

### Dependency Management

Go modules handle external dependencies automatically:
```bash
# Ensure dependencies are up-to-date
go mod tidy
```

This command downloads required packages and updates go.mod and go.sum files accordingly. The indirect dependencies (`golang.org/x/net` and `github.com/montanaflynn/stats`) are resolved transitively.

### Development Workflow

The Makefile supports a complete development cycle:
```bash
# Format code
make fmt

# Check for issues
make vet

# Run tests
make test

# Build final binary
make build
```

This standardized workflow ensures code quality and consistency across development environments.

**Section sources**
- [Makefile](file://Makefile#L1-L50)
- [go.mod](file://go.mod#L1-L10)
- [README.md](file://README.md#L300-L332)

## Common Setup Issues

Several common issues may arise during installation and setup.

### Missing Go Modules

If dependencies fail to download:
```bash
# Clear module cache
go clean -modcache

# Re-download dependencies
go mod download

# Tidy module files
go mod tidy
```

Network restrictions or proxy settings might interfere with module downloads. Set appropriate environment variables if behind a corporate firewall:
```bash
export GOPROXY=https://proxy.golang.org,direct
export GONOSUMDB=your-private-repo.com
```

### Permission Errors

Permission issues during compilation typically occur when writing to protected directories:
```bash
# Ensure write permissions in project directory
chmod -R u+w .

# Or build to user-writable location
go build -o ~/bin/proxy-benchmark
```

On Unix-like systems, avoid using sudo with go build as it can create permission conflicts later.

### Go Version Mismatch

Verify your Go version meets requirements:
```bash
# Check current version
go version

# Expected output format
# go version go1.23.x darwin/amd64
```

If the version is too old, update Go from the official website or package manager. The project specifies Go 1.23.0 in go.mod, so older versions may encounter compatibility issues.

### Missing Git Repository

The build.sh script attempts to retrieve Git metadata. If working outside a Git repository:
```bash
# Initialize Git repository
git init
git add .
git commit -m "Initial commit"

# Or modify build.sh to handle non-Git environments
```

Alternatively, the script will use "dev" as the version identifier when Git is unavailable.

**Section sources**
- [go.mod](file://go.mod#L3-L3)
- [build.sh](file://build.sh#L6-L10)
- [README.md](file://README.md#L250-L260)

## Troubleshooting Tips

Effective troubleshooting strategies help resolve common setup problems.

### Verify Installation

After building, confirm successful compilation:
```bash
# Check if binary exists and is executable
ls -la proxy-benchmark

# Display file information
file proxy-benchmark

# Test basic execution
./proxy-benchmark -h
```

A successful build should produce an executable binary that runs without immediate errors.

### Configuration Validation

Ensure the configuration file is correctly formatted:
```bash
# Validate JSON syntax
cat config.json | python -m json.tool

# Or use jq
jq . config.json
```

Common configuration issues include:
- Invalid proxy string format (must be protocol:host:port:username:password:status)
- Missing required fields in benchmark settings
- Incorrect JSON syntax or encoding

### Dependency Verification

Confirm all dependencies are properly resolved:
```bash
# List direct and indirect dependencies
go list -m all

# Check for missing modules
go mod verify

# Download missing modules
go mod download
```

The go.mod file shows both direct and indirect dependencies, with the actual required packages being pulled in transitively.

### Build Script Debugging

Enable verbose output in the build script by modifying build.sh:
```bash
#!/bin/bash
set -x  # Add this line for debug output
set -e
```

This reveals each command executed and helps identify where failures occur during the build process.

**Section sources**
- [config.go](file://config.go#L32-L47)
- [proxy.go](file://proxy.go#L8-L15)
- [build.sh](file://build.sh#L2-L5)

## Performance Considerations

Different operating systems and configurations impact build and runtime performance.

### Operating System Differences

- **Linux/macOS**: Generally faster build times due to efficient filesystem and process handling
- **Windows**: Slightly slower builds due to longer path resolution and antivirus scanning
- **ARM vs AMD64**: Native compilation on Apple Silicon (ARM64) may be faster than cross-compilation

### Build Optimization

To optimize build performance:
```bash
# Use Go build cache
export GOCACHE=$HOME/.cache/go-build

# Enable parallel compilation
export GOMAXPROCS=$(nproc)

# Use incremental builds when possible
go build  # Subsequent builds are faster
```

The build.sh script already implements parallelization through concurrent goroutines in the benchmark engine, but build performance benefits from proper Go environment tuning.

### Resource Usage

During compilation:
- **Memory**: Requires approximately 500MB-1GB RAM
- **CPU**: Benefits from multiple cores for faster compilation
- **Disk Space**: Binary size varies by platform (typically 10-20MB)

For large-scale deployments, consider pre-building binaries on powerful machines and distributing them to target systems.

**Section sources**
- [build.sh](file://build.sh#L1-L37)
- [benchmark.go](file://benchmark.go#L18-L36)
- [README.md](file://README.md#L250-L260)