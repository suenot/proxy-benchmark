# Metrics Collection

<cite>
**Referenced Files in This Document **   
- [metrics.go](file://metrics.go)
- [benchmark.go](file://benchmark.go)
- [statistics.go](file://statistics.go)
</cite>

## Table of Contents
1. [Metrics Collection](#metrics-collection)
2. [Core Components](#core-components)
3. [Architecture Overview](#architecture-overview)
4. [Detailed Component Analysis](#detailed-component-analysis)
5. [Performance Considerations](#performance-considerations)

## Core Components

The metrics collection mechanism is centered around the `Metrics` struct defined in `metrics.go`, which serves as a thread-safe container for collecting and storing performance data during proxy benchmarking. The system collects three primary types of timing metrics: request durations, ping (TCP connection) times, and derived proxy processing times.

The `Metrics` struct incorporates a `sync.Mutex` field to ensure thread safety when multiple goroutines concurrently update metrics during high-concurrency benchmarking phases. This synchronization mechanism prevents race conditions that could corrupt data integrity when recording measurements from parallel operations.

Three supporting structs—`RequestMetrics`, `PingMetrics`, and `DerivedMetrics`—organize different categories of collected data, while the `Statistics` struct holds calculated statistical values such as minimum, maximum, mean, median, standard deviation, and configurable percentiles.

**Section sources**
- [metrics.go](file://metrics.go#L1-L122)

## Architecture Overview

The metrics collection system operates within a multi-phase benchmarking workflow orchestrated by the `BenchmarkEngine`. During execution, metrics are collected across distinct phases: warmup, ping measurement, and request benchmarking. Each phase contributes specific timing data that ultimately feeds into comprehensive performance analysis.

```mermaid
graph TD
A[BenchmarkEngine.Run] --> B[Initialize Metrics]
B --> C[runWarmup]
C --> D[runPingMeasurement]
D --> E[runRequestBenchmarking]
E --> F[calculateDerivedMetrics]
F --> G[calculateStatistics]
G --> H[Generate Report]
subgraph "Metrics Collection"
M1[NewMetrics]
M2[AddRequestTime]
M3[AddPingTime]
M4[AddDerivedTime]
end
D --> M3
E --> M2
F --> M4
G --> UpdateMetricsStatistics
```

**Diagram sources **
- [benchmark.go](file://benchmark.go#L39-L75)
- [metrics.go](file://metrics.go#L48-L61)

## Detailed Component Analysis

### Metrics Struct Analysis

The `Metrics` struct serves as the central data container for all performance measurements associated with a single proxy. It maintains separate collections for different metric types and ensures thread-safe access through mutex synchronization.

#### Class Diagram of Metrics Structure
```mermaid
classDiagram
class Metrics {
+string ProxyString
+RequestMetrics RequestMetrics
+PingMetrics PingMetrics
+DerivedMetrics DerivedMetrics
-sync.Mutex mu
+AddRequestTime(duration time.Duration, success bool)
+AddPingTime(duration time.Duration)
+AddDerivedTime(duration int64)
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
Metrics --> RequestMetrics : "contains"
Metrics --> PingMetrics : "contains"
Metrics --> DerivedMetrics : "contains"
RequestMetrics --> Statistics : "references"
PingMetrics --> Statistics : "references"
DerivedMetrics --> Statistics : "references"
```

**Diagram sources **
- [metrics.go](file://metrics.go#L10-L45)

### Data Collection Workflow

The metrics system captures performance data at various stages of the benchmarking process, with each measurement type serving a specific analytical purpose.

#### Sequence Diagram of Metric Updates
```mermaid
sequenceDiagram
participant Engine as BenchmarkEngine
participant Metrics as Metrics
participant Client as HTTP/SOCKS Client
Engine->>Metrics : NewMetrics(proxy)
loop For each request
Client->>Engine : MakeRequest()
Engine->>Metrics : AddRequestTime(duration, success)
Note over Metrics : Mutex protects<br/>append operation
end
loop For each ping
Engine->>Metrics : AddPingTime(duration)
Note over Metrics : Mutex protects<br/>append operation
end
Engine->>Metrics : GetRequestTimes()
Engine->>Metrics : GetPingTimes()
Engine->>Metrics : AddDerivedTime(derived)
Note over Metrics : Mutex protects<br/>append operation
```

**Diagram sources **
- [metrics.go](file://metrics.go#L64-L83)
- [benchmark.go](file://benchmark.go#L174-L187)

### Integration with Benchmarking Process

The metrics collection system is tightly integrated with the benchmarking engine, receiving timing data from various phases of execution. The `BenchmarkEngine.Run` method orchestrates the entire process, initializing metrics structures and coordinating data collection across concurrent goroutines.

During the ping measurement phase, direct TCP connection times to proxies are recorded using `AddPingTime`. In the request benchmarking phase, complete HTTP/SOCKS request durations are captured via `AddRequestTime`, including success/failure status. After raw data collection, derived metrics are calculated by subtracting twice the ping time from request times, representing the actual proxy processing overhead.

#### Flowchart of Metrics Lifecycle
```mermaid
flowchart TD
Start([Start Benchmark]) --> Init["Initialize Metrics<br/>for each proxy"]
Init --> Warmup["Run Warmup Phase"]
Warmup --> Ping["Run Ping Measurement"]
Ping --> RecordPing["Metrics.AddPingTime()"]
RecordPing --> Request["Run Request Benchmarking"]
Request --> RecordRequest["Metrics.AddRequestTime()"]
RecordRequest --> Derive["Calculate Derived Metrics"]
Derive --> AddDerived["Metrics.AddDerivedTime()"]
AddDerived --> Stats["Calculate Statistics"]
Stats --> Update["UpdateMetricsStatistics()"]
Update --> Complete([Benchmark Complete])
```

**Diagram sources **
- [benchmark.go](file://benchmark.go#L39-L75)
- [metrics.go](file://metrics.go#L86-L95)

## Performance Considerations

The metrics collection system employs several design patterns to balance accuracy and performance under high-concurrency workloads. The use of `sync.Mutex` ensures data integrity but introduces potential contention points when numerous goroutines simultaneously attempt to update metrics.

Each metric update operation follows a consistent pattern: acquire mutex lock, modify the underlying slice or counter, then release the lock. While this approach guarantees thread safety, frequent updates in high-throughput scenarios may impact overall benchmark performance due to lock contention.

To mitigate these effects, the system minimizes critical section duration by performing only essential operations within locked regions—primarily appending to slices and incrementing counters. More intensive operations like statistical calculations are deferred until after data collection completes, reducing pressure on the mutex during active benchmarking phases.

The design also considers memory efficiency by pre-allocating slices where possible and avoiding unnecessary copying during read operations. When retrieving metric data for statistical analysis, copies are made to prevent blocking write operations, maintaining responsiveness even during report generation.

**Section sources**
- [metrics.go](file://metrics.go#L64-L95)
- [benchmark.go](file://benchmark.go#L240-L255)