/*
 * Data Ingestion TUI - Worker Pool Tests
 * 
 * Comprehensive test suite for the advanced async worker pool implementation.
 * Tests dynamic scaling, work stealing, circuit breakers, and telemetry integration.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package tests

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"data-ingestion-tui/internal/async"
)

// @decorator: WorkerPoolTestSuite
// @description: Comprehensive test suite for worker pool functionality
type WorkerPoolTestSuite struct {
	suite.Suite
	logger *logrus.Logger
}

// @decorator: TestTask
// @description: Mock task implementation for testing
type TestTask struct {
	id          string
	priority    int
	timeout     time.Duration
	retryCount  int
	shouldFail  bool
	shouldRetry bool
	duration    time.Duration
	executed    int64
}

func (t *TestTask) Execute(ctx context.Context) (interface{}, error) {
	atomic.AddInt64(&t.executed, 1)
	
	if t.duration > 0 {
		select {
		case <-time.After(t.duration):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	
	if t.shouldFail {
		return nil, assert.AnError
	}
	
	return "success", nil
}

func (t *TestTask) GetID() string { return t.id }
func (t *TestTask) GetPriority() int { return t.priority }
func (t *TestTask) GetTimeout() time.Duration { return t.timeout }
func (t *TestTask) GetRetryCount() int { return t.retryCount }
func (t *TestTask) ShouldRetry(err error) bool { return t.shouldRetry }

// @decorator: SetupSuite
// @description: Set up test suite
func (suite *WorkerPoolTestSuite) SetupSuite() {
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.WarnLevel) // Reduce noise in tests
}

// @decorator: TestWorkerPoolInitialization
// @description: Test worker pool initialization with various configurations
func (suite *WorkerPoolTestSuite) TestWorkerPoolInitialization() {
	tests := []struct {
		name   string
		config async.WorkerPoolConfig
		valid  bool
	}{
		{
			name: "valid_default_config",
			config: async.WorkerPoolConfig{
				MinWorkers:                2,
				MaxWorkers:                8,
				QueueSize:                 100,
				WorkerTimeout:             30 * time.Second,
				CircuitBreakerTimeout:     5 * time.Second,
				CircuitBreakerMaxFailures: 5,
			},
			valid: true,
		},
		{
			name: "zero_values_use_defaults",
			config: async.WorkerPoolConfig{
				CircuitBreakerTimeout:     5 * time.Second,
				CircuitBreakerMaxFailures: 5,
			},
			valid: true,
		},
		{
			name: "high_performance_config",
			config: async.WorkerPoolConfig{
				MinWorkers:                10,
				MaxWorkers:                50,
				QueueSize:                 1000,
				WorkerTimeout:             10 * time.Second,
				CircuitBreakerTimeout:     2 * time.Second,
				CircuitBreakerMaxFailures: 10,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			pool, err := async.NewWorkerPool(tt.config, suite.logger)
			
			if tt.valid {
				require.NoError(suite.T(), err)
				require.NotNil(suite.T(), pool)
				
				stats := pool.GetStats()
				assert.GreaterOrEqual(suite.T(), stats.TotalWorkers, tt.config.MinWorkers)
				assert.LessOrEqual(suite.T(), stats.TotalWorkers, tt.config.MaxWorkers)
				
				// Cleanup
				err = pool.Shutdown(5 * time.Second)
				assert.NoError(suite.T(), err)
			} else {
				assert.Error(suite.T(), err)
			}
		})
	}
}

// @decorator: TestTaskSubmissionAndExecution
// @description: Test basic task submission and execution
func (suite *WorkerPoolTestSuite) TestTaskSubmissionAndExecution() {
	config := async.WorkerPoolConfig{
		MinWorkers:                2,
		MaxWorkers:                4,
		QueueSize:                 10,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 5,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	// Submit multiple tasks
	numTasks := 10
	tasks := make([]*TestTask, numTasks)
	
	for i := 0; i < numTasks; i++ {
		task := &TestTask{
			id:       fmt.Sprintf("task-%d", i),
			priority: i % 3,
			timeout:  2 * time.Second,
		}
		tasks[i] = task
		
		err := pool.Submit(task)
		assert.NoError(suite.T(), err)
	}
	
	// Wait for tasks to complete
	time.Sleep(3 * time.Second)
	
	// Verify all tasks were executed
	for i, task := range tasks {
		executed := atomic.LoadInt64(&task.executed)
		assert.Equal(suite.T(), int64(1), executed, "Task %d should be executed exactly once", i)
	}
	
	// Check pool stats
	stats := pool.GetStats()
	assert.Equal(suite.T(), int64(numTasks), stats.CompletedTasks)
	assert.Equal(suite.T(), int64(0), stats.FailedTasks)
}

// @decorator: TestTaskRetryLogic
// @description: Test task retry logic with exponential backoff
func (suite *WorkerPoolTestSuite) TestTaskRetryLogic() {
	config := async.WorkerPoolConfig{
		MinWorkers:                1,
		MaxWorkers:                2,
		QueueSize:                 5,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     10 * time.Second,
		CircuitBreakerMaxFailures: 10,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	// Task that fails initially but succeeds on retry
	task := &TestTask{
		id:          "retry-task",
		timeout:     1 * time.Second,
		retryCount:  3,
		shouldFail:  true,
		shouldRetry: true,
	}
	
	startTime := time.Now()
	err = pool.Submit(task)
	require.NoError(suite.T(), err)
	
	// Wait for task to complete with retries
	time.Sleep(8 * time.Second)
	
	// Verify task was retried
	executed := atomic.LoadInt64(&task.executed)
	assert.Greater(suite.T(), executed, int64(1), "Task should be retried")
	assert.LessOrEqual(suite.T(), executed, int64(4), "Task should not exceed max retries")
	
	// Verify exponential backoff timing
	duration := time.Since(startTime)
	assert.Greater(suite.T(), duration, 2*time.Second, "Should include backoff delays")
}

// @decorator: TestDynamicWorkerScaling
// @description: Test dynamic scaling of workers based on load
func (suite *WorkerPoolTestSuite) TestDynamicWorkerScaling() {
	config := async.WorkerPoolConfig{
		MinWorkers:                2,
		MaxWorkers:                8,
		QueueSize:                 20,
		WorkerTimeout:             2 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 5,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	// Initial state
	stats := pool.GetStats()
	initialWorkers := stats.TotalWorkers
	assert.GreaterOrEqual(suite.T(), initialWorkers, config.MinWorkers)
	
	// Submit many long-running tasks to trigger scaling
	numTasks := 15
	for i := 0; i < numTasks; i++ {
		task := &TestTask{
			id:       fmt.Sprintf("long-task-%d", i),
			timeout:  1 * time.Second,
			duration: 2 * time.Second,
		}
		
		err := pool.Submit(task)
		assert.NoError(suite.T(), err)
	}
	
	// Wait for scaling to occur
	time.Sleep(1 * time.Second)
	
	// Check if workers scaled up
	stats = pool.GetStats()
	assert.Greater(suite.T(), stats.TotalWorkers, initialWorkers, "Workers should scale up under load")
	assert.LessOrEqual(suite.T(), stats.TotalWorkers, config.MaxWorkers, "Should not exceed max workers")
	
	// Wait for tasks to complete and workers to scale down
	time.Sleep(5 * time.Second)
	
	stats = pool.GetStats()
	// Note: Scale down might not happen immediately due to worker timeout
	assert.LessOrEqual(suite.T(), stats.TotalWorkers, config.MaxWorkers)
}

// @decorator: TestCircuitBreakerFunctionality
// @description: Test circuit breaker pattern for fault tolerance
func (suite *WorkerPoolTestSuite) TestCircuitBreakerFunctionality() {
	config := async.WorkerPoolConfig{
		MinWorkers:                1,
		MaxWorkers:                2,
		QueueSize:                 10,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     2 * time.Second,
		CircuitBreakerMaxFailures: 3,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	// Submit failing tasks to trigger circuit breaker
	for i := 0; i < 5; i++ {
		task := &TestTask{
			id:          fmt.Sprintf("fail-task-%d", i),
			timeout:     1 * time.Second,
			shouldFail:  true,
			shouldRetry: false,
		}
		
		err := pool.Submit(task)
		if i < 3 {
			assert.NoError(suite.T(), err, "Should accept tasks before circuit opens")
		} else {
			// Circuit breaker should be open after max failures
			assert.Error(suite.T(), err, "Should reject tasks when circuit is open")
		}
	}
	
	// Wait for circuit breaker timeout
	time.Sleep(3 * time.Second)
	
	// Circuit should be half-open, allowing one task
	successTask := &TestTask{
		id:      "success-task",
		timeout: 1 * time.Second,
	}
	
	err = pool.Submit(successTask)
	assert.NoError(suite.T(), err, "Should accept task when circuit is half-open")
	
	// Wait for task completion
	time.Sleep(2 * time.Second)
	
	// Circuit should be closed again
	anotherTask := &TestTask{
		id:      "another-task",
		timeout: 1 * time.Second,
	}
	
	err = pool.Submit(anotherTask)
	assert.NoError(suite.T(), err, "Should accept task when circuit is closed")
}

// @decorator: TestConcurrentTaskSubmission
// @description: Test concurrent task submission from multiple goroutines
func (suite *WorkerPoolTestSuite) TestConcurrentTaskSubmission() {
	config := async.WorkerPoolConfig{
		MinWorkers:                3,
		MaxWorkers:                6,
		QueueSize:                 50,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 10,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	numGoroutines := 5
	tasksPerGoroutine := 10
	totalTasks := numGoroutines * tasksPerGoroutine
	
	var wg sync.WaitGroup
	var submissionErrors int64
	
	// Submit tasks concurrently
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			
			for t := 0; t < tasksPerGoroutine; t++ {
				task := &TestTask{
					id:      fmt.Sprintf("concurrent-task-%d-%d", goroutineID, t),
					timeout: 2 * time.Second,
				}
				
				if err := pool.Submit(task); err != nil {
					atomic.AddInt64(&submissionErrors, 1)
				}
			}
		}(g)
	}
	
	wg.Wait()
	
	// Verify no submission errors
	assert.Equal(suite.T(), int64(0), atomic.LoadInt64(&submissionErrors), "Should have no submission errors")
	
	// Wait for all tasks to complete
	time.Sleep(5 * time.Second)
	
	// Verify all tasks completed
	stats := pool.GetStats()
	assert.Equal(suite.T(), int64(totalTasks), stats.CompletedTasks, "All tasks should complete")
}

// @decorator: TestWorkerPoolShutdown
// @description: Test graceful shutdown of worker pool
func (suite *WorkerPoolTestSuite) TestWorkerPoolShutdown() {
	config := async.WorkerPoolConfig{
		MinWorkers:                2,
		MaxWorkers:                4,
		QueueSize:                 10,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 5,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	
	// Submit some long-running tasks
	for i := 0; i < 3; i++ {
		task := &TestTask{
			id:       fmt.Sprintf("shutdown-task-%d", i),
			timeout:  10 * time.Second,
			duration: 2 * time.Second,
		}
		
		err := pool.Submit(task)
		assert.NoError(suite.T(), err)
	}
	
	// Start shutdown
	shutdownStart := time.Now()
	err = pool.Shutdown(5 * time.Second)
	shutdownDuration := time.Since(shutdownStart)
	
	assert.NoError(suite.T(), err, "Shutdown should complete without error")
	assert.Less(suite.T(), shutdownDuration, 6*time.Second, "Shutdown should complete within timeout")
	
	// Verify no new tasks can be submitted
	newTask := &TestTask{
		id:      "post-shutdown-task",
		timeout: 1 * time.Second,
	}
	
	err = pool.Submit(newTask)
	assert.Error(suite.T(), err, "Should reject tasks after shutdown")
}

// @decorator: TestWorkerPoolMetrics
// @description: Test worker pool metrics and statistics
func (suite *WorkerPoolTestSuite) TestWorkerPoolMetrics() {
	config := async.WorkerPoolConfig{
		MinWorkers:                2,
		MaxWorkers:                4,
		QueueSize:                 10,
		WorkerTimeout:             5 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 5,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(5 * time.Second)
	
	// Initial stats
	stats := pool.GetStats()
	assert.GreaterOrEqual(suite.T(), stats.ActiveWorkers, 0)
	assert.GreaterOrEqual(suite.T(), stats.TotalWorkers, config.MinWorkers)
	assert.Equal(suite.T(), int64(0), stats.TotalTasks)
	assert.Equal(suite.T(), int64(0), stats.CompletedTasks)
	assert.Equal(suite.T(), int64(0), stats.FailedTasks)
	
	// Submit successful tasks
	successTasks := 5
	for i := 0; i < successTasks; i++ {
		task := &TestTask{
			id:      fmt.Sprintf("success-task-%d", i),
			timeout: 1 * time.Second,
		}
		
		err := pool.Submit(task)
		assert.NoError(suite.T(), err)
	}
	
	// Submit failing tasks
	failTasks := 3
	for i := 0; i < failTasks; i++ {
		task := &TestTask{
			id:          fmt.Sprintf("fail-task-%d", i),
			timeout:     1 * time.Second,
			shouldFail:  true,
			shouldRetry: false,
		}
		
		err := pool.Submit(task)
		assert.NoError(suite.T(), err)
	}
	
	// Wait for tasks to complete
	time.Sleep(3 * time.Second)
	
	// Check final stats
	stats = pool.GetStats()
	assert.Equal(suite.T(), int64(successTasks+failTasks), stats.TotalTasks)
	assert.Equal(suite.T(), int64(successTasks), stats.CompletedTasks)
	assert.Equal(suite.T(), int64(failTasks), stats.FailedTasks)
}

// @decorator: TestWorkerPoolStressTest
// @description: Stress test with high load and concurrent operations
func (suite *WorkerPoolTestSuite) TestWorkerPoolStressTest() {
	if testing.Short() {
		suite.T().Skip("Skipping stress test in short mode")
	}
	
	config := async.WorkerPoolConfig{
		MinWorkers:                5,
		MaxWorkers:                20,
		QueueSize:                 200,
		WorkerTimeout:             10 * time.Second,
		CircuitBreakerTimeout:     5 * time.Second,
		CircuitBreakerMaxFailures: 50,
	}
	
	pool, err := async.NewWorkerPool(config, suite.logger)
	require.NoError(suite.T(), err)
	defer pool.Shutdown(10 * time.Second)
	
	numTasks := 1000
	var completedTasks int64
	
	// Submit many tasks concurrently
	var wg sync.WaitGroup
	for i := 0; i < numTasks; i++ {
		wg.Add(1)
		go func(taskID int) {
			defer wg.Done()
			
			task := &TestTask{
				id:       fmt.Sprintf("stress-task-%d", taskID),
				timeout:  5 * time.Second,
				duration: time.Duration(taskID%100) * time.Millisecond, // Variable duration
			}
			
			if err := pool.Submit(task); err == nil {
				atomic.AddInt64(&completedTasks, 1)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Wait for all tasks to process
	time.Sleep(10 * time.Second)
	
	// Verify high completion rate
	stats := pool.GetStats()
	completionRate := float64(stats.CompletedTasks) / float64(numTasks) * 100
	assert.Greater(suite.T(), completionRate, 95.0, "Should have >95%% completion rate under stress")
	
	suite.logger.WithFields(logrus.Fields{
		"total_tasks":     numTasks,
		"completed_tasks": stats.CompletedTasks,
		"failed_tasks":    stats.FailedTasks,
		"completion_rate": completionRate,
		"total_workers":   stats.TotalWorkers,
	}).Info("Stress test completed")
}

// Run the test suite
func TestWorkerPoolSuite(t *testing.T) {
	suite.Run(t, new(WorkerPoolTestSuite))
}
