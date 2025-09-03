package main

import (
	"fmt"
	"github.com/montanaflynn/stats"
)

// CalculateStatistics calculates statistical metrics for a set of values
func CalculateStatistics(values []int64, config *StatisticsConfig) *Statistics {
	if len(values) == 0 {
		return nil
	}

	// Convert int64 to float64 for stats library
	floatValues := make([]float64, len(values))
	for i, v := range values {
		floatValues[i] = float64(v)
	}

	stat := &Statistics{}

	// Calculate min and max
	min, _ := stats.Min(floatValues)
	max, _ := stats.Max(floatValues)
	stat.Min = int64(min)
	stat.Max = int64(max)

	// Calculate mean if requested
	if config.Mean {
		mean, _ := stats.Mean(floatValues)
		stat.Mean = mean
	}

	// Calculate median if requested
	if config.Median {
		median, _ := stats.Median(floatValues)
		stat.Median = median
	}

	// Calculate standard deviation
	stdDev, _ := stats.StandardDeviation(floatValues)
	stat.StdDev = stdDev

	// Calculate percentiles if requested
	if len(config.Percentiles) > 0 {
		stat.Percentiles = make(map[string]float64)
		for _, p := range config.Percentiles {
			value, _ := stats.Percentile(floatValues, p)
			stat.Percentiles[fmt.Sprintf("%.1f", p)] = value
		}
	}

	return stat
}

// UpdateMetricsStatistics calculates and updates statistics for all metrics
func UpdateMetricsStatistics(metrics *Metrics, config *StatisticsConfig) {
	metrics.RequestMetrics.Statistics = CalculateStatistics(metrics.GetRequestTimes(), config)
	metrics.PingMetrics.Statistics = CalculateStatistics(metrics.GetPingTimes(), config)
	metrics.DerivedMetrics.Statistics = CalculateStatistics(metrics.GetDerivedTimes(), config)
}
