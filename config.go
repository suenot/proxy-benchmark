package main

import (
	"encoding/json"
	"os"
)

// Config represents the application configuration
type Config struct {
	Proxies    []string         `json:"proxies"`
	Benchmark  BenchmarkConfig  `json:"benchmark"`
	Statistics StatisticsConfig `json:"statistics"`
}

// BenchmarkConfig holds benchmark-specific configuration
type BenchmarkConfig struct {
	Requests           int                   `json:"requests"`
	IntervalMs         int                   `json:"interval_ms"`
	WarmupRequests     int                   `json:"warmup_requests"`
	TargetURL          string                `json:"target_url"`
	Concurrency        int                   `json:"concurrency"`
	TimeoutMs          int                   `json:"timeout_ms"`
	ResponseValidation *ResponseValidation   `json:"response_validation,omitempty"`
}

// ResponseValidation holds response validation configuration
type ResponseValidation struct {
	Enabled bool              `json:"enabled"`
	Checks  []ValidationCheck `json:"checks"`
}

// ValidationCheck defines a single validation rule
type ValidationCheck struct {
	Path  string      `json:"path"`
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

// StatisticsConfig holds statistics configuration
type StatisticsConfig struct {
	Percentiles []float64 `json:"percentiles"`
	Mean        bool      `json:"mean"`
	Median      bool      `json:"median"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filepath string) (*Config, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
