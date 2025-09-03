package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Parse command line flags
	configPath := flag.String("config", "config.json", "Path to the configuration file")
	flag.Parse()

	// Check if config file exists
	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		log.Fatalf("Configuration file not found: %s", *configPath)
	}

	// Load configuration
	fmt.Printf("Loading configuration from %s...\n", *configPath)
	config, err := LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Set default values if not specified
	if config.Benchmark.Requests == 0 {
		config.Benchmark.Requests = 100
	}
	if config.Benchmark.IntervalMs == 0 {
		config.Benchmark.IntervalMs = 5000
	}
	if config.Benchmark.WarmupRequests == 0 {
		config.Benchmark.WarmupRequests = 10
	}
	if config.Benchmark.TargetURL == "" {
		config.Benchmark.TargetURL = "https://httpbin.org/get"
	}
	if config.Benchmark.Concurrency == 0 {
		config.Benchmark.Concurrency = 10
	}
	if config.Benchmark.TimeoutMs == 0 {
		config.Benchmark.TimeoutMs = 30000
	}

	// Create benchmark engine
	fmt.Println("Initializing benchmark engine...")
	engine, err := NewBenchmarkEngine(config)
	if err != nil {
		log.Fatalf("Failed to create benchmark engine: %v", err)
	}

	// Run benchmark
	if err := engine.Run(); err != nil {
		log.Fatalf("Benchmark failed: %v", err)
	}

	// Generate and save report
	fmt.Println("Generating report...")
	reporter := NewReporter()
	results := engine.GetResults()
	report := reporter.GenerateReport(results)

	fmt.Println("Saving results to result.json...")
	if err := reporter.SaveReport(report, "result.json"); err != nil {
		log.Fatalf("Failed to save report: %v", err)
	}

	// Generate and save short summary
	fmt.Println("Generating short summary...")
	shortSummary := reporter.GenerateShortSummary(results)

	fmt.Println("Saving short summary to results_short.json...")
	if err := reporter.SaveShortSummary(shortSummary, "results_short.json"); err != nil {
		log.Fatalf("Failed to save short summary: %v", err)
	}

	fmt.Println("Benchmark results saved to result.json")
	fmt.Println("Short summary saved to results_short.json")
}
