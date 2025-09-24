package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

// BenchmarkEngine orchestrates the benchmarking process
type BenchmarkEngine struct {
	config  *Config
	proxies []*Proxy
	metrics map[string]*Metrics
	mu      sync.Mutex
}

// NewBenchmarkEngine creates a new benchmark engine
func NewBenchmarkEngine(config *Config) (*BenchmarkEngine, error) {
	proxies := make([]*Proxy, 0)
	for _, proxyStr := range config.Proxies {
		proxy, err := ParseProxy(proxyStr)
		if err != nil {
			fmt.Printf("Warning: skipping invalid proxy %s: %v\n", proxyStr, err)
			continue
		}
		if proxy.IsValid() {
			proxies = append(proxies, proxy)
		}
	}

	return &BenchmarkEngine{
		config:  config,
		proxies: proxies,
		metrics: make(map[string]*Metrics),
	}, nil
}

// Run executes the complete benchmark process
func (b *BenchmarkEngine) Run() error {
	fmt.Println("Starting proxy benchmark...")

	// Initialize metrics for each proxy
	for _, proxy := range b.proxies {
		b.metrics[proxy.String()] = NewMetrics(proxy.String())
	}

	// Run warmup phase
	fmt.Println("Running warmup phase...")
	if err := b.runWarmup(); err != nil {
		return fmt.Errorf("warmup phase failed: %w", err)
	}

	// Run ping measurement phase
	fmt.Println("Running ping measurement phase...")
	if err := b.runPingMeasurement(); err != nil {
		return fmt.Errorf("ping measurement phase failed: %w", err)
	}

	// Run request benchmarking phase
	fmt.Println("Running request benchmarking phase...")
	if err := b.runRequestBenchmarking(); err != nil {
		return fmt.Errorf("request benchmarking phase failed: %w", err)
	}

	// Calculate derived metrics
	fmt.Println("Calculating derived metrics...")
	b.calculateDerivedMetrics()

	// Calculate statistics
	fmt.Println("Calculating statistics...")
	b.calculateStatistics()

	fmt.Println("Benchmark completed successfully!")
	return nil
}

// runWarmup executes warmup requests for each proxy
func (b *BenchmarkEngine) runWarmup() error {
	var wg sync.WaitGroup

	for _, proxy := range b.proxies {
		wg.Add(1)
		go func(p *Proxy) {
			defer wg.Done()
			b.runWarmupForProxy(p)
		}(proxy)
	}

	wg.Wait()
	return nil
}

// runWarmupForProxy executes warmup requests for a single proxy
func (b *BenchmarkEngine) runWarmupForProxy(proxy *Proxy) {
	fmt.Printf("Running warmup for proxy %s...\n", proxy.Address())

	timeout := time.Duration(b.config.Benchmark.TimeoutMs) * time.Millisecond

	for i := 0; i < b.config.Benchmark.WarmupRequests; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var err error
		switch proxy.Protocol {
		case "http", "https":
			client, err := NewHTTPClient(proxy, timeout)
			if err != nil {
				fmt.Printf("Failed to create HTTP client for proxy %s: %v\n", proxy.Address(), err)
				return
			}
			body, err := client.MakeRequest(ctx, b.config.Benchmark.TargetURL)
			if err == nil && b.config.Benchmark.ResponseValidation != nil && b.config.Benchmark.ResponseValidation.Enabled {
				err = b.validateResponse(body)
			}
		case "socks":
			client, err := NewSOCKS5Client(proxy, timeout)
			if err != nil {
				fmt.Printf("Failed to create SOCKS5 client for proxy %s: %v\n", proxy.Address(), err)
				return
			}
			body, err := client.MakeRequest(ctx, b.config.Benchmark.TargetURL)
			if err == nil && b.config.Benchmark.ResponseValidation != nil && b.config.Benchmark.ResponseValidation.Enabled {
				err = b.validateResponse(body)
			}
		default:
			fmt.Printf("Unsupported protocol for proxy %s: %s\n", proxy.Address(), proxy.Protocol)
			return
		}

		if err != nil {
			fmt.Printf("Warmup request failed for proxy %s: %v\n", proxy.Address(), err)
		}
	}
}

// runPingMeasurement executes ping measurements for each proxy
func (b *BenchmarkEngine) runPingMeasurement() error {
	var wg sync.WaitGroup

	for _, proxy := range b.proxies {
		wg.Add(1)
		go func(p *Proxy) {
			defer wg.Done()
			b.runPingMeasurementForProxy(p)
		}(proxy)
	}

	wg.Wait()
	return nil
}

// runPingMeasurementForProxy executes ping measurements for a single proxy
func (b *BenchmarkEngine) runPingMeasurementForProxy(proxy *Proxy) {
	fmt.Printf("Running ping measurement for proxy %s...\n", proxy.Address())

	timeout := time.Duration(b.config.Benchmark.TimeoutMs) * time.Millisecond
	interval := time.Duration(b.config.Benchmark.IntervalMs) * time.Millisecond
	pingClient := NewPingClient(timeout)

	for i := 0; i < b.config.Benchmark.Requests; i++ {
		if i > 0 {
			time.Sleep(interval)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Measure direct TCP connection time to proxy
		duration, err := pingClient.PingProxy(ctx, proxy)
		if err != nil {
			fmt.Printf("Ping failed for proxy %s: %v\n", proxy.Address(), err)
			b.metrics[proxy.String()].AddPingTime(0) // Add 0 for failed requests
		} else {
			b.metrics[proxy.String()].AddPingTime(duration)
		}
	}
}

// runRequestBenchmarking executes request benchmarking for each proxy
func (b *BenchmarkEngine) runRequestBenchmarking() error {
	var wg sync.WaitGroup

	for _, proxy := range b.proxies {
		wg.Add(1)
		go func(p *Proxy) {
			defer wg.Done()
			b.runRequestBenchmarkingForProxy(p)
		}(proxy)
	}

	wg.Wait()
	return nil
}

// runRequestBenchmarkingForProxy executes request benchmarking for a single proxy
func (b *BenchmarkEngine) runRequestBenchmarkingForProxy(proxy *Proxy) {
	fmt.Printf("Running request benchmarking for proxy %s...\n", proxy.Address())

	timeout := time.Duration(b.config.Benchmark.TimeoutMs) * time.Millisecond
	interval := time.Duration(b.config.Benchmark.IntervalMs) * time.Millisecond

	for i := 0; i < b.config.Benchmark.Requests; i++ {
		if i > 0 {
			time.Sleep(interval)
		}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		start := time.Now()
		var err error
		switch proxy.Protocol {
		case "http", "https":
			client, err := NewHTTPClient(proxy, timeout)
			if err != nil {
				fmt.Printf("Failed to create HTTP client for proxy %s: %v\n", proxy.Address(), err)
				b.metrics[proxy.String()].AddRequestTime(0, false)
				continue
			}
			body, err := client.MakeRequest(ctx, b.config.Benchmark.TargetURL)
			if err == nil && b.config.Benchmark.ResponseValidation != nil && b.config.Benchmark.ResponseValidation.Enabled {
				err = b.validateResponse(body)
			}
		case "socks":
			client, err := NewSOCKS5Client(proxy, timeout)
			if err != nil {
				fmt.Printf("Failed to create SOCKS5 client for proxy %s: %v\n", proxy.Address(), err)
				b.metrics[proxy.String()].AddRequestTime(0, false)
				continue
			}
			body, err := client.MakeRequest(ctx, b.config.Benchmark.TargetURL)
			if err == nil && b.config.Benchmark.ResponseValidation != nil && b.config.Benchmark.ResponseValidation.Enabled {
				err = b.validateResponse(body)
			}
		default:
			fmt.Printf("Unsupported protocol for proxy %s: %s\n", proxy.Address(), proxy.Protocol)
			b.metrics[proxy.String()].AddRequestTime(0, false)
			continue
		}

		duration := time.Since(start)
		if err != nil {
			if b.config.Benchmark.ResponseValidation != nil && b.config.Benchmark.ResponseValidation.Enabled {
				fmt.Printf("Request/Validation failed for proxy %s: %v\n", proxy.Address(), err)
			} else {
				fmt.Printf("Request failed for proxy %s: %v\n", proxy.Address(), err)
			}
			b.metrics[proxy.String()].AddRequestTime(duration, false)
		} else {
			b.metrics[proxy.String()].AddRequestTime(duration, true)
		}
	}
}

// calculateDerivedMetrics calculates derived processing times
func (b *BenchmarkEngine) calculateDerivedMetrics() {
	for _, metrics := range b.metrics {
		requestTimes := metrics.GetRequestTimes()
		pingTimes := metrics.GetPingTimes()

		// Calculate derived times (request time - ping*2)
		for i := 0; i < len(requestTimes) && i < len(pingTimes); i++ {
			derivedTime := requestTimes[i] - (pingTimes[i] * 2)
			// Ensure derived time is not negative
			if derivedTime < 0 {
				derivedTime = 0
			}
			metrics.AddDerivedTime(derivedTime)
		}
	}
}

// calculateStatistics calculates statistics for all metrics
func (b *BenchmarkEngine) calculateStatistics() {
	for _, metrics := range b.metrics {
		UpdateMetricsStatistics(metrics, &b.config.Statistics)
	}
}

// validateResponse validates the response body against configured checks
func (b *BenchmarkEngine) validateResponse(body []byte) error {
	if b.config.Benchmark.ResponseValidation == nil || !b.config.Benchmark.ResponseValidation.Enabled {
		return nil
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return fmt.Errorf("failed to parse JSON response: %w", err)
	}

	for _, check := range b.config.Benchmark.ResponseValidation.Checks {
		value, err := getNestedValue(data, check.Path)
		if err != nil {
			return fmt.Errorf("validation failed for path '%s': %w", check.Path, err)
		}

		if err := validateType(value, check.Type, check.Value); err != nil {
			return fmt.Errorf("validation failed for path '%s': %w", check.Path, err)
		}
	}

	return nil
}

// getNestedValue retrieves a value from a nested map using dot notation path
func getNestedValue(data map[string]interface{}, path string) (interface{}, error) {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("path not found: %s", part)
			}
		} else {
			return nil, fmt.Errorf("cannot navigate through non-object at %s", part)
		}
	}

	return current, nil
}

// validateType checks if a value matches the expected type and optional value
func validateType(value interface{}, expectedType string, expectedValue interface{}) error {
	switch expectedType {
	case "boolean":
		boolVal, ok := value.(bool)
		if !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
		if expectedValue != nil {
			if expectedBool, ok := expectedValue.(bool); ok {
				if boolVal != expectedBool {
					return fmt.Errorf("expected value %v, got %v", expectedBool, boolVal)
				}
			}
		}
	case "number":
		_, ok := value.(float64)
		if !ok {
			_, ok = value.(int)
			if !ok {
				return fmt.Errorf("expected number, got %T", value)
			}
		}
		if expectedValue != nil {
			// Compare numeric values if provided
			var numVal float64
			switch v := value.(type) {
			case float64:
				numVal = v
			case int:
				numVal = float64(v)
			}
			var expectedNum float64
			switch v := expectedValue.(type) {
			case float64:
				expectedNum = v
			case int:
				expectedNum = float64(v)
			}
			if numVal != expectedNum {
				return fmt.Errorf("expected value %v, got %v", expectedNum, numVal)
			}
		}
	case "string":
		strVal, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
		if expectedValue != nil {
			if expectedStr, ok := expectedValue.(string); ok {
				if strVal != expectedStr {
					return fmt.Errorf("expected value %v, got %v", expectedStr, strVal)
				}
			}
		}
	case "array":
		_, ok := value.([]interface{})
		if !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	case "object":
		_, ok := value.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	default:
		return fmt.Errorf("unknown type: %s", expectedType)
	}

	return nil
}

// GetResults returns the benchmark results
func (b *BenchmarkEngine) GetResults() map[string]*Metrics {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Return a copy of the metrics
	results := make(map[string]*Metrics)
	for k, v := range b.metrics {
		results[k] = v
	}
	return results
}
