package main

import (
	"sync"
	"time"
)

// Metrics holds all metrics for a single proxy
type Metrics struct {
	ProxyString    string         `json:"proxy"`
	RequestMetrics RequestMetrics `json:"request_metrics"`
	PingMetrics    PingMetrics    `json:"ping_metrics"`
	DerivedMetrics DerivedMetrics `json:"derived_metrics"`
	mu             sync.Mutex
}

// RequestMetrics holds request timing metrics
type RequestMetrics struct {
	Total      int         `json:"total"`
	Successful int         `json:"successful"`
	Failed     int         `json:"failed"`
	Times      []int64     `json:"times"`
	Statistics *Statistics `json:"statistics,omitempty"`
}

// PingMetrics holds ping timing metrics
type PingMetrics struct {
	Times      []int64     `json:"times"`
	Statistics *Statistics `json:"statistics,omitempty"`
}

// DerivedMetrics holds derived timing metrics (request time - ping*2)
type DerivedMetrics struct {
	ProcessingTimes []int64     `json:"processing_times"`
	Statistics      *Statistics `json:"statistics,omitempty"`
}

// Statistics holds calculated statistical values
type Statistics struct {
	Min         int64               `json:"min"`
	Max         int64               `json:"max"`
	Mean        float64             `json:"mean,omitempty"`
	Median      float64             `json:"median,omitempty"`
	StdDev      float64             `json:"std_dev"`
	Percentiles map[string]float64 `json:"percentiles,omitempty"`
}

// NewMetrics creates a new Metrics instance for a proxy
func NewMetrics(proxyString string) *Metrics {
	return &Metrics{
		ProxyString: proxyString,
		RequestMetrics: RequestMetrics{
			Times: make([]int64, 0),
		},
		PingMetrics: PingMetrics{
			Times: make([]int64, 0),
		},
		DerivedMetrics: DerivedMetrics{
			ProcessingTimes: make([]int64, 0),
		},
	}
}

// AddRequestTime adds a request time measurement
func (m *Metrics) AddRequestTime(duration time.Duration, success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RequestMetrics.Total++
	if success {
		m.RequestMetrics.Successful++
		m.RequestMetrics.Times = append(m.RequestMetrics.Times, duration.Milliseconds())
	} else {
		m.RequestMetrics.Failed++
	}
}

// AddPingTime adds a ping time measurement
func (m *Metrics) AddPingTime(duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.PingMetrics.Times = append(m.PingMetrics.Times, duration.Milliseconds())
}

// AddDerivedTime adds a derived processing time measurement
func (m *Metrics) AddDerivedTime(duration int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.DerivedMetrics.ProcessingTimes = append(m.DerivedMetrics.ProcessingTimes, duration)
}

// GetRequestTimes returns a copy of request times
func (m *Metrics) GetRequestTimes() []int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	times := make([]int64, len(m.RequestMetrics.Times))
	copy(times, m.RequestMetrics.Times)
	return times
}

// GetPingTimes returns a copy of ping times
func (m *Metrics) GetPingTimes() []int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	times := make([]int64, len(m.PingMetrics.Times))
	copy(times, m.PingMetrics.Times)
	return times
}

// GetDerivedTimes returns a copy of derived processing times
func (m *Metrics) GetDerivedTimes() []int64 {
	m.mu.Lock()
	defer m.mu.Unlock()

	times := make([]int64, len(m.DerivedMetrics.ProcessingTimes))
	copy(times, m.DerivedMetrics.ProcessingTimes)
	return times
}
