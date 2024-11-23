package handlers

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/adedejiosvaldo/scalable_api/db"
	"github.com/gin-gonic/gin"
)

type BenchmarkResult struct {
	TotalRequests      int     `json:"total_requests"`
	SuccessfulRequests int     `json:"successful_requests"`
	FailedRequests     int     `json:"failed_requests"`
	TotalDuration      string  `json:"total_duration"`  // Changed to string
	AverageLatency     string  `json:"average_latency"` // Changed to string
	MaxLatency         string  `json:"max_latency"`     // Changed to string
	MinLatency         string  `json:"min_latency"`     // Changed to string
	RequestsPerSecond  float64 `json:"requests_per_second"`
	SuccessRate        string  `json:"success_rate"`       // New field
	ThroughputPerMin   float64 `json:"throughput_per_min"` // New field
}

func formatDuration(d time.Duration) string {
	if d.Hours() > 1 {
		return fmt.Sprintf("%.2f hours", d.Hours())
	}
	if d.Minutes() > 1 {
		return fmt.Sprintf("%.2f minutes", d.Minutes())
	}
	if d.Seconds() > 1 {
		return fmt.Sprintf("%.2f seconds", d.Seconds())
	}
	if d.Milliseconds() > 1 {
		return fmt.Sprintf("%d milliseconds", d.Milliseconds())
	}
	return fmt.Sprintf("%d microseconds", d.Microseconds())
}

func calculateSuccessRate(successful, total int) string {
	if total == 0 {
		return "0%"
	}
	rate := (float64(successful) / float64(total)) * 100
	return fmt.Sprintf("%.2f%%", rate)
}

func runConcurrentBenchmark(
	numRequests int,
	concurrency int,
	operation func(context.Context) error,
) *BenchmarkResult {
	result := &BenchmarkResult{
		TotalRequests: numRequests,
	}

	var wg sync.WaitGroup
	requestChan := make(chan struct{}, numRequests)
	resultChan := make(chan time.Duration, numRequests)

	start := time.Now()
	var maxLatency time.Duration
	minLatency := time.Hour // Initialize with a large value

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range requestChan {
				requestStart := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := operation(ctx)
				duration := time.Since(requestStart)

				if err != nil {
					log.Printf("Request failed: %v", err)
					result.FailedRequests++
					continue
				}
				resultChan <- duration
			}
		}()
	}

	for i := 0; i < numRequests; i++ {
		requestChan <- struct{}{}
	}
	close(requestChan)

	wg.Wait()
	close(resultChan)

	var totalLatency time.Duration
	result.SuccessfulRequests = len(resultChan)

	for latency := range resultChan {
		totalLatency += latency
		if latency > maxLatency {
			maxLatency = latency
		}
		if latency < minLatency {
			minLatency = latency
		}
	}

	totalDuration := time.Since(start)
	avgLatency := totalLatency / time.Duration(result.SuccessfulRequests)

	// Format all durations to readable strings
	result.TotalDuration = formatDuration(totalDuration)
	result.AverageLatency = formatDuration(avgLatency)
	result.MaxLatency = formatDuration(maxLatency)
	result.MinLatency = formatDuration(minLatency)

	// Calculate performance metrics
	result.RequestsPerSecond = float64(result.SuccessfulRequests) / totalDuration.Seconds()
	result.ThroughputPerMin = result.RequestsPerSecond * 60
	result.SuccessRate = calculateSuccessRate(result.SuccessfulRequests, result.TotalRequests)

	return result
}

func SQLReadBenchmarkHandler(c *gin.Context) {
	numRequests := 100000 // Replace with query param for dynamic input
	concurrency := 50     // Replace with query param for dynamic input

	// Generate a valid ID for reading (assuming we have data in the database)
	readID := 1 // You might want to make this dynamic or random within a range

	result := runConcurrentBenchmark(numRequests, concurrency, func(ctx context.Context) error {
		_, err := db.ReadFromPostgres(readID)
		return err
	})

	c.JSON(200, gin.H{
		"operation": "read",
		"metrics":   result,
	})
}

func SQLWriteBenchmarkHandler(c *gin.Context) {
	numRequests := 100000
	concurrency := 50

	result := runConcurrentBenchmark(numRequests, concurrency, func(ctx context.Context) error {
		return db.InsertToPostgres("test data")
	})

	c.JSON(200, gin.H{
		"operation": "write",
		"metrics":   result,
	})
}
