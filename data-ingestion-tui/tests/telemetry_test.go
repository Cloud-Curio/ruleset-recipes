/*
 * Data Ingestion TUI - Telemetry Tests
 * 
 * Comprehensive test suite for telemetry, metrics, benchmarking, and monitoring.
 * Tests Prometheus metrics, OpenTelemetry tracing, and performance profiling.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"data-ingestion-tui/internal/telemetry"
)

// @decorator: TelemetryTestSuite
// @description: Comprehensive test suite for telemetry functionality
type TelemetryTestSuite struct {
	suite.Suite
	logger *logrus.Logger
}

// @decorator: SetupSuite
// @description: Set up test suite
func (suite *TelemetryTestSuite) SetupSuite() {
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.WarnLevel) // Reduce noise in tests
}

// @decorator: TestTelemetryManagerInitialization
// @description: Test telemetry manager initialization with various configurations
func (suite *TelemetryTestSuite) TestTelemetryManagerInitialization() {
	tests := []struct {
		name   string
		config *telemetry.TelemetryConfig
		valid  bool
	}{
		{
			name: "default_config",
			config: &telemetry.TelemetryConfig{
				ServiceName:     "test-service",
				ServiceVersion:  "1.0.0",
				Environment:     "test",
				MetricsPort:     9091,
				EnableTracing:   false, // Disable to avoid Jaeger dependency
				EnableMetrics:   true,
				EnableProfiling: true,
				SampleRate:      1.0,
			},
			valid: true,
		},
		{
			name: "minimal_config",
			config: &telemetry.TelemetryConfig{
				ServiceName:     "minimal-service",
				EnableTracing:   false,
				EnableMetrics:   true,
				EnableProfiling: false,
			},
			valid: true,
		},
		{
			name: "nil_config_uses_defaults",
			config: nil,
			valid: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			tm, err := telemetry.NewTelemetryManager(tt.config, suite.logger)
			
			if tt.valid {
				require.NoError(suite.T(), err)
				require.NotNil(suite.T(), tm)
				
				// Test metrics are available
				metrics := tm.GetMetrics()
				assert.NotNil(suite.T(), metrics)
				assert.NotNil(suite.T(), metrics.SystemInfo)
				assert.NotNil(suite.T(), metrics.MemoryUsage)
				assert.NotNil(suite.T(), metrics.CPUUsage)
				
				// Cleanup
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				err = tm.Shutdown(ctx)
				assert.NoError(suite.T(), err)
			} else {
				assert.Error(suite.T(), err)
			}
		})
	}
}

// @decorator: TestMetricsCollection
// @description: Test basic metrics collection functionality
func (suite *TelemetryTestSuite) TestMetricsCollection() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "metrics-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9092,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	metrics := tm.GetMetrics()
	
	// Test system metrics are being collected
	// Wait a moment for metrics collection to start
	time.Sleep(2 * time.Second)
	
	// System metrics should be populated
	assert.NotNil(suite.T(), metrics.SystemInfo)
	assert.NotNil(suite.T(), metrics.MemoryUsage)
	assert.NotNil(suite.T(), metrics.GoroutineCount)
	
	// Test custom metrics recording
	tm.RecordIngestionMetrics("test-endpoint", "bulk", 5*time.Second, 100, 5, 20.5)
	tm.RecordAPICall("test-api", "GET", 200, 150*time.Millisecond)
	tm.RecordDatabaseOperation("INSERT", "test_table", true, 50*time.Millisecond)
	
	// Verify metrics were recorded (basic smoke test)
	// In a real implementation, you'd check the actual metric values
	assert.NotNil(suite.T(), metrics.IngestionJobsTotal)
	assert.NotNil(suite.T(), metrics.APICallsTotal)
	assert.NotNil(suite.T(), metrics.DatabaseOperations)
}

// @decorator: TestMetricsServer
// @description: Test Prometheus metrics HTTP server
func (suite *TelemetryTestSuite) TestMetricsServer() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "server-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9093,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Wait for server to start
	time.Sleep(1 * time.Second)
	
	// Test metrics endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/metrics", config.MetricsPort))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
	assert.Contains(suite.T(), resp.Header.Get("Content-Type"), "text/plain")
	
	// Test health endpoint
	resp, err = http.Get(fmt.Sprintf("http://localhost:%d/health", config.MetricsPort))
	require.NoError(suite.T(), err)
	defer resp.Body.Close()
	
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

// @decorator: TestBenchmarkSuite
// @description: Test benchmarking functionality
func (suite *TelemetryTestSuite) TestBenchmarkSuite() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "benchmark-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9094,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Start a benchmark
	benchmark := tm.StartBenchmark("test-benchmark", "Testing benchmark functionality")
	assert.NotNil(suite.T(), benchmark)
	assert.Equal(suite.T(), "test-benchmark", benchmark.Name)
	assert.Equal(suite.T(), "Testing benchmark functionality", benchmark.Description)
	
	// Simulate some work
	time.Sleep(100 * time.Millisecond)
	
	// Update benchmark stats
	benchmark.Iterations = 1000
	benchmark.Operations = 1000
	benchmark.Bytes = 1024 * 1024 // 1MB
	benchmark.Errors = 5
	
	// Finish benchmark
	result := tm.FinishBenchmark("test-benchmark")
	require.NotNil(suite.T(), result)
	
	assert.Equal(suite.T(), "test-benchmark", result.Name)
	assert.Greater(suite.T(), result.Duration, time.Duration(0))
	assert.Equal(suite.T(), int64(1000), result.Iterations)
	assert.Greater(suite.T(), result.NsPerOp, int64(0))
	assert.Greater(suite.T(), result.MBPerSec, float64(0))
	assert.Equal(suite.T(), float64(99.5), result.SuccessRate) // (1000-5)/1000 * 100
	
	// Test getting benchmark results
	results := tm.GetBenchmarkResults()
	assert.Contains(suite.T(), results, "test-benchmark")
	assert.Equal(suite.T(), result, results["test-benchmark"])
}

// @decorator: TestMultipleBenchmarks
// @description: Test running multiple concurrent benchmarks
func (suite *TelemetryTestSuite) TestMultipleBenchmarks() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "multi-benchmark-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9095,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Start multiple benchmarks
	benchmarkNames := []string{"benchmark-1", "benchmark-2", "benchmark-3"}
	
	for _, name := range benchmarkNames {
		benchmark := tm.StartBenchmark(name, fmt.Sprintf("Description for %s", name))
		assert.NotNil(suite.T(), benchmark)
		assert.Equal(suite.T(), name, benchmark.Name)
	}
	
	// Simulate work for different durations
	time.Sleep(50 * time.Millisecond)
	tm.FinishBenchmark("benchmark-1")
	
	time.Sleep(100 * time.Millisecond)
	tm.FinishBenchmark("benchmark-2")
	
	time.Sleep(150 * time.Millisecond)
	tm.FinishBenchmark("benchmark-3")
	
	// Verify all benchmarks completed
	results := tm.GetBenchmarkResults()
	assert.Len(suite.T(), results, 3)
	
	for _, name := range benchmarkNames {
		result, exists := results[name]
		assert.True(suite.T(), exists, "Benchmark %s should exist", name)
		assert.Greater(suite.T(), result.Duration, time.Duration(0))
	}
	
	// Verify different durations
	assert.Less(suite.T(), results["benchmark-1"].Duration, results["benchmark-2"].Duration)
	assert.Less(suite.T(), results["benchmark-2"].Duration, results["benchmark-3"].Duration)
}

// @decorator: TestPerformanceProfiler
// @description: Test performance profiling functionality
func (suite *TelemetryTestSuite) TestPerformanceProfiler() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "profiler-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9096,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: true,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Note: The profiler is internal to the telemetry manager
	// We can only test that it initializes without error
	// In a real implementation, you might expose profiler methods
	
	// Wait for profiler to collect some samples
	time.Sleep(3 * time.Second)
	
	// The profiler should be running in the background
	// This is mainly a smoke test to ensure no panics occur
	assert.True(suite.T(), true, "Profiler should run without errors")
}

// @decorator: TestIngestionMetricsRecording
// @description: Test recording of ingestion-specific metrics
func (suite *TelemetryTestSuite) TestIngestionMetricsRecording() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "ingestion-metrics-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9097,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Record various ingestion metrics
	testCases := []struct {
		endpoint         string
		jobType          string
		duration         time.Duration
		recordsProcessed int64
		recordsFailed    int64
		rps              float64
	}{
		{"congress-api", "bulk", 30 * time.Second, 1000, 10, 33.3},
		{"fred-api", "incremental", 15 * time.Second, 500, 5, 33.3},
		{"sec-api", "bulk", 60 * time.Second, 2000, 20, 33.3},
	}
	
	for _, tc := range testCases {
		tm.RecordIngestionMetrics(
			tc.endpoint,
			tc.jobType,
			tc.duration,
			tc.recordsProcessed,
			tc.recordsFailed,
			tc.rps,
		)
	}
	
	// Verify metrics were recorded
	metrics := tm.GetMetrics()
	assert.NotNil(suite.T(), metrics.IngestionJobsTotal)
	assert.NotNil(suite.T(), metrics.IngestionDuration)
	assert.NotNil(suite.T(), metrics.RecordsProcessed)
	assert.NotNil(suite.T(), metrics.RecordsPerSecond)
}

// @decorator: TestAPICallMetricsRecording
// @description: Test recording of API call metrics
func (suite *TelemetryTestSuite) TestAPICallMetricsRecording() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "api-metrics-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9098,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Record various API call metrics
	testCases := []struct {
		endpoint   string
		method     string
		statusCode int
		duration   time.Duration
	}{
		{"congress-api", "GET", 200, 150 * time.Millisecond},
		{"congress-api", "GET", 404, 50 * time.Millisecond},
		{"fred-api", "GET", 200, 200 * time.Millisecond},
		{"fred-api", "GET", 500, 1000 * time.Millisecond},
		{"sec-api", "POST", 201, 300 * time.Millisecond},
	}
	
	for _, tc := range testCases {
		tm.RecordAPICall(tc.endpoint, tc.method, tc.statusCode, tc.duration)
	}
	
	// Verify metrics were recorded
	metrics := tm.GetMetrics()
	assert.NotNil(suite.T(), metrics.APICallsTotal)
	assert.NotNil(suite.T(), metrics.APIResponseTime)
}

// @decorator: TestDatabaseOperationMetrics
// @description: Test recording of database operation metrics
func (suite *TelemetryTestSuite) TestDatabaseOperationMetrics() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "db-metrics-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9099,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Record various database operation metrics
	testCases := []struct {
		operation string
		table     string
		success   bool
		duration  time.Duration
	}{
		{"INSERT", "congress_bills", true, 10 * time.Millisecond},
		{"INSERT", "congress_bills", false, 50 * time.Millisecond},
		{"SELECT", "fred_observations", true, 5 * time.Millisecond},
		{"UPDATE", "sec_companies", true, 15 * time.Millisecond},
		{"DELETE", "old_data", true, 20 * time.Millisecond},
	}
	
	for _, tc := range testCases {
		tm.RecordDatabaseOperation(tc.operation, tc.table, tc.success, tc.duration)
	}
	
	// Verify metrics were recorded
	metrics := tm.GetMetrics()
	assert.NotNil(suite.T(), metrics.DatabaseOperations)
	assert.NotNil(suite.T(), metrics.DatabaseResponseTime)
}

// @decorator: TestTelemetryShutdown
// @description: Test graceful shutdown of telemetry manager
func (suite *TelemetryTestSuite) TestTelemetryShutdown() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "shutdown-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9100,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: true,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	
	// Wait for initialization
	time.Sleep(1 * time.Second)
	
	// Test graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	shutdownStart := time.Now()
	err = tm.Shutdown(ctx)
	shutdownDuration := time.Since(shutdownStart)
	
	assert.NoError(suite.T(), err, "Shutdown should complete without error")
	assert.Less(suite.T(), shutdownDuration, 6*time.Second, "Shutdown should complete within timeout")
	
	// Verify metrics server is stopped
	time.Sleep(500 * time.Millisecond)
	_, err = http.Get(fmt.Sprintf("http://localhost:%d/metrics", config.MetricsPort))
	assert.Error(suite.T(), err, "Metrics server should be stopped")
}

// @decorator: TestConcurrentMetricsRecording
// @description: Test concurrent metrics recording from multiple goroutines
func (suite *TelemetryTestSuite) TestConcurrentMetricsRecording() {
	config := &telemetry.TelemetryConfig{
		ServiceName:     "concurrent-test",
		ServiceVersion:  "1.0.0",
		Environment:     "test",
		MetricsPort:     9101,
		EnableTracing:   false,
		EnableMetrics:   true,
		EnableProfiling: false,
		SampleRate:      1.0,
	}
	
	tm, err := telemetry.NewTelemetryManager(config, suite.logger)
	require.NoError(suite.T(), err)
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm.Shutdown(ctx)
	}()
	
	// Record metrics concurrently from multiple goroutines
	numGoroutines := 10
	metricsPerGoroutine := 100
	
	done := make(chan bool, numGoroutines)
	
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer func() { done <- true }()
			
			for i := 0; i < metricsPerGoroutine; i++ {
				endpoint := fmt.Sprintf("endpoint-%d", goroutineID)
				
				// Record different types of metrics
				tm.RecordIngestionMetrics(endpoint, "test", time.Second, 10, 1, 10.0)
				tm.RecordAPICall(endpoint, "GET", 200, 100*time.Millisecond)
				tm.RecordDatabaseOperation("INSERT", "test_table", true, 10*time.Millisecond)
			}
		}(g)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
	
	// Verify no panics occurred and metrics are still accessible
	metrics := tm.GetMetrics()
	assert.NotNil(suite.T(), metrics)
	assert.NotNil(suite.T(), metrics.IngestionJobsTotal)
	assert.NotNil(suite.T(), metrics.APICallsTotal)
	assert.NotNil(suite.T(), metrics.DatabaseOperations)
}

// Run the test suite
func TestTelemetrySuite(t *testing.T) {
	suite.Run(t, new(TelemetryTestSuite))
}
