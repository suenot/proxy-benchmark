#!/bin/bash

echo "=== Running Validation Tests ==="
echo ""

echo "1. Unit Tests:"
go test -v -run TestValidateResponse
echo ""

echo "2. Benchmark Tests:"
go test -bench=BenchmarkValidateResponse -benchmem
echo ""

echo "3. Testing with JSONPlaceholder API (valid config):"
echo "   Testing endpoint: https://jsonplaceholder.typicode.com/posts/1"
timeout 5 ./proxy-benchmark test-configs/jsonplaceholder-config.json 2>&1 | grep -E "benchmarking|Calculating|completed" | head -5
echo ""

echo "4. Testing with invalid validation rules (should fail):"
echo "   This config has a non-existent field 'nonExistentField'"
timeout 5 ./proxy-benchmark test-configs/jsonplaceholder-fail-config.json 2>&1 | grep -i "validation\|failed" | head -5
if [ $? -eq 0 ]; then
    echo "   ✓ Validation errors detected as expected"
else
    echo "   ⚠ No validation errors found (check might have issues)"
fi
echo ""

echo "5. Testing with GitHub API (if proxy allows):"
echo "   Testing endpoint: https://api.github.com/users/octocat"
timeout 5 ./proxy-benchmark test-configs/github-api-config.json 2>&1 | grep -E "benchmarking|Calculating|completed" | head -5
echo ""

echo "=== Tests Complete ==="