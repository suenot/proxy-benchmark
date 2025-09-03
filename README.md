# Proxies Benchmark

A comprehensive benchmarking tool for testing and comparing proxy servers performance. This tool measures latency, throughput, and reliability of HTTP/HTTPS/SOCKS5 proxies through concurrent testing.

## Features

- **Multi-Protocol Support**: Benchmarks HTTP, HTTPS, and SOCKS5 proxies
- **Comprehensive Metrics**: Measures ping time, request time, and derived processing time
- **Concurrent Testing**: Tests multiple proxies simultaneously for efficient benchmarking
- **Statistical Analysis**: Calculates mean, median, and customizable percentiles
- **Warmup Phase**: Ensures stable connections before benchmarking
- **Detailed Reports**: Generates both detailed and summary JSON reports
- **Configurable Parameters**: Customizable timeouts, intervals, and request counts

## Installation

### Prerequisites

- Go 1.19 or higher
- Make (optional, for using Makefile)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/suenot/proxy-benchmark.git
cd proxy-benchmark

# Build using Make
make build

# Or build directly with Go
go build -o proxy-benchmark
```

## Configuration

Create a `config.json` file based on the provided `config.example.json`:

```json
{
  "proxies": [
    "socks:proxy1.example.com:1080:username:password:enabled",
    "http:proxy2.example.com:8080:username:password:enabled",
    "https:proxy3.example.com:443:username:password:enabled"
  ],
  "benchmark": {
    "requests": 100,
    "interval_ms": 5000,
    "warmup_requests": 10,
    "target_url": "https://httpbin.org/get",
    "concurrency": 10,
    "timeout_ms": 30000
  },
  "statistics": {
    "percentiles": [90, 95, 99],
    "mean": true,
    "median": true
  }
}
```

### Configuration Parameters

#### Proxy Format
Each proxy should be specified as a string with the following format:
```
protocol:host:port:username:password:status
```

- **protocol**: `http`, `https`, or `socks` (SOCKS5)
- **host**: Proxy server hostname or IP
- **port**: Proxy server port
- **username**: Authentication username
- **password**: Authentication password
- **status**: `enabled` or `disabled`

#### Benchmark Settings

| Parameter | Description | Default |
|-----------|-------------|---------|
| `requests` | Number of benchmark requests per proxy | 100 |
| `interval_ms` | Delay between requests in milliseconds | 5000 |
| `warmup_requests` | Number of warmup requests before benchmarking | 10 |
| `target_url` | URL to test through proxies | https://httpbin.org/get |
| `concurrency` | Number of concurrent workers | 10 |
| `timeout_ms` | Request timeout in milliseconds | 30000 |

#### Statistics Settings

| Parameter | Description |
|-----------|-------------|
| `percentiles` | Array of percentile values to calculate (e.g., [90, 95, 99]) |
| `mean` | Calculate mean values (true/false) |
| `median` | Calculate median values (true/false) |

## Usage

### Basic Usage

```bash
# Run with default config.json
./proxy-benchmark

# Run with custom configuration file
./proxy-benchmark -config custom-config.json
```

### Output Files

The benchmark generates two output files:

1. **`result.json`**: Detailed benchmark results with all metrics
2. **`results_short.json`**: Condensed summary for quick overview

## Benchmark Algorithm

The benchmarking process follows a sophisticated multi-phase approach:

```mermaid
flowchart TD
    Start([Start Benchmark]) --> LoadConfig[Load Configuration]
    LoadConfig --> ParseProxies[Parse & Validate Proxies]
    ParseProxies --> InitMetrics[Initialize Metrics Storage]
    
    InitMetrics --> WarmupPhase[Warmup Phase]
    WarmupPhase --> WarmupLoop{For Each Proxy}
    WarmupLoop --> WarmupRequests[Send Warmup Requests<br/>Concurrent Execution]
    WarmupRequests --> WarmupNext{More Proxies?}
    WarmupNext -->|Yes| WarmupLoop
    WarmupNext -->|No| PingPhase
    
    PingPhase[Ping Measurement Phase] --> PingLoop{For Each Proxy}
    PingLoop --> PingMeasure[Measure TCP Connection Time<br/>to Proxy Server]
    PingMeasure --> StorePing[Store Ping Metrics]
    StorePing --> PingInterval[Wait Interval]
    PingInterval --> PingNext{More Pings?}
    PingNext -->|Yes| PingMeasure
    PingNext -->|No| PingProxyNext{More Proxies?}
    PingProxyNext -->|Yes| PingLoop
    PingProxyNext -->|No| RequestPhase
    
    RequestPhase[Request Benchmarking Phase] --> RequestLoop{For Each Proxy}
    RequestLoop --> ProxyType{Proxy Protocol?}
    ProxyType -->|HTTP/HTTPS| HTTPClient[Create HTTP Client]
    ProxyType -->|SOCKS5| SOCKSClient[Create SOCKS5 Client]
    
    HTTPClient --> MakeRequest[Make Request to Target URL]
    SOCKSClient --> MakeRequest
    
    MakeRequest --> MeasureTime[Measure Total Request Time]
    MeasureTime --> StoreMetrics[Store Request Metrics<br/>Success/Failure Status]
    StoreMetrics --> RequestInterval[Wait Interval]
    RequestInterval --> RequestNext{More Requests?}
    RequestNext -->|Yes| MakeRequest
    RequestNext -->|No| RequestProxyNext{More Proxies?}
    RequestProxyNext -->|Yes| RequestLoop
    RequestProxyNext -->|No| DerivedMetrics
    
    DerivedMetrics[Calculate Derived Metrics] --> DerivedCalc[Request Time - 2�Ping Time<br/>= Processing Time]
    DerivedCalc --> Statistics[Calculate Statistics]
    
    Statistics --> CalcStats[Calculate for Each Metric:<br/>" Mean<br/>" Median<br/>" Percentiles<br/>" Min/Max<br/>" Success Rate]
    
    CalcStats --> GenerateReport[Generate Reports]
    GenerateReport --> SaveFull[Save result.json<br/>Full Metrics]
    SaveFull --> SaveShort[Save results_short.json<br/>Summary Statistics]
    SaveShort --> End([Benchmark Complete])
    
    style Start fill:#e1f5fe
    style End fill:#c8e6c9
    style WarmupPhase fill:#fff3e0
    style PingPhase fill:#fce4ec
    style RequestPhase fill:#f3e5f5
    style DerivedMetrics fill:#e8f5e9
    style Statistics fill:#e0f2f1
```

### Algorithm Phases Explained

1. **Configuration Loading**: Reads and validates the configuration file
2. **Proxy Parsing**: Validates proxy formats and filters enabled proxies
3. **Warmup Phase**: Establishes initial connections to ensure stable performance
4. **Ping Measurement**: Measures raw TCP connection time to proxy servers
5. **Request Benchmarking**: Performs actual HTTP/HTTPS requests through proxies
6. **Derived Metrics**: Calculates processing time by subtracting network latency
7. **Statistical Analysis**: Computes comprehensive statistics for all metrics

## Metrics Collected

### Primary Metrics

- **Ping Time**: Direct TCP connection time to the proxy server
- **Request Time**: Total time for a request through the proxy
- **Derived Time**: Estimated processing time (Request Time - 2�Ping Time)
- **Success Rate**: Percentage of successful requests

### Statistical Calculations

For each metric, the tool calculates:
- Minimum and Maximum values
- Mean (average)
- Median
- Custom percentiles (e.g., P90, P95, P99)
- Standard deviation
- Success/failure counts

## Example Output

### results_short.json
```json
{
  "socks:proxy1.example.com:1080:user:pass:enabled": {
    "ping_ms": {
      "mean": 45.2,
      "median": 44.0,
      "p90": 52.0,
      "p95": 55.0,
      "p99": 62.0
    },
    "request_ms": {
      "mean": 234.5,
      "median": 220.0,
      "p90": 280.0,
      "p95": 310.0,
      "p99": 380.0
    },
    "success_rate": 98.5
  }
}
```

## Development

### Project Structure

```
main.go              # Entry point
config.go            # Configuration structures and loading
benchmark.go         # Core benchmarking engine
proxy.go             # Proxy parsing and management
http_client.go       # HTTP/HTTPS proxy client
socks5_client.go     # SOCKS5 proxy client
ping.go              # TCP ping implementation
metrics.go           # Metrics collection and storage
statistics.go        # Statistical calculations
reporter.go          # Report generation
config.example.json  # Configuration template
```

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestMetrics
```

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o proxy-benchmark-linux

# Windows
GOOS=windows GOARCH=amd64 go build -o proxy-benchmark.exe

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o proxy-benchmark-mac

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o proxy-benchmark-mac-arm64
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Guidelines

1. Follow Go best practices and conventions
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure all tests pass before submitting PR
5. Use meaningful commit messages

## Use Cases

- **Proxy Provider Comparison**: Compare performance across different proxy providers
- **Network Monitoring**: Monitor proxy performance over time
- **Load Testing**: Test proxy behavior under various load conditions
- **Geographic Performance**: Measure proxy performance from different locations
- **Service Selection**: Choose optimal proxies for specific use cases

## Troubleshooting

### Common Issues

1. **Connection Timeouts**: Increase `timeout_ms` in configuration
2. **Authentication Failures**: Verify proxy credentials are correct
3. **High Failure Rate**: Check proxy server status and network connectivity
4. **Memory Usage**: Reduce `concurrency` for large proxy lists

### Debug Mode

Enable verbose output by modifying the code to add debug logging:
```go
// Add to main.go for debug output
fmt.Printf("Debug: %v\n", debugInfo)
```

## License

This project is open source and available under the [MIT License](LICENSE).

## Author

**suenot** - [GitHub Profile](https://github.com/suenot)

## Acknowledgments

- HTTP testing endpoint provided by [httpbin.org](https://httpbin.org)
- SOCKS5 implementation inspired by Go's x/net/proxy package
- Statistical algorithms based on standard mathematical formulas

## Support

For issues, questions, or suggestions, please [open an issue](https://github.com/suenot/proxy-benchmark/issues) on GitHub.