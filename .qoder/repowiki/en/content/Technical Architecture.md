# Technical Architecture

<cite>
**Referenced Files in This Document **   
- [main.go](file://main.go)
- [benchmark.go](file://benchmark.go)
- [config.go](file://config.go)
- [metrics.go](file://metrics.go)
- [reporter.go](file://reporter.go)
- [http_client.go](file://http_client.go)
- [socks5_client.go](file://socks5_client.go)
- [ping.go](file://ping.go)
- [statistics.go](file://statistics.go)
</cite>

## Table of Contents
1. [Introduction](#introduction)
2. [Project Structure](#project-structure)
3. [Core Components](#core-components)
4. [Architecture Overview](#architecture-overview)
5. [Detailed Component Analysis](#detailed-component-analysis)
6. [Dependency Analysis](#dependency-analysis)
7. [Performance Considerations](#performance-considerations)
8. [Troubleshooting Guide](#troubleshooting-guide)
9. [Conclusion](#conclusion)

## Introduction
The proxy-benchmark application is designed to evaluate the performance of multiple proxy servers through a structured benchmarking process. It follows a modular monolith architecture with clearly separated concerns across configuration, benchmarking, metrics collection, and reporting layers. The system leverages established design patterns such as Factory, Strategy, Observer, and Dependency Injection to ensure extensibility, testability, and maintainability. Concurrency is managed using goroutines per proxy and mutex-protected shared state, enabling efficient parallel execution while preserving data integrity.

## Project Structure

```mermaid
flowchart TD
A[main.go] --> B[benchmark.go]
B --> C[config.go]
B --> D[metrics.go]
B --> E[reporter.go]
B --> F[http_client.go]
B --> G[socks5_client.go]
B --> H[ping.go]
D --> I[statistics.go]
```

**Diagram sources**
- [main.go](file://main.go#L9-L80)
- [benchmark.go](file://benchmark.go#L10-L15)

**Section sources**
- [main.go](file://main.go#L9-L80)
- [benchmark.go](file://benchmark.go#L10-L15)

## Core Components

The core components of the proxy-benchmark system are organized around four primary architectural layers: Configuration, Benchmarking, Metrics, and Reporting. These layers interact in a unidirectional flow from initialization to result generation. The `BenchmarkEngine` orchestrates the entire process by coordinating client creation, executing benchmark phases, collecting metrics, and preparing results for reporting. Each component receives its dependencies explicitly via constructor injection, promoting loose coupling and ease of testing.

**Section sources**
- [benchmark.go](file://benchmark.go#L10-L15)
- [config.go](file://config.go#L8-L12)
- [metrics.go](file://metrics.go#L8-L14)
- [reporter.go](file://reporter.go#L29-L29)

## Architecture Overview

```mermaid
graph TB
subgraph "Configuration Layer"
Config[Config]
end
subgraph "Benchmarking Layer"
Engine[BenchmarkEngine]
HTTPClient[HTTPClient]
SOCKS5Client[SOCKS5Client]
PingClient[PingClient]
end
subgraph "Metrics Layer"
Metrics[Metrics]
Statistics[Statistics]
end
subgraph "Reporting Layer"
Reporter[Reporter]
end
Config --> Engine
Engine --> HTTPClient
Engine --> SOCKS5Client
Engine --> PingClient
Engine --> Metrics
Metrics --> Statistics
Engine --> Reporter
Reporter --> Output[(JSON Reports)]
style Engine stroke:#f66,stroke-width:2px
style Metrics stroke:#66f,stroke-width:2px
style Reporter stroke:#6f6,stroke-width:2px
```

**Diagram sources**
- [benchmark.go](file://benchmark.go#L10-L15)
- [config.go](file://config.go#L8-L12)
- [metrics.go](file://metrics.go#L8-L14)
- [reporter.go](file://reporter.go#L29-L29)

## Detailed Component Analysis

### Benchmark Engine Analysis

The `BenchmarkEngine` serves as the central orchestrator of the benchmarking workflow. It implements the Factory Pattern through `NewBenchmarkEngine`, which parses and validates proxy configurations before instantiation. The engine manages three distinct benchmarking phases—warmup, ping measurement, and request benchmarking—executed sequentially. During each phase, it applies the Strategy Pattern by conditionally selecting the appropriate client implementation based on the proxy protocol (HTTP/HTTPS vs SOCKS5).

```mermaid
sequenceDiagram
participant Main as main.go
participant Engine as BenchmarkEngine
participant Proxy as Proxy
participant Client as Client Interface
participant Metrics as Metrics
participant Reporter as Reporter
Main->>Engine : Run()
Engine->>Engine : Initialize metrics
loop For each proxy
Engine->>Engine : Goroutine per proxy
Engine->>Client : Select strategy by protocol
Client->>Proxy : MakeRequest or Ping
Client-->>Metrics : Add timing data
end
Engine->>Metrics : Calculate derived metrics
Engine->>Metrics : Update statistics
Engine-->>Main : GetResults()
Main->>Reporter : GenerateReport(metrics)
Reporter-->>Main : Return JSON report
```

**Diagram sources**
- [benchmark.go](file://benchmark.go#L39-L75)
- [http_client.go](file://http_client.go#L17-L36)
- [socks5_client.go](file://socks5_client.go#L16-L40)
- [ping.go](file://ping.go#L15-L19)

**Section sources**
- [benchmark.go](file://benchmark.go#L39-L75)
- [benchmark.go](file://benchmark.go#L78-L91)
- [benchmark.go](file://benchmark.go#L131-L144)
- [benchmark.go](file://benchmark.go#L174-L187)

### Metrics System Analysis

The metrics subsystem employs the Observer Pattern during benchmark phases, where time measurements are observed and recorded asynchronously by each `Metrics` instance. Each proxy has its own `Metrics` object that collects raw timings for requests and pings, then derives processing times by subtracting round-trip network latency. All write operations on metrics are protected by a mutex (`mu`) to ensure thread safety during concurrent access from multiple goroutines.

```mermaid
classDiagram
class Metrics {
+string ProxyString
+RequestMetrics RequestMetrics
+PingMetrics PingMetrics
+DerivedMetrics DerivedMetrics
-sync.Mutex mu
+AddRequestTime(duration, success)
+AddPingTime(duration)
+AddDerivedTime(duration)
+GetRequestTimes() []int64
+GetPingTimes() []int64
+GetDerivedTimes() []int64
}
class RequestMetrics {
+int Total
+int Successful
+int Failed
+[]int64 Times
+*Statistics Statistics
}
class PingMetrics {
+[]int64 Times
+*Statistics Statistics
}
class DerivedMetrics {
+[]int64 ProcessingTimes
+*Statistics Statistics
}
class Statistics {
+int64 Min
+int64 Max
+float64 Mean
+float64 Median
+float64 StdDev
+map[string]float64 Percentiles
}
Metrics --> RequestMetrics
Metrics --> PingMetrics
Metrics --> DerivedMetrics
RequestMetrics --> Statistics
PingMetrics --> Statistics
DerivedMetrics --> Statistics
```

**Diagram sources**
- [metrics.go](file://metrics.go#L8-L14)
- [metrics.go](file://metrics.go#L48-L61)

**Section sources**
- [metrics.go](file://metrics.go#L8-L122)

### Reporting Layer Analysis

The reporting module uses the Factory Pattern via `NewReporter()` to instantiate a reporter responsible for transforming raw metrics into structured JSON reports. Two output formats are supported: a full `BenchmarkResult` containing detailed statistics per proxy, and a concise `ShortSummary` that includes only mean values. This dual-reporting approach enables both deep analysis and quick comparisons.

```mermaid
flowchart TD
Start([Generate Report]) --> CheckMetrics{"Metrics Available?"}
CheckMetrics --> |Yes| ExtractData["Extract Proxy Metrics"]
ExtractData --> FormatFull["Format Full BenchmarkResult"]
ExtractData --> FormatShort["Format ShortSummary with Means"]
FormatFull --> SaveJSON["Save to result.json"]
FormatShort --> SaveShort["Save to results_short.json"]
SaveJSON --> End1([Report Generated])
SaveShort --> End2([Summary Saved])
style Start fill:#f9f,stroke:#333
style End1 fill:#bbf,stroke:#333,color:#fff
style End2 fill:#bbf,stroke:#333,color:#fff
```

**Diagram sources**
- [reporter.go](file://reporter.go#L37-L54)
- [reporter.go](file://reporter.go#L57-L72)
- [reporter.go](file://reporter.go#L75-L85)

**Section sources**
- [reporter.go](file://reporter.go#L37-L85)

## Dependency Analysis

```mermaid
graph LR
A[main.go] --> B[BenchmarkEngine]
B --> C[Config]
B --> D[Metrics]
B --> E[HTTPClient]
B --> F[SOCKS5Client]
B --> G[PingClient]
D --> H[Statistics]
B --> I[Reporter]
I --> J[JSON Encoding]
style A fill:#f96,stroke:#333
style B fill:#69f,stroke:#333,color:#fff
style C fill:#6f9,stroke:#333
style D fill:#6cf,stroke:#333
style I fill:#6c6,stroke:#333
```

**Diagram sources**
- [main.go](file://main.go#L9-L80)
- [benchmark.go](file://benchmark.go#L10-L15)
- [config.go](file://config.go#L8-L12)
- [metrics.go](file://metrics.go#L8-L14)
- [reporter.go](file://reporter.go#L29-L29)

**Section sources**
- [main.go](file://main.go#L9-L80)
- [benchmark.go](file://benchmark.go#L10-L15)

## Performance Considerations

The concurrency model utilizes one goroutine per proxy during each benchmark phase, allowing parallel execution without blocking. Shared state in `Metrics` is protected by mutexes rather than channels, favoring simplicity and direct control over communication complexity. While this choice increases contention risk under high concurrency, the current design assumes moderate numbers of proxies and intervals between requests, making mutex overhead acceptable. The use of Go's native sync primitives aligns with the project’s goal of readability and maintainability over maximum throughput.

The decision to use JSON for configuration and reporting enhances interoperability and human-readability, facilitating integration with external tools and manual inspection. Although binary formats could reduce I/O size, JSON was selected for its ubiquity and ease of debugging in operational contexts.

**Section sources**
- [benchmark.go](file://benchmark.go#L78-L91)
- [metrics.go](file://metrics.go#L13-L13)
- [main.go](file://main.go#L9-L80)

## Troubleshooting Guide

Common issues typically arise from invalid proxy configurations, network timeouts, or unsupported protocols. The application logs warnings when parsing invalid proxies and skips them gracefully. Timeouts are configurable via `timeout_ms` in the config file, and increasing this value may resolve transient connectivity issues. Protocol support is limited to "http", "https", and "socks"—any other values will trigger an unsupported protocol warning.

When benchmark results show consistently high latencies or failures, verify:
- Proxy credentials and addresses
- Network reachability to proxy endpoints
- Target URL accessibility
- Sufficient timeout settings relative to network conditions

Logs are printed directly to stdout with contextual messages indicating phase progression and errors, aiding real-time diagnosis.

**Section sources**
- [benchmark.go](file://benchmark.go#L39-L75)
- [proxy.go](file://proxy.go#L18-L32)
- [main.go](file://main.go#L9-L80)

## Conclusion

The proxy-benchmark system exemplifies a well-structured modular monolith with clear separation of concerns and thoughtful application of design patterns. Its architecture balances performance, correctness, and maintainability through strategic use of concurrency, dependency injection, and layered abstraction. By adhering to Go idioms and prioritizing code clarity, the system remains accessible for extension and adaptation to new benchmarking requirements.