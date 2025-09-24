package main

import (
	"os"
	"strings"
)

// GetTestProxies returns test proxies from TEST_PROXIES environment variable
// If the environment variable is not set, returns the fallback values
func GetTestProxies(fallback []string) []string {
	value := os.Getenv("TEST_PROXIES")
	if value == "" {
		return fallback
	}
	
	proxies := strings.Split(value, ",")
	// Trim whitespace from each proxy
	for i, proxy := range proxies {
		proxies[i] = strings.TrimSpace(proxy)
	}
	
	return proxies
}

// LoadTestConfig loads test configuration with environment variable support
func LoadTestConfig() *Config {
	return &Config{
		Proxies: GetTestProxies([]string{
			"http:proxy.example.com:8080:username:password:enabled",
			"socks:proxy.example.com:1080:username:password:enabled",
		}),
		Benchmark: BenchmarkConfig{
			Requests:       2,
			IntervalMs:     1000,
			WarmupRequests: 1,
			TargetURL:      "https://httpbin.org/get",
			Concurrency:    1,
			TimeoutMs:      5000,
		},
		Statistics: StatisticsConfig{
			Percentiles: []float64{90, 95, 99},
			Mean:        true,
			Median:      true,
		},
	}
}

// LoadGitHubTestConfig loads GitHub API test configuration
func LoadGitHubTestConfig() *Config {
	// Use first proxy from TEST_PROXIES or fallback
	proxies := GetTestProxies([]string{"http:proxy.example.com:8080:username:password:enabled"})
	return &Config{
		Proxies: []string{proxies[0]}, // Use first proxy
		Benchmark: BenchmarkConfig{
			Requests:       2,
			IntervalMs:     1000,
			WarmupRequests: 1,
			TargetURL:      "https://api.github.com/users/octocat",
			Concurrency:    1,
			TimeoutMs:      5000,
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{Path: "login", Type: "string", Value: "octocat"},
					{Path: "id", Type: "number"},
					{Path: "type", Type: "string", Value: "User"},
					{Path: "site_admin", Type: "boolean"},
					{Path: "public_repos", Type: "number"},
				},
			},
		},
		Statistics: StatisticsConfig{
			Percentiles: []float64{90, 95, 99},
			Mean:        true,
			Median:      true,
		},
	}
}

// LoadJSONPlaceholderTestConfig loads JSONPlaceholder test configuration
func LoadJSONPlaceholderTestConfig() *Config {
	// Use first proxy from TEST_PROXIES or fallback
	proxies := GetTestProxies([]string{"http:proxy.example.com:8080:username:password:enabled"})
	return &Config{
		Proxies: []string{proxies[0]}, // Use first proxy
		Benchmark: BenchmarkConfig{
			Requests:       2,
			IntervalMs:     1000,
			WarmupRequests: 1,
			TargetURL:      "https://jsonplaceholder.typicode.com/posts/1",
			Concurrency:    1,
			TimeoutMs:      5000,
			ResponseValidation: &ResponseValidation{
				Enabled: true,
				Checks: []ValidationCheck{
					{Path: "userId", Type: "number"},
					{Path: "id", Type: "number", Value: float64(1)},
					{Path: "title", Type: "string"},
					{Path: "body", Type: "string"},
				},
			},
		},
		Statistics: StatisticsConfig{
			Percentiles: []float64{90, 95, 99},
			Mean:        true,
			Median:      true,
		},
	}
}