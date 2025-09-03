package main

import (
	"encoding/json"
	"os"
	"time"
)

// BenchmarkResult represents the complete benchmark result
type BenchmarkResult struct {
	Timestamp time.Time       `json:"timestamp"`
	Proxies   []*ProxyMetrics `json:"proxies"`
}

// ShortSummary represents a concise summary with only mean delivered per proxy
type ShortSummary struct {
	Timestamp time.Time                    `json:"timestamp"`
	Proxies   map[string]float64          `json:"proxies"`
}

// ProxyMetrics represents metrics for a single proxy
type ProxyMetrics struct {
	ProxyString    string         `json:"proxy"`
	RequestMetrics RequestMetrics `json:"request_metrics"`
	PingMetrics    PingMetrics    `json:"ping_metrics"`
	DerivedMetrics DerivedMetrics `json:"derived_metrics"`
}

// Reporter generates benchmark reports
type Reporter struct{}

// NewReporter creates a new reporter
func NewReporter() *Reporter {
	return &Reporter{}
}

// GenerateReport generates a benchmark report from metrics
func (r *Reporter) GenerateReport(metrics map[string]*Metrics) *BenchmarkResult {
	result := &BenchmarkResult{
		Timestamp: time.Now(),
		Proxies:   make([]*ProxyMetrics, 0),
	}

	for _, m := range metrics {
		proxyMetrics := &ProxyMetrics{
			ProxyString:    m.ProxyString,
			RequestMetrics: m.RequestMetrics,
			PingMetrics:    m.PingMetrics,
			DerivedMetrics: m.DerivedMetrics,
		}
		result.Proxies = append(result.Proxies, proxyMetrics)
	}

	return result
}

// GenerateShortSummary generates a short summary with only mean delivered per proxy
func (r *Reporter) GenerateShortSummary(metrics map[string]*Metrics) *ShortSummary {
	summary := &ShortSummary{
		Timestamp: time.Now(),
		Proxies:   make(map[string]float64),
	}

	for _, m := range metrics {
		if m.DerivedMetrics.Statistics != nil {
			summary.Proxies[m.ProxyString] = m.DerivedMetrics.Statistics.Mean
		} else {
			summary.Proxies[m.ProxyString] = 0.0
		}
	}

	return summary
}

// SaveReport saves the benchmark report to a JSON file
func (r *Reporter) SaveReport(result *BenchmarkResult, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// SaveShortSummary saves the short summary to a JSON file
func (r *Reporter) SaveShortSummary(summary *ShortSummary, filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(summary)
}
