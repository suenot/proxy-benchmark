# Environment Variable Configuration

This project supports loading proxy credentials from environment variables to avoid hardcoding sensitive information in configuration files.

## Setup

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

## Environment Variables

Only one environment variable is needed:

### Proxy Configuration
- `TEST_PROXIES` - Comma-separated list of proxy credentials
  - Format: `protocol:host:port:username:password:enabled,protocol:host:port:username:password:enabled`
  - Example: `http:proxy1.com:8080:user:pass:enabled,socks:proxy2.com:1080:user:pass:enabled`

### Target URLs
Target URLs are hardcoded in the templates and don't require environment variables:
- GitHub API: `https://api.github.com/users/octocat`
- JSONPlaceholder: `https://jsonplaceholder.typicode.com/posts/1`
- Default test: `https://httpbin.org/get`

## Template System

Configuration templates use environment variable substitution:

- `test-configs/*.template.json` - Template files with `${ENV_VAR}` placeholders
- `test-configs/*.json` - Generated from templates (gitignored)

## Usage in Code

The `env_config.go` file provides helper functions:

```go
// Get proxies from TEST_PROXIES environment variable
proxies := GetTestProxies([]string{"fallback-proxy"})

// Load complete test configuration
config := LoadGitHubTestConfig()
```

## Security

- ✅ `.env` files are gitignored
- ✅ Template files contain no credentials
- ✅ Generated config files are gitignored
- ✅ Example files are safe to commit

## Testing

Run environment-based tests:

```bash
# With real credentials
./test_with_env.sh

# Validation tests only
go test -v -run TestValidateResponse

# Manual configuration
TEST_PROXIES="your-proxy" go run . -config test-configs/custom-config.json
```

## Example

```bash
# 1. Copy example file
cp .env.example .env

# 2. Edit with your real proxies
echo 'TEST_PROXIES=http:real-proxy.com:8080:user:pass:enabled' > .env

# 3. Run tests
./test_with_env.sh
```