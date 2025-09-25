# Request Benchmarking

<cite>
**Referenced Files in This Document**   
- [benchmark.go](file://benchmark.go)
- [http_client.go](file://http_client.go)
- [socks5_client.go](file://socks5_client.go)
- [metrics.go](file://metrics.go)
- [config.go](file://config.go)
- [proxy.go](file://proxy.go)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [Request Execution Workflow](#request-execution-workflow)
3. [Protocol-Specific Client Implementation](#protocol-specific-client-implementation)
4. [Concurrency and Parallelization Strategy](#concurrency-and-parallelization-strategy)
5. [Metrics Collection and Atomic Updates](#metrics-collection-and-atomic-updates)
6. [Error Handling and Recovery Mechanisms](#error-handling-and-recovery-mechanisms)
7. [Timeout and Interval Configuration](#timeout-and-interval-configuration)
8. [Challenges in Benchmark Accuracy](#challenges-in-benchmark-accuracy)

## Introduction
The request benchmarking phase is a core component of the proxy performance evaluation system, designed to measure end-to-end response times for HTTP/HTTPS and SOCKS5 proxy types under controlled conditions. This process is orchestrated by the `BenchmarkEngine` in `benchmark.go`, which coordinates parallel execution across multiple proxies using protocol-specific clients implemented in `http_client.go` and `socks5_client.go`. The system measures full round-trip durations while handling errors, tracking success rates, and updating metrics safely across concurrent goroutines. Each proxy undergoes warmup, ping measurement, and actual request benchmarking phases before derived statistics are calculated.

**Section sources**
- [benchmark.go](file://benchmark.go#L39-L75)

## Request Execution Workflow
The request benchmarking workflow begins with the `runRequestBenchmarking` method, which launches a separate goroutine for each proxy via `runRequestBenchmarkingForProxy`. For every configured request iteration, the system enforces a time interval between requests (except the first), creates a context with timeout constraints, and records the start time before initiating the request. After completion, the elapsed duration is computed and recorded along with success status. This loop ensures consistent pacing and accurate timing measurements across all proxies.

```mermaid
flowchart TD
Start([Start Benchmark]) --> Init["Initialize Metrics for Proxies"]
Init --> Warmup["Run Warmup Phase"]
Warmup --> Ping["Run Ping Measurement Phase"]
Ping --> Request["Run Request Benchmarking Phase"]
Request --> Derive["Calculate Derived Metrics"]
Derive --> Stats["Calculate Statistics"]
Stats --> Complete([Benchmark Completed])
subgraph "Per-Proxy Request Loop"
RequestLoop["For each request:"] --> SleepCheck{"i > 0?"}
SleepCheck --> |Yes| Sleep["time.Sleep(interval)"]
Sleep --> Context["Create Context with Timeout"]
SleepCheck --> |No| Context
Context --> StartTime["Record Start Time"]
StartTime --> ProtocolCheck{"Protocol Type"}
ProtocolCheck --> |HTTP/HTTPS| CreateHTTP["NewHTTPClient()"]
ProtocolCheck --> |SOCKS5| CreateSOCKS5["NewSOCKS5Client()"]
CreateHTTP --> MakeRequest["client.MakeRequest()"]
CreateSOCKS5 --> MakeRequest
MakeRequest --> EndTime["Record End Time"]
EndTime --> Record["Add Duration to Metrics"]
Record --> Next[Next Request]
end
```

**Diagram sources**
- [benchmark.go](file://benchmark.go#L190-L237)
- [benchmark.go](file://benchmark.go#L174-L187)

**Section sources**
- [benchmark.go](file://benchmark.go#L190-L237)
- [benchmark.go](file://benchmark.go#L174-L187)

## Protocol-Specific Client Implementation
The system implements two distinct client types to handle different proxy protocols: `HTTPClient` for HTTP/HTTPS proxies and `SOCKS5Client` for SOCKS5 proxies. Both conform to a common interface through their `MakeRequest` methods but differ fundamentally in how they route traffic.

### HTTP/HTTPS Client Implementation
The `HTTPClient` uses Go's standard `http.Transport` with a proxy URL that includes authentication credentials embedded in the format `http://user:pass@host:port`. It configures an `http.Client` with this transport and applies a global timeout. The `MakeRequest` method leverages `http.NewRequestWithContext` to ensure cancellation and timeout propagation during the request lifecycle.

```mermaid
classDiagram
class HTTPClient {
+*http.Client client
+time.Duration timeout
+MakeRequest(ctx, url) (time.Duration, error)
}
HTTPClient --> http.Client : "uses"
http.Client --> http.Transport : "uses"
http.Transport --> "ProxyURL" : "routes via"
```

**Diagram sources**
- [http_client.go](file://http_client.go#L11-L14)
- [http_client.go](file://http_client.go#L17-L36)

### SOCKS5 Client Implementation
The `SOCKS5Client` relies on the `golang.org/x/net/proxy` package to create a custom dialer using `proxy.SOCKS5`, which handles authentication via `proxy.Auth`. Instead of setting a proxy URL, it replaces the `Dial` function in `http.Transport` with the SOCKS5 dialer’s `Dial` method, enabling direct TCP-level proxying. This approach allows standard HTTP clients to operate over SOCKS5 tunnels transparently.

```mermaid
classDiagram
class SOCKS5Client {
+*http.Client client
+time.Duration timeout
+MakeRequest(ctx, url) (time.Duration, error)
}
SOCKS5Client --> http.Client : "uses"
http.Client --> http.Transport : "uses"
http.Transport --> "Dial" : "replaced with"
"Dial" --> SOCKS5Dialer : "delegates to"
SOCKS5Dialer --> proxy.SOCKS5 : "created from"
```

**Diagram sources**
- [socks5_client.go](file://socks5_client.go#L10-L13)
- [socks5_client.go](file://socks5_client.go#L16-L40)

**Section sources**
- [http_client.go](file://http_client.go#L17-L36)
- [socks5_client.go](file://socks5_client.go#L16-L40)

## Concurrency and Parallelization Strategy
To maximize efficiency and simulate real-world usage patterns, the benchmark engine executes requests for each proxy in parallel using goroutines coordinated by `sync.WaitGroup`. The `runRequestBenchmarking` method spawns one goroutine per proxy, ensuring isolation and preventing bottlenecks. Within each goroutine, requests are executed sequentially with configurable intervals to avoid overwhelming target servers or triggering rate limits.

```mermaid
sequenceDiagram
participant Engine as BenchmarkEngine
participant Proxy1 as Proxy Worker 1
participant Proxy2 as Proxy Worker 2
participant WaitGroup as sync.WaitGroup
Engine->>WaitGroup : Add(number of proxies)
Engine->>Proxy1 : go runRequestBenchmarkingForProxy()
Engine->>Proxy2 : go runRequestBenchmarkingForProxy()
Proxy1->>WaitGroup : defer Done()
Proxy2->>WaitGroup : defer Done()
Engine->>WaitGroup : Wait()
WaitGroup-->>Engine : All workers complete
```

This design enables true concurrency without race conditions due to isolated metric storage per proxy and atomic updates protected by mutexes within the `Metrics` structure.

**Diagram sources**
- [benchmark.go](file://benchmark.go#L174-L187)
- [benchmark.go](file://benchmark.go#L190-L237)

**Section sources**
- [benchmark.go](file://benchmark.go#L174-L187)

## Metrics Collection and Atomic Updates
All performance data is captured in the `Metrics` struct, which maintains separate time series for request durations, ping times, and derived processing times. Each update operation—such as `AddRequestTime` or `AddPingTime`—is protected by a `sync.Mutex` to prevent data races when multiple goroutines attempt to write simultaneously. Successful responses increment counters and store millisecond-resolution durations, while failed attempts record zero or measured partial durations.

```mermaid
classDiagram
class Metrics {
+string ProxyString
+RequestMetrics RequestMetrics
+PingMetrics PingMetrics
+DerivedMetrics DerivedMetrics
+sync.Mutex mu
+AddRequestTime(duration, success)
+AddPingTime(duration)
+AddDerivedTime(duration)
}
class RequestMetrics {
+int Total
+int Successful
+int Failed
+[]int64 Times
}
class PingMetrics {
+[]int64 Times
}
class DerivedMetrics {
+[]int64 ProcessingTimes
}
Metrics --> RequestMetrics
Metrics --> PingMetrics
Metrics --> DerivedMetrics
```

After collecting raw timings, the system calculates derived metrics such as estimated backend processing time by subtracting twice the ping duration (round-trip to proxy) from total request time.

**Diagram sources**
- [metrics.go](file://metrics.go#L8-L14)
- [metrics.go](file://metrics.go#L17-L23)

**Section sources**
- [metrics.go](file://metrics.go#L64-L75)
- [benchmark.go](file://benchmark.go#L239-L250)

## Error Handling and Recovery Mechanisms
The system incorporates robust error handling at multiple levels. During client creation, failures due to malformed proxy URLs or invalid configurations result in immediate logging and skipping of the request, with a zero-duration failure recorded. During request execution, any network-level or protocol-level error (e.g., TLS handshake failure, connection reset) is caught, logged with the proxy address, and counted as a failed attempt. No retries are performed to preserve measurement integrity, but warmup requests help mitigate cold-start effects.

Failed pings are treated similarly, with zero values stored to maintain alignment between request and ping time arrays for later derivation calculations. This ensures statistical consistency even in the presence of intermittent connectivity issues.

**Section sources**
- [benchmark.go](file://benchmark.go#L190-L237)
- [benchmark.go](file://benchmark.go#L147-L171)

## Timeout and Interval Configuration
All operations respect configurable timeouts defined in `config.Benchmark.TimeoutMs`, converted into `time.Duration` for use with `context.WithTimeout`. This prevents individual requests from hanging indefinitely and ensures predictable execution times. Additionally, `IntervalMs` governs the delay between successive requests to the same proxy, simulating realistic load patterns and avoiding rate limiting.

These values are passed directly to both `NewHTTPClient` and `NewSOCKS5Client`, where they configure the underlying `http.Client`'s timeout behavior. Context cancellation (`defer cancel()`) ensures resources are released promptly upon completion or timeout.

**Section sources**
- [benchmark.go](file://benchmark.go#L190-L237)
- [config.go](file://config.go#L20-L25)

## Challenges in Benchmark Accuracy
Several factors can impact benchmark accuracy:

- **TLS Handshake Failures**: Common with misconfigured HTTPS proxies; mitigated by warmup phase but still contribute to failure counts.
- **Proxy Authentication Requirements**: Incorrect credentials lead to immediate rejection; the current implementation assumes correct credentials in input.
- **Rate Limiting**: Without adaptive backoff, repeated requests may trigger throttling, inflating latency measurements.
- **Network Variability**: External network conditions affect ping and request times independently of proxy performance.
- **Clock Synchronization**: High-precision timing assumes monotonic clocks; virtualized environments may introduce jitter.

Future improvements could include retry logic with exponential backoff, dynamic adjustment of request intervals, and validation of proxy reachability prior to benchmarking.

**Section sources**
- [benchmark.go](file://benchmark.go#L190-L237)
- [ping.go](file://ping.go#L25-L40)
- [http_client.go](file://http_client.go#L39-L55)
- [socks5_client.go](file://socks5_client.go#L43-L59)