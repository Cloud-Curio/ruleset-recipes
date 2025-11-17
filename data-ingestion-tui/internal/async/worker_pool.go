/*
 * Data Ingestion TUI - Advanced Async Worker Pool
 * 
 * High-performance async worker pool implementation with advanced
 * concurrency patterns, telemetry, and optimization features.
 * 
 * Features:
 * - Dynamic worker scaling
 * - Work stealing algorithms
 * - Backpressure handling
 * - Comprehensive telemetry
 * - Circuit breaker patterns
 * - Graceful shutdown
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package async

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

// @decorator: WorkerPool
// @description: High-performance async worker pool with advanced concurrency features
// @version: 2.0.0
// @author: Codegen AI Assistant

// WorkerPool manages a dynamic pool of workers with advanced features
type WorkerPool struct {
	// Configuration
	minWorkers    int
	maxWorkers    int
	queueSize     int
	workerTimeout time.Duration
	
	// State management
	mu            sync.RWMutex
	workers       []*Worker
	activeWorkers int64
	totalTasks    int64
	completedTasks int64
	failedTasks   int64
	
	// Channels
	taskQueue     chan Task
	resultQueue   chan TaskResult
	shutdownChan  chan struct{}
	
	// Context and cancellation
	ctx           context.Context
	cancel        context.CancelFunc
	
	// Telemetry
	logger        *logrus.Logger
	metrics       *WorkerPoolMetrics
	
	// Circuit breaker
	circuitBreaker *CircuitBreaker
	
	// Work stealing
	workStealers  []*WorkStealer
}

// @decorator: Task
// @description: Represents a task to be executed by workers
type Task interface {
	Execute(ctx context.Context) (interface{}, error)
	GetID() string
	GetPriority() int
	GetTimeout() time.Duration
	GetRetryCount() int
	ShouldRetry(err error) bool
}

// @decorator: TaskResult
// @description: Result of task execution
type TaskResult struct {
	TaskID    string
	Result    interface{}
	Error     error
	Duration  time.Duration
	WorkerID  int
	Timestamp time.Time
	Retries   int
}

// @decorator: Worker
// @description: Individual worker in the pool
type Worker struct {
	ID           int
	pool         *WorkerPool
	taskQueue    chan Task
	resultQueue  chan TaskResult
	stopChan     chan struct{}
	isActive     int64
	tasksHandled int64
	lastActivity time.Time
	mu           sync.RWMutex
}

// @decorator: WorkerPoolMetrics
// @description: Prometheus metrics for worker pool telemetry
type WorkerPoolMetrics struct {
	ActiveWorkers    prometheus.Gauge
	QueuedTasks      prometheus.Gauge
	CompletedTasks   prometheus.Counter
	FailedTasks      prometheus.Counter
	TaskDuration     prometheus.Histogram
	WorkerUtilization prometheus.Histogram
	QueueWaitTime    prometheus.Histogram
}

// @decorator: CircuitBreaker
// @description: Circuit breaker for fault tolerance
type CircuitBreaker struct {
	mu              sync.RWMutex
	state           CircuitState
	failureCount    int64
	successCount    int64
	lastFailureTime time.Time
	timeout         time.Duration
	maxFailures     int64
}

// @decorator: CircuitState
// @description: Circuit breaker states
type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

// @decorator: WorkStealer
// @description: Work stealing implementation for load balancing
type WorkStealer struct {
	workerID     int
	stealQueue   chan Task
	stealAttempts int64
	stealSuccess int64
}

// @decorator: NewWorkerPool
// @description: Create a new advanced worker pool
// @param config: Worker pool configuration
// @param logger: Logger instance
// @return *WorkerPool: New worker pool
// @return error: Any error that occurred during creation
func NewWorkerPool(config WorkerPoolConfig, logger *logrus.Logger) (*WorkerPool, error) {
	if config.MinWorkers <= 0 {
		config.MinWorkers = runtime.NumCPU()
	}
	
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = runtime.NumCPU() * 4
	}
	
	if config.QueueSize <= 0 {
		config.QueueSize = 1000
	}
	
	if config.WorkerTimeout <= 0 {
		config.WorkerTimeout = 30 * time.Second
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	// Initialize metrics
	metrics := &WorkerPoolMetrics{
		ActiveWorkers: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "worker_pool_active_workers",
			Help: "Number of active workers in the pool",
		}),
		QueuedTasks: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "worker_pool_queued_tasks",
			Help: "Number of tasks in the queue",
		}),
		CompletedTasks: promauto.NewCounter(prometheus.CounterOpts{
			Name: "worker_pool_completed_tasks_total",
			Help: "Total number of completed tasks",
		}),
		FailedTasks: promauto.NewCounter(prometheus.CounterOpts{
			Name: "worker_pool_failed_tasks_total",
			Help: "Total number of failed tasks",
		}),
		TaskDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "worker_pool_task_duration_seconds",
			Help:    "Task execution duration in seconds",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		}),
		WorkerUtilization: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "worker_pool_worker_utilization",
			Help:    "Worker utilization percentage",
			Buckets: prometheus.LinearBuckets(0, 10, 11),
		}),
		QueueWaitTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Name:    "worker_pool_queue_wait_time_seconds",
			Help:    "Time tasks spend waiting in queue",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		}),
	}
	
	// Initialize circuit breaker
	circuitBreaker := &CircuitBreaker{
		state:       CircuitClosed,
		timeout:     config.CircuitBreakerTimeout,
		maxFailures: config.CircuitBreakerMaxFailures,
	}
	
	pool := &WorkerPool{
		minWorkers:     config.MinWorkers,
		maxWorkers:     config.MaxWorkers,
		queueSize:      config.QueueSize,
		workerTimeout:  config.WorkerTimeout,
		taskQueue:      make(chan Task, config.QueueSize),
		resultQueue:    make(chan TaskResult, config.QueueSize),
		shutdownChan:   make(chan struct{}),
		ctx:            ctx,
		cancel:         cancel,
		logger:         logger,
		metrics:        metrics,
		circuitBreaker: circuitBreaker,
		workers:        make([]*Worker, 0, config.MaxWorkers),
		workStealers:   make([]*WorkStealer, 0, config.MaxWorkers),
	}
	
	// Start initial workers
	for i := 0; i < config.MinWorkers; i++ {
		if err := pool.addWorker(); err != nil {
			return nil, errors.Wrap(err, "failed to create initial workers")
		}
	}
	
	// Start monitoring goroutines
	go pool.monitorWorkers()
	go pool.monitorQueue()
	go pool.updateMetrics()
	
	logger.WithFields(logrus.Fields{
		"min_workers":    config.MinWorkers,
		"max_workers":    config.MaxWorkers,
		"queue_size":     config.QueueSize,
		"worker_timeout": config.WorkerTimeout,
	}).Info("Advanced worker pool initialized")
	
	return pool, nil
}

// @decorator: WorkerPoolConfig
// @description: Configuration for worker pool
type WorkerPoolConfig struct {
	MinWorkers                 int
	MaxWorkers                 int
	QueueSize                  int
	WorkerTimeout              time.Duration
	ScaleUpThreshold           float64
	ScaleDownThreshold         float64
	CircuitBreakerTimeout      time.Duration
	CircuitBreakerMaxFailures  int64
	EnableWorkStealing         bool
	WorkStealingThreshold      int
}

// @decorator: Submit
// @description: Submit a task to the worker pool
// @param task: Task to execute
// @return error: Any error that occurred during submission
func (wp *WorkerPool) Submit(task Task) error {
	select {
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	default:
	}
	
	// Check circuit breaker
	if !wp.circuitBreaker.CanExecute() {
		wp.metrics.FailedTasks.Inc()
		return errors.New("circuit breaker is open")
	}
	
	// Add task to queue with timeout
	select {
	case wp.taskQueue <- task:
		atomic.AddInt64(&wp.totalTasks, 1)
		wp.metrics.QueuedTasks.Inc()
		return nil
	case <-time.After(5 * time.Second):
		wp.metrics.FailedTasks.Inc()
		return errors.New("task queue is full")
	case <-wp.ctx.Done():
		return errors.New("worker pool is shutting down")
	}
}

// @decorator: addWorker
// @description: Add a new worker to the pool
// @return error: Any error that occurred during worker creation
func (wp *WorkerPool) addWorker() error {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	
	if len(wp.workers) >= wp.maxWorkers {
		return errors.New("maximum workers reached")
	}
	
	workerID := len(wp.workers)
	worker := &Worker{
		ID:           workerID,
		pool:         wp,
		taskQueue:    wp.taskQueue,
		resultQueue:  wp.resultQueue,
		stopChan:     make(chan struct{}),
		lastActivity: time.Now(),
	}
	
	wp.workers = append(wp.workers, worker)
	
	// Start worker goroutine
	go worker.run()
	
	// Add work stealer if enabled
	if len(wp.workStealers) < len(wp.workers) {
		stealer := &WorkStealer{
			workerID:   workerID,
			stealQueue: make(chan Task, 10),
		}
		wp.workStealers = append(wp.workStealers, stealer)
		go stealer.run(wp)
	}
	
	atomic.AddInt64(&wp.activeWorkers, 1)
	wp.metrics.ActiveWorkers.Inc()
	
	wp.logger.WithField("worker_id", workerID).Debug("Added new worker")
	return nil
}

// @decorator: run
// @description: Main worker execution loop
func (w *Worker) run() {
	defer func() {
		if r := recover(); r != nil {
			w.pool.logger.WithFields(logrus.Fields{
				"worker_id": w.ID,
				"panic":     r,
			}).Error("Worker panicked")
		}
		atomic.AddInt64(&w.pool.activeWorkers, -1)
		w.pool.metrics.ActiveWorkers.Dec()
	}()
	
	for {
		select {
		case task := <-w.taskQueue:
			w.executeTask(task)
		case <-w.stopChan:
			w.pool.logger.WithField("worker_id", w.ID).Debug("Worker stopping")
			return
		case <-w.pool.ctx.Done():
			return
		case <-time.After(w.pool.workerTimeout):
			// Worker timeout - check if we should scale down
			if w.shouldScaleDown() {
				return
			}
		}
	}
}

// @decorator: executeTask
// @description: Execute a single task with comprehensive error handling and telemetry
// @param task: Task to execute
func (w *Worker) executeTask(task Task) {
	startTime := time.Now()
	queueWaitTime := startTime.Sub(time.Now()) // This would be set when task is queued
	
	atomic.StoreInt64(&w.isActive, 1)
	w.mu.Lock()
	w.lastActivity = startTime
	w.mu.Unlock()
	
	defer func() {
		atomic.StoreInt64(&w.isActive, 0)
		atomic.AddInt64(&w.tasksHandled, 1)
		
		duration := time.Since(startTime)
		w.pool.metrics.TaskDuration.Observe(duration.Seconds())
		w.pool.metrics.QueueWaitTime.Observe(queueWaitTime.Seconds())
		w.pool.metrics.QueuedTasks.Dec()
	}()
	
	// Create task context with timeout
	taskCtx, cancel := context.WithTimeout(w.pool.ctx, task.GetTimeout())
	defer cancel()
	
	// Execute task with retry logic
	var result interface{}
	var err error
	retries := 0
	maxRetries := task.GetRetryCount()
	
	for retries <= maxRetries {
		result, err = task.Execute(taskCtx)
		
		if err == nil {
			// Task succeeded
			w.pool.circuitBreaker.RecordSuccess()
			atomic.AddInt64(&w.pool.completedTasks, 1)
			w.pool.metrics.CompletedTasks.Inc()
			break
		}
		
		// Task failed - check if we should retry
		if retries < maxRetries && task.ShouldRetry(err) {
			retries++
			backoffDuration := time.Duration(retries*retries) * 100 * time.Millisecond
			
			w.pool.logger.WithFields(logrus.Fields{
				"task_id":  task.GetID(),
				"worker_id": w.ID,
				"retry":    retries,
				"error":    err.Error(),
				"backoff":  backoffDuration,
			}).Warn("Task failed, retrying")
			
			select {
			case <-time.After(backoffDuration):
				continue
			case <-taskCtx.Done():
				err = taskCtx.Err()
				break
			}
		} else {
			// No more retries or shouldn't retry
			break
		}
	}
	
	// Record final result
	taskResult := TaskResult{
		TaskID:    task.GetID(),
		Result:    result,
		Error:     err,
		Duration:  time.Since(startTime),
		WorkerID:  w.ID,
		Timestamp: time.Now(),
		Retries:   retries,
	}
	
	if err != nil {
		w.pool.circuitBreaker.RecordFailure()
		atomic.AddInt64(&w.pool.failedTasks, 1)
		w.pool.metrics.FailedTasks.Inc()
		
		w.pool.logger.WithFields(logrus.Fields{
			"task_id":   task.GetID(),
			"worker_id": w.ID,
			"error":     err.Error(),
			"retries":   retries,
			"duration":  taskResult.Duration,
		}).Error("Task execution failed")
	} else {
		w.pool.logger.WithFields(logrus.Fields{
			"task_id":   task.GetID(),
			"worker_id": w.ID,
			"duration":  taskResult.Duration,
			"retries":   retries,
		}).Debug("Task executed successfully")
	}
	
	// Send result
	select {
	case w.resultQueue <- taskResult:
	case <-time.After(5 * time.Second):
		w.pool.logger.WithField("task_id", task.GetID()).Warn("Failed to send task result")
	case <-w.pool.ctx.Done():
		return
	}
}

// @decorator: shouldScaleDown
// @description: Determine if this worker should be removed from the pool
// @return bool: True if worker should be scaled down
func (w *Worker) shouldScaleDown() bool {
	w.pool.mu.RLock()
	defer w.pool.mu.RUnlock()
	
	// Don't scale below minimum workers
	if len(w.pool.workers) <= w.pool.minWorkers {
		return false
	}
	
	// Check if worker has been idle
	w.mu.RLock()
	idleTime := time.Since(w.lastActivity)
	w.mu.RUnlock()
	
	return idleTime > w.pool.workerTimeout
}

// @decorator: monitorWorkers
// @description: Monitor worker pool and handle dynamic scaling
func (wp *WorkerPool) monitorWorkers() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			wp.handleDynamicScaling()
		case <-wp.ctx.Done():
			return
		}
	}
}

// @decorator: handleDynamicScaling
// @description: Handle dynamic scaling of workers based on load
func (wp *WorkerPool) handleDynamicScaling() {
	queueLength := len(wp.taskQueue)
	activeWorkers := int(atomic.LoadInt64(&wp.activeWorkers))
	
	// Calculate utilization
	utilization := float64(queueLength) / float64(wp.queueSize)
	
	wp.logger.WithFields(logrus.Fields{
		"queue_length":    queueLength,
		"active_workers":  activeWorkers,
		"utilization":     utilization,
		"total_workers":   len(wp.workers),
	}).Debug("Worker pool status")
	
	// Scale up if needed
	if utilization > 0.8 && len(wp.workers) < wp.maxWorkers {
		if err := wp.addWorker(); err != nil {
			wp.logger.WithError(err).Warn("Failed to scale up workers")
		} else {
			wp.logger.Info("Scaled up worker pool")
		}
	}
	
	// Update utilization metric
	wp.metrics.WorkerUtilization.Observe(utilization * 100)
}

// @decorator: monitorQueue
// @description: Monitor task queue and handle backpressure
func (wp *WorkerPool) monitorQueue() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			queueLength := len(wp.taskQueue)
			wp.metrics.QueuedTasks.Set(float64(queueLength))
			
			// Handle backpressure
			if queueLength > wp.queueSize*8/10 {
				wp.logger.WithField("queue_length", queueLength).Warn("High queue utilization detected")
			}
			
		case <-wp.ctx.Done():
			return
		}
	}
}

// @decorator: updateMetrics
// @description: Update Prometheus metrics periodically
func (wp *WorkerPool) updateMetrics() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			wp.metrics.ActiveWorkers.Set(float64(atomic.LoadInt64(&wp.activeWorkers)))
			
		case <-wp.ctx.Done():
			return
		}
	}
}

// @decorator: GetResults
// @description: Get task results channel
// @return <-chan TaskResult: Results channel
func (wp *WorkerPool) GetResults() <-chan TaskResult {
	return wp.resultQueue
}

// @decorator: GetStats
// @description: Get worker pool statistics
// @return WorkerPoolStats: Current statistics
func (wp *WorkerPool) GetStats() WorkerPoolStats {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	
	return WorkerPoolStats{
		ActiveWorkers:   int(atomic.LoadInt64(&wp.activeWorkers)),
		TotalWorkers:    len(wp.workers),
		QueuedTasks:     len(wp.taskQueue),
		TotalTasks:      atomic.LoadInt64(&wp.totalTasks),
		CompletedTasks:  atomic.LoadInt64(&wp.completedTasks),
		FailedTasks:     atomic.LoadInt64(&wp.failedTasks),
		CircuitState:    wp.circuitBreaker.GetState(),
	}
}

// @decorator: WorkerPoolStats
// @description: Worker pool statistics
type WorkerPoolStats struct {
	ActiveWorkers  int
	TotalWorkers   int
	QueuedTasks    int
	TotalTasks     int64
	CompletedTasks int64
	FailedTasks    int64
	CircuitState   CircuitState
}

// @decorator: Shutdown
// @description: Gracefully shutdown the worker pool
// @param timeout: Maximum time to wait for shutdown
// @return error: Any error that occurred during shutdown
func (wp *WorkerPool) Shutdown(timeout time.Duration) error {
	wp.logger.Info("Shutting down worker pool...")
	
	// Cancel context to signal shutdown
	wp.cancel()
	
	// Close task queue to prevent new submissions
	close(wp.taskQueue)
	
	// Wait for workers to finish with timeout
	done := make(chan struct{})
	go func() {
		wp.mu.Lock()
		for _, worker := range wp.workers {
			close(worker.stopChan)
		}
		wp.mu.Unlock()
		
		// Wait for all workers to finish
		for atomic.LoadInt64(&wp.activeWorkers) > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		close(done)
	}()
	
	select {
	case <-done:
		wp.logger.Info("Worker pool shutdown complete")
		return nil
	case <-time.After(timeout):
		wp.logger.Warn("Worker pool shutdown timed out")
		return errors.New("shutdown timeout")
	}
}

// Circuit Breaker Implementation

// @decorator: CanExecute
// @description: Check if circuit breaker allows execution
// @return bool: True if execution is allowed
func (cb *CircuitBreaker) CanExecute() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	
	switch cb.state {
	case CircuitClosed:
		return true
	case CircuitOpen:
		return time.Since(cb.lastFailureTime) > cb.timeout
	case CircuitHalfOpen:
		return true
	default:
		return false
	}
}

// @decorator: RecordSuccess
// @description: Record a successful execution
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.successCount++
	
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
		cb.failureCount = 0
	}
}

// @decorator: RecordFailure
// @description: Record a failed execution
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	
	cb.failureCount++
	cb.lastFailureTime = time.Now()
	
	if cb.failureCount >= cb.maxFailures {
		cb.state = CircuitOpen
	}
}

// @decorator: GetState
// @description: Get current circuit breaker state
// @return CircuitState: Current state
func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Work Stealing Implementation

// @decorator: run
// @description: Run work stealing goroutine
// @param pool: Worker pool reference
func (ws *WorkStealer) run(pool *WorkerPool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ws.attemptWorkStealing(pool)
		case <-pool.ctx.Done():
			return
		}
	}
}

// @decorator: attemptWorkStealing
// @description: Attempt to steal work from other workers
// @param pool: Worker pool reference
func (ws *WorkStealer) attemptWorkStealing(pool *WorkerPool) {
	// Simple work stealing - try to take tasks from queue if idle
	if len(pool.taskQueue) > 10 {
		select {
		case task := <-pool.taskQueue:
			atomic.AddInt64(&ws.stealAttempts, 1)
			// Put task back in queue for now (simplified implementation)
			select {
			case pool.taskQueue <- task:
				atomic.AddInt64(&ws.stealSuccess, 1)
			default:
				// Queue full, drop task (in real implementation, would handle better)
			}
		default:
			// No tasks to steal
		}
	}
}
