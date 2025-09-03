package main

import (
	"testing"
	"time"
)

func TestMetrics(t *testing.T) {
	// Create a new metrics instance
	metrics := NewMetrics("test-proxy")

	// Add some request times
	metrics.AddRequestTime(100*time.Millisecond, true)
	metrics.AddRequestTime(200*time.Millisecond, true)
	metrics.AddRequestTime(150*time.Millisecond, false) // Failed request

	// Add some ping times
	metrics.AddPingTime(50 * time.Millisecond)
	metrics.AddPingTime(60 * time.Millisecond)
	metrics.AddPingTime(40 * time.Millisecond)

	// Add some derived times
	metrics.AddDerivedTime(100)
	metrics.AddDerivedTime(120)
	metrics.AddDerivedTime(90)

	// Check request metrics
	if metrics.RequestMetrics.Total != 3 {
		t.Errorf("Expected Total=3, got %d", metrics.RequestMetrics.Total)
	}
	if metrics.RequestMetrics.Successful != 2 {
		t.Errorf("Expected Successful=2, got %d", metrics.RequestMetrics.Successful)
	}
	if metrics.RequestMetrics.Failed != 1 {
		t.Errorf("Expected Failed=1, got %d", metrics.RequestMetrics.Failed)
	}
	if len(metrics.RequestMetrics.Times) != 2 {
		t.Errorf("Expected 2 request times (successful only), got %d", len(metrics.RequestMetrics.Times))
	}

	// Check ping metrics
	if len(metrics.PingMetrics.Times) != 3 {
		t.Errorf("Expected 3 ping times, got %d", len(metrics.PingMetrics.Times))
	}

	// Check derived metrics
	if len(metrics.DerivedMetrics.ProcessingTimes) != 3 {
		t.Errorf("Expected 3 derived times, got %d", len(metrics.DerivedMetrics.ProcessingTimes))
	}
}

func TestStatisticsCalculation(t *testing.T) {
	// Create test data
	values := []int64{100, 200, 150, 175, 125}
	
	// Create statistics config
	config := &StatisticsConfig{
		Percentiles: []float64{90, 95, 99},
		Mean:        true,
		Median:      true,
	}

	// Calculate statistics
	stats := CalculateStatistics(values, config)

	// Basic checks
	if stats == nil {
		t.Fatal("Expected statistics, got nil")
	}
	
	if stats.Min != 100 {
		t.Errorf("Expected Min=100, got %d", stats.Min)
	}
	
	if stats.Max != 200 {
		t.Errorf("Expected Max=200, got %d", stats.Max)
	}
	
	// Check that mean and median were calculated
	if stats.Mean == 0 {
		t.Error("Expected mean to be calculated")
	}
	
	if stats.Median == 0 {
		t.Error("Expected median to be calculated")
	}
	
	// Check that percentiles were calculated
	if len(stats.Percentiles) == 0 {
		t.Error("Expected percentiles to be calculated")
	}
}