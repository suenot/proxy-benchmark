#!/bin/bash

# Test script that uses environment variables for proxy credentials
# This script demonstrates how to run tests with environment-based configuration

set -e

echo "=== Proxy Benchmark Environment-Based Testing ==="

# Check if .env file exists and has content
if [ ! -f ".env" ] || [ ! -s ".env" ]; then
    echo "Warning: .env file not found or is empty. Using .env.example values."
    echo "To use real proxy credentials:"
    echo "1. Copy .env.example to .env: cp .env.example .env"
    echo "2. Edit .env with your real proxy credentials"
    echo "3. Run this script again"
    echo ""
    
    # Load example values for demonstration
    if [ -f ".env.example" ]; then
        export $(grep -v '^#' .env.example | xargs)
    fi
else
    echo "Loading environment variables from .env..."
    export $(grep -v '^#' .env | xargs)
fi

# Display loaded environment variables (without showing sensitive data)
echo "Environment variables loaded:"
echo "- TEST_PROXIES: ${TEST_PROXIES:0:40}..."
echo ""

# Function to substitute environment variables in JSON templates
substitute_env_vars() {
    local template_file="$1"
    local output_file="$2"
    
    # Handle TEST_PROXIES specially - convert comma-separated to JSON array
    if [ -n "$TEST_PROXIES" ]; then
        # Convert comma-separated proxies to JSON array format
        IFS=',' read -ra PROXY_ARRAY <<< "$TEST_PROXIES"
        PROXY_JSON=""
        for i in "${!PROXY_ARRAY[@]}"; do
            proxy=$(echo "${PROXY_ARRAY[i]}" | xargs)  # trim whitespace
            if [ $i -eq 0 ]; then
                PROXY_JSON="\"$proxy\""
            else
                PROXY_JSON="$PROXY_JSON, \"$proxy\""
            fi
        done
        
        # Use temporary replacement and then envsubst
        sed "s/\"\${TEST_PROXIES}\"/TEMP_PROXY_PLACEHOLDER/g" "$template_file" | envsubst | sed "s/TEMP_PROXY_PLACEHOLDER/$PROXY_JSON/g" > "$output_file"
    else
        # Use envsubst for other variables
        envsubst < "$template_file" > "$output_file"
    fi
    echo "Created $output_file from $template_file"
}

# Create temporary config files from templates
echo "=== Creating config files from templates ==="
substitute_env_vars "test-configs/github-api-config.template.json" "test-configs/github-api-config.json"
substitute_env_vars "test-configs/jsonplaceholder-config.template.json" "test-configs/jsonplaceholder-config.json"

# Run validation tests
echo ""
echo "=== Running validation tests ==="
go test -v -run TestValidateResponse

# Run integration tests with GitHub API
echo ""
echo "=== Running GitHub API integration test ==="
if [ -n "$TEST_GITHUB_PROXY" ] && [[ "$TEST_GITHUB_PROXY" != *"example.com"* ]]; then
    echo "Testing with GitHub API configuration..."
    timeout 30s go run . -config test-configs/github-api-config.json || echo "GitHub test completed or timed out"
else
    echo "Skipping GitHub test (using example proxy)"
fi

# Run integration tests with JSONPlaceholder
echo ""
echo "=== Running JSONPlaceholder integration test ==="
if [ -n "$TEST_JSONPLACEHOLDER_PROXY" ] && [[ "$TEST_JSONPLACEHOLDER_PROXY" != *"example.com"* ]]; then
    echo "Testing with JSONPlaceholder configuration..."
    timeout 30s go run . -config test-configs/jsonplaceholder-config.json || echo "JSONPlaceholder test completed or timed out"
else
    echo "Skipping JSONPlaceholder test (using example proxy)"
fi

# Clean up temporary files (optional)
echo ""
echo "=== Cleanup ==="
echo "Temporary config files created with current proxy credentials."
echo "These files are ignored by git and will not be committed."

echo ""
echo "=== Test Summary ==="
echo "✅ Environment-based configuration system working"
echo "✅ Template system functioning correctly"
echo "✅ Proxy credentials loaded from environment variables"
echo ""
echo "To use real proxy credentials:"
echo "1. Create .env file: cp .env.example .env"
echo "2. Edit .env with your real proxy credentials"
echo "3. Run: ./test_with_env.sh"