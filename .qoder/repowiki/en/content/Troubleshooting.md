# Troubleshooting

<cite>
**Referenced Files in This Document**   
- [main.go](file://main.go)
- [benchmark.go](file://benchmark.go)
- [config.go](file://config.go)
- [http_client.go](file://http_client.go)
- [socks5_client.go](file://socks5_client.go)
- [ping.go](file://ping.go)
- [proxy.go](file://proxy.go)
- [config.json](file://config.json)
</cite>

## Table of Contents
1. [Connection Timeouts](#connection-timeouts)  
2. [Authentication Failures](#authentication-failures)  
3. [Malformed Proxy Strings](#malformed-proxy-strings)  
4. [Empty Result Files](#empty-result-files)  
5. [High-Latency Proxies and Skewed Results](#high-latency-proxies-and-skewed-results)

## Connection Timeouts

Connection timeouts occur when the benchmark tool fails to establish a TCP connection or complete an HTTP request within the configured time limit. These are typically reported during warmup, ping, or request phases with messages such as `"failed to connect to proxy"` or `"Request failed for proxy"`.

### Root Causes
- **Insufficient timeout settings**: The `timeout_ms` value in `config.json` may be too low for slow proxies.
- **Network connectivity issues**: The target proxy server might be unreachable due to firewall rules, network outages, or DNS resolution problems.
- **Dialer configuration**: In `ping.go`, the `net.Dialer` uses the timeout from `PingClient.timeout`, which is derived from `config.Benchmark.TimeoutMs`. If this is too short, even valid proxies may fail.

### Diagnostic Steps
1. Check if the proxy host is reachable using external tools like `telnet` or `nc`.
2. Verify that no local firewall or corporate network policy blocks outbound connections to the proxy port.
3. Enable verbose logging by adding print statements in `runWarmupForProxy()` or `runRequestBenchmarkingForProxy()` to trace failure points.

### Solutions
- Increase the `timeout_ms` value in `config.json` (default: 5000 ms).
- Ensure the target URL (`target_url`) is accessible without the proxy.
- Retry with a known-working proxy to isolate whether the issue is systemic or proxy-specific.

**Section sources**
- [ping.go](file://ping.go#L22-L40)
- [benchmark.go](file://benchmark.go#L94-L128)
- [benchmark.go](file://benchmark.go#L190-L237)
- [config.json](file://config.json)

## Authentication Failures

Authentication failures arise when credentials provided in the proxy string are rejected by the proxy server. These appear as errors in client creation, such as `"Failed to create HTTP client for proxy"` or `"dial tcp: connection refused"` during SOCKS5 handshake.

### Root Causes
- Incorrect username or password in the proxy string.
- Malformed proxy URL construction in `NewHTTPClient()` where credentials are embedded directly into the `http://user:pass@host:port` format.
- SOCKS5 authentication rejection in `NewSOCKS5Client()`, which uses `golang.org/x/net/proxy.Auth`.

### Diagnostic Steps
1. Validate credentials independently using tools like `curl`:
   ```bash
   curl --proxy http://username:password@host:port https://httpbin.org/get
   ```
2. Confirm that the proxy supports the specified protocol (`http`, `https`, `socks`).

### Solutions
- Correct the username and password fields in the proxy string within `config.json`.
- Ensure special characters in credentials are properly URL-encoded.
- Test with a different credential set or proxy provider to confirm validity.

**Section sources**
- [http_client.go](file://http_client.go#L17-L36)
- [socks5_client.go](file://socks5_client.go#L16-L40)
- [config.json](file://config.json)

## Malformed Proxy Strings

Malformed proxy strings cause parsing errors during initialization, leading to skipped proxies or runtime panics.

### Root Causes
- Invalid format in the `proxies` array: each entry must follow the exact six-part structure:
  ```
  protocol:host:port:username:password:status
  ```
- Missing or extra colons (`:`) causing `ParseProxy()` in `proxy.go` to return an error.
- Disabled status (`"disabled"`) preventing usage even if structurally valid.

### Diagnostic Steps
1. Review all entries in the `proxies` list in `config.json`.
2. Look for warnings like `"Warning: skipping invalid proxy"` in the console output from `NewBenchmarkEngine()`.

### Solutions
- Format all proxy strings strictly according to the required schema.
- Use `config.example.json` as a template for correct syntax.
- Ensure every proxy ends with `:enabled` to be included in testing.

**Section sources**
- [proxy.go](file://proxy.go#L8-L15)
- [proxy.go](file://proxy.go#L17-L48)
- [config.json](file://config.json)
- [config.example.json](file://config.example.json)

## Empty Result Files

Empty result files (`result.json`, `results_short.json`) indicate that no benchmarks completed successfully, often due to early-stage failures.

### Root Causes
- Configuration file not found or unreadable — handled in `main.go` with `os.Stat()` check.
- No valid proxies parsed — if all proxies are malformed or disabled, `BenchmarkEngine` initializes with zero proxies.
- Early termination in `Run()` due to warmup, ping, or request phase returning an error.

### Diagnostic Steps
1. Confirm `config.json` exists at the specified path or use the `-config` flag correctly.
2. Check console output for fatal errors like `"Configuration file not found"` or `"Failed to load configuration"`.
3. Verify that at least one proxy has `:enabled` status and correct formatting.

### Solutions
- Provide a valid configuration file path via command-line argument:  
  ```bash
  go run main.go -config test-config.json
  ```
- Ensure at least one properly formatted and enabled proxy is present.
- Add debug prints in `LoadConfig()` to verify JSON decoding success.

**Section sources**
- [main.go](file://main.go#L16-L80)
- [config.go](file://config.go#L32-L47)
- [benchmark.go](file://benchmark.go#L18)
- [config.json](file://config.json)

## High-Latency Proxies and Skewed Results

High-latency proxies can distort benchmark outcomes, especially when measuring derived metrics like processing time (request time minus round-trip ping × 2).

### Root Causes
- Network congestion or geographical distance increasing ping times.
- Overloaded proxy servers introducing variable delays.
- Inconsistent interval timing between requests due to fixed `interval_ms` regardless of response time.

### Diagnostic Steps
1. Examine `ping_metrics` in the output JSON to identify outliers.
2. Compare `request_metrics` against `derived_metrics` to detect excessive processing overhead.
3. Run multiple iterations to assess consistency across runs.

### Solutions
- Adjust `interval_ms` to allow sufficient cooling between requests.
- Filter out high-ping proxies before final analysis using post-processing scripts.
- Increase sample size (`requests`) to improve statistical significance despite noise.

```mermaid
flowchart TD
A["Start Benchmark"] --> B{Valid Config?}
B --> |No| C["Exit with Error"]
B --> |Yes| D["Parse Proxies"]
D --> E{Any Valid Enabled Proxies?}
E --> |No| F["Skip All Proxies"]
E --> |Yes| G["Run Warmup Phase"]
G --> H["Run Ping Measurement"]
H --> I["Run Request Benchmarking"]
I --> J["Generate Report"]
J --> K["Save result.json & results_short.json"]
```

**Diagram sources**
- [main.go](file://main.go#L16-L80)
- [benchmark.go](file://benchmark.go#L39-L62)
- [reporter.go](file://reporter.go#L45-L99)

**Section sources**
- [benchmark.go](file://benchmark.go#L147-L171)
- [metrics.go](file://metrics.go#L80-L100)
- [statistics.go](file://statistics.go)
- [config.json](file://config.json)