/*
 * Data Ingestion TUI - Bulk Ingestion Engine
 * 
 * This package provides comprehensive bulk data ingestion capabilities
 * with support for multiple data sources, parallel processing, error
 * handling, and data validation.
 * 
 * Features:
 * - Parallel bulk ingestion with worker pools
 * - Rate limiting and retry logic
 * - Data validation and transformation
 * - Progress tracking and monitoring
 * - Comprehensive error handling
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"

	"data-ingestion-tui/internal/config"
	"data-ingestion-tui/internal/database"
	"data-ingestion-tui/internal/models"
)

// @decorator: Engine
// @description: Main bulk ingestion engine with parallel processing
// @version: 1.0.0
// @author: Codegen AI Assistant

// Engine manages bulk data ingestion operations
type Engine struct {
	config     *config.Config
	db         *database.DB
	logger     *logrus.Logger
	
	// Worker pool management
	workerPool   *WorkerPool
	rateLimiter  *rate.Limiter
	
	// Progress tracking
	progressChan chan *IngestionProgress
	errorChan    chan error
	
	// State management
	mu       sync.RWMutex
	running  bool
	jobs     map[string]*IngestionJob
	
	// Data processors
	processors map[string]DataProcessor
	validators map[string]DataValidator
}

// @decorator: IngestionJob
// @description: Represents a bulk ingestion job
type IngestionJob struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Endpoint    string                 `json:"endpoint"`
	Status      JobStatus              `json:"status"`
	Progress    *IngestionProgress     `json:"progress"`
	Config      *JobConfig             `json:"config"`
	StartTime   time.Time              `json:"start_time"`
	EndTime     *time.Time             `json:"end_time,omitempty"`
	Error       string                 `json:"error,omitempty"`
	Results     *IngestionResults      `json:"results,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// @decorator: JobStatus
// @description: Status enumeration for ingestion jobs
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
	JobStatusPaused    JobStatus = "paused"
)

// @decorator: JobConfig
// @description: Configuration for bulk ingestion jobs
type JobConfig struct {
	// Data source configuration
	Endpoint     string            `json:"endpoint"`
	APIKey       string            `json:"api_key,omitempty"`
	BaseURL      string            `json:"base_url"`
	Headers      map[string]string `json:"headers,omitempty"`
	
	// Bulk processing configuration
	BatchSize    int               `json:"batch_size"`
	MaxWorkers   int               `json:"max_workers"`
	RateLimit    int               `json:"rate_limit"` // requests per second
	Timeout      time.Duration     `json:"timeout"`
	
	// Data processing configuration
	DataFormat   string            `json:"data_format"` // json, xml, csv
	Validation   bool              `json:"validation"`
	Transform    bool              `json:"transform"`
	Deduplicate  bool              `json:"deduplicate"`
	
	// Storage configuration
	TableName    string            `json:"table_name"`
	CreateTable  bool              `json:"create_table"`
	TruncateFirst bool             `json:"truncate_first"`
	
	// Retry configuration
	MaxRetries   int               `json:"max_retries"`
	RetryDelay   time.Duration     `json:"retry_delay"`
	
	// Custom parameters
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

// @decorator: IngestionProgress
// @description: Progress tracking for ingestion jobs
type IngestionProgress struct {
	JobID           string    `json:"job_id"`
	TotalRecords    int64     `json:"total_records"`
	ProcessedRecords int64    `json:"processed_records"`
	SuccessfulRecords int64   `json:"successful_records"`
	FailedRecords   int64     `json:"failed_records"`
	CurrentBatch    int       `json:"current_batch"`
	TotalBatches    int       `json:"total_batches"`
	RecordsPerSecond float64  `json:"records_per_second"`
	EstimatedTimeRemaining time.Duration `json:"estimated_time_remaining"`
	LastUpdate      time.Time `json:"last_update"`
}

// @decorator: IngestionResults
// @description: Final results of an ingestion job
type IngestionResults struct {
	TotalRecords      int64         `json:"total_records"`
	SuccessfulRecords int64         `json:"successful_records"`
	FailedRecords     int64         `json:"failed_records"`
	DuplicateRecords  int64         `json:"duplicate_records"`
	Duration          time.Duration `json:"duration"`
	AverageRPS        float64       `json:"average_rps"`
	ErrorSummary      []ErrorSummary `json:"error_summary,omitempty"`
	DataQualityReport *DataQualityReport `json:"data_quality_report,omitempty"`
}

// @decorator: ErrorSummary
// @description: Summary of errors encountered during ingestion
type ErrorSummary struct {
	ErrorType    string `json:"error_type"`
	Count        int64  `json:"count"`
	SampleError  string `json:"sample_error"`
	FirstOccurrence time.Time `json:"first_occurrence"`
}

// @decorator: DataQualityReport
// @description: Data quality metrics for ingested data
type DataQualityReport struct {
	Completeness    float64            `json:"completeness"`
	Accuracy        float64            `json:"accuracy"`
	Consistency     float64            `json:"consistency"`
	Validity        float64            `json:"validity"`
	FieldStatistics map[string]FieldStats `json:"field_statistics"`
}

// @decorator: FieldStats
// @description: Statistics for individual fields
type FieldStats struct {
	TotalValues   int64   `json:"total_values"`
	NullValues    int64   `json:"null_values"`
	UniqueValues  int64   `json:"unique_values"`
	MinLength     int     `json:"min_length,omitempty"`
	MaxLength     int     `json:"max_length,omitempty"`
	AvgLength     float64 `json:"avg_length,omitempty"`
}

// @decorator: New
// @description: Create a new ingestion engine
// @param config: Application configuration
// @param db: Database connection
// @param logger: Logger instance
// @return *Engine: New ingestion engine
// @return error: Any error that occurred during creation
func New(config *config.Config, db *database.DB, logger *logrus.Logger) (*Engine, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}
	
	if db == nil {
		return nil, errors.New("database connection cannot be nil")
	}
	
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}
	
	// Create rate limiter (default: 10 requests per second)
	rateLimiter := rate.NewLimiter(rate.Limit(10), 20)
	
	// Create worker pool
	workerPool := NewWorkerPool(5, logger) // Default: 5 workers
	
	engine := &Engine{
		config:       config,
		db:          db,
		logger:      logger,
		workerPool:  workerPool,
		rateLimiter: rateLimiter,
		progressChan: make(chan *IngestionProgress, 100),
		errorChan:   make(chan error, 100),
		jobs:        make(map[string]*IngestionJob),
		processors:  make(map[string]DataProcessor),
		validators:  make(map[string]DataValidator),
	}
	
	// Initialize built-in processors and validators
	if err := engine.initializeProcessors(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize data processors")
	}
	
	if err := engine.initializeValidators(); err != nil {
		return nil, errors.Wrap(err, "failed to initialize data validators")
	}
	
	logger.Info("Bulk ingestion engine initialized successfully")
	return engine, nil
}

// @decorator: StartBulkIngestion
// @description: Start a bulk ingestion job
// @param ctx: Context for cancellation
// @param jobConfig: Job configuration
// @return string: Job ID
// @return error: Any error that occurred during job start
func (e *Engine) StartBulkIngestion(ctx context.Context, jobConfig *JobConfig) (string, error) {
	if jobConfig == nil {
		return "", errors.New("job configuration cannot be nil")
	}
	
	// Validate job configuration
	if err := e.validateJobConfig(jobConfig); err != nil {
		return "", errors.Wrap(err, "invalid job configuration")
	}
	
	// Create new job
	job := &IngestionJob{
		ID:        generateJobID(),
		Name:      fmt.Sprintf("Bulk ingestion from %s", jobConfig.Endpoint),
		Endpoint:  jobConfig.Endpoint,
		Status:    JobStatusPending,
		Config:    jobConfig,
		StartTime: time.Now(),
		Progress: &IngestionProgress{
			LastUpdate: time.Now(),
		},
		Metadata: make(map[string]interface{}),
	}
	
	// Store job
	e.mu.Lock()
	e.jobs[job.ID] = job
	e.mu.Unlock()
	
	// Start job asynchronously
	go e.executeJob(ctx, job)
	
	e.logger.WithFields(logrus.Fields{
		"job_id":   job.ID,
		"endpoint": jobConfig.Endpoint,
	}).Info("Started bulk ingestion job")
	
	return job.ID, nil
}

// @decorator: executeJob
// @description: Execute a bulk ingestion job
// @param ctx: Context for cancellation
// @param job: Job to execute
func (e *Engine) executeJob(ctx context.Context, job *IngestionJob) {
	defer func() {
		if r := recover(); r != nil {
			e.logger.WithField("job_id", job.ID).Errorf("Job panicked: %v", r)
			job.Status = JobStatusFailed
			job.Error = fmt.Sprintf("Job panicked: %v", r)
			endTime := time.Now()
			job.EndTime = &endTime
		}
	}()
	
	// Update job status
	job.Status = JobStatusRunning
	job.Progress.JobID = job.ID
	
	e.logger.WithField("job_id", job.ID).Info("Executing bulk ingestion job")
	
	// Get data source
	dataSource, err := e.getDataSource(job.Config.Endpoint)
	if err != nil {
		e.failJob(job, errors.Wrap(err, "failed to get data source"))
		return
	}
	
	// Initialize progress tracking
	startTime := time.Now()
	
	// Execute ingestion based on data source type
	var results *IngestionResults
	switch dataSource.Type {
	case "rest_api":
		results, err = e.executeRESTIngestion(ctx, job, dataSource)
	case "graphql":
		results, err = e.executeGraphQLIngestion(ctx, job, dataSource)
	case "file":
		results, err = e.executeFileIngestion(ctx, job, dataSource)
	default:
		err = errors.Errorf("unsupported data source type: %s", dataSource.Type)
	}
	
	// Handle results
	if err != nil {
		e.failJob(job, err)
		return
	}
	
	// Complete job
	endTime := time.Now()
	job.EndTime = &endTime
	job.Status = JobStatusCompleted
	job.Results = results
	job.Results.Duration = endTime.Sub(startTime)
	
	e.logger.WithFields(logrus.Fields{
		"job_id":            job.ID,
		"duration":          job.Results.Duration,
		"total_records":     job.Results.TotalRecords,
		"successful_records": job.Results.SuccessfulRecords,
		"failed_records":    job.Results.FailedRecords,
	}).Info("Bulk ingestion job completed successfully")
}

// @decorator: executeRESTIngestion
// @description: Execute REST API bulk ingestion
// @param ctx: Context for cancellation
// @param job: Job to execute
// @param dataSource: Data source configuration
// @return *IngestionResults: Ingestion results
// @return error: Any error that occurred during ingestion
func (e *Engine) executeRESTIngestion(ctx context.Context, job *IngestionJob, dataSource *DataSource) (*IngestionResults, error) {
	results := &IngestionResults{
		ErrorSummary: make([]ErrorSummary, 0),
	}
	
	// Create HTTP client with timeout
	client := NewHTTPClient(job.Config.Timeout)
	
	// Set up rate limiter for this job
	jobRateLimiter := rate.NewLimiter(rate.Limit(job.Config.RateLimit), job.Config.RateLimit*2)
	
	// Get total record count (if supported by API)
	totalRecords, err := e.getTotalRecordCount(client, job, dataSource)
	if err != nil {
		e.logger.WithError(err).Warn("Could not determine total record count")
		totalRecords = -1 // Unknown
	}
	
	job.Progress.TotalRecords = totalRecords
	if totalRecords > 0 {
		job.Progress.TotalBatches = int((totalRecords + int64(job.Config.BatchSize) - 1) / int64(job.Config.BatchSize))
	}
	
	// Process data in batches
	offset := 0
	batchNum := 0
	
	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		default:
		}
		
		// Rate limiting
		if err := jobRateLimiter.Wait(ctx); err != nil {
			return results, err
		}
		
		// Fetch batch
		batchData, hasMore, err := e.fetchBatch(ctx, client, job, dataSource, offset, job.Config.BatchSize)
		if err != nil {
			e.logger.WithError(err).WithFields(logrus.Fields{
				"job_id": job.ID,
				"batch":  batchNum,
				"offset": offset,
			}).Error("Failed to fetch batch")
			
			// Add to error summary
			e.addToErrorSummary(results, "fetch_error", err.Error())
			
			// Check if we should retry or fail
			if job.Config.MaxRetries > 0 {
				time.Sleep(job.Config.RetryDelay)
				job.Config.MaxRetries--
				continue
			} else {
				return results, err
			}
		}
		
		if len(batchData) == 0 {
			break
		}
		
		// Process batch
		batchResults, err := e.processBatch(ctx, job, batchData, batchNum)
		if err != nil {
			e.logger.WithError(err).WithField("job_id", job.ID).Error("Failed to process batch")
			e.addToErrorSummary(results, "process_error", err.Error())
		} else {
			// Update results
			results.TotalRecords += batchResults.TotalRecords
			results.SuccessfulRecords += batchResults.SuccessfulRecords
			results.FailedRecords += batchResults.FailedRecords
			results.DuplicateRecords += batchResults.DuplicateRecords
		}
		
		// Update progress
		batchNum++
		offset += job.Config.BatchSize
		job.Progress.CurrentBatch = batchNum
		job.Progress.ProcessedRecords = results.TotalRecords
		job.Progress.SuccessfulRecords = results.SuccessfulRecords
		job.Progress.FailedRecords = results.FailedRecords
		job.Progress.LastUpdate = time.Now()
		
		// Calculate records per second
		elapsed := time.Since(job.StartTime)
		if elapsed > 0 {
			job.Progress.RecordsPerSecond = float64(results.TotalRecords) / elapsed.Seconds()
			
			// Estimate time remaining
			if totalRecords > 0 && job.Progress.RecordsPerSecond > 0 {
				remaining := totalRecords - results.TotalRecords
				job.Progress.EstimatedTimeRemaining = time.Duration(float64(remaining)/job.Progress.RecordsPerSecond) * time.Second
			}
		}
		
		// Send progress update
		select {
		case e.progressChan <- job.Progress:
		default:
		}
		
		// Check if we have more data
		if !hasMore {
			break
		}
	}
	
	// Calculate final statistics
	if job.Progress.RecordsPerSecond > 0 {
		results.AverageRPS = job.Progress.RecordsPerSecond
	}
	
	return results, nil
}

// @decorator: processBatch
// @description: Process a batch of data records
// @param ctx: Context for cancellation
// @param job: Job configuration
// @param batchData: Raw batch data
// @param batchNum: Batch number
// @return *IngestionResults: Batch processing results
// @return error: Any error that occurred during processing
func (e *Engine) processBatch(ctx context.Context, job *IngestionJob, batchData []byte, batchNum int) (*IngestionResults, error) {
	results := &IngestionResults{}
	
	// Parse data based on format
	var records []map[string]interface{}
	var err error
	
	switch job.Config.DataFormat {
	case "json":
		err = json.Unmarshal(batchData, &records)
	case "csv":
		records, err = e.parseCSV(batchData)
	case "xml":
		records, err = e.parseXML(batchData)
	default:
		return nil, errors.Errorf("unsupported data format: %s", job.Config.DataFormat)
	}
	
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse batch data")
	}
	
	results.TotalRecords = int64(len(records))
	
	// Process records in parallel using worker pool
	recordChan := make(chan map[string]interface{}, len(records))
	resultChan := make(chan *RecordResult, len(records))
	
	// Start workers
	for i := 0; i < job.Config.MaxWorkers; i++ {
		go e.processRecordWorker(ctx, job, recordChan, resultChan)
	}
	
	// Send records to workers
	for _, record := range records {
		select {
		case recordChan <- record:
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}
	close(recordChan)
	
	// Collect results
	for i := 0; i < len(records); i++ {
		select {
		case result := <-resultChan:
			if result.Success {
				results.SuccessfulRecords++
			} else {
				results.FailedRecords++
			}
			if result.Duplicate {
				results.DuplicateRecords++
			}
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}
	
	return results, nil
}

// @decorator: RecordResult
// @description: Result of processing a single record
type RecordResult struct {
	Success   bool
	Duplicate bool
	Error     string
}

// @decorator: processRecordWorker
// @description: Worker function for processing individual records
// @param ctx: Context for cancellation
// @param job: Job configuration
// @param recordChan: Channel of records to process
// @param resultChan: Channel for results
func (e *Engine) processRecordWorker(ctx context.Context, job *IngestionJob, recordChan <-chan map[string]interface{}, resultChan chan<- *RecordResult) {
	for {
		select {
		case record, ok := <-recordChan:
			if !ok {
				return
			}
			
			result := &RecordResult{Success: true}
			
			// Validate record
			if job.Config.Validation {
				if validator, exists := e.validators[job.Config.Endpoint]; exists {
					if err := validator.Validate(record); err != nil {
						result.Success = false
						result.Error = err.Error()
						resultChan <- result
						continue
					}
				}
			}
			
			// Transform record
			if job.Config.Transform {
				if processor, exists := e.processors[job.Config.Endpoint]; exists {
					transformedRecord, err := processor.Transform(record)
					if err != nil {
						result.Success = false
						result.Error = err.Error()
						resultChan <- result
						continue
					}
					record = transformedRecord
				}
			}
			
			// Check for duplicates
			if job.Config.Deduplicate {
				isDuplicate, err := e.checkDuplicate(job.Config.TableName, record)
				if err != nil {
					e.logger.WithError(err).Warn("Failed to check for duplicate")
				} else if isDuplicate {
					result.Duplicate = true
					resultChan <- result
					continue
				}
			}
			
			// Store record
			if err := e.storeRecord(job.Config.TableName, record); err != nil {
				result.Success = false
				result.Error = err.Error()
			}
			
			resultChan <- result
			
		case <-ctx.Done():
			return
		}
	}
}

// @decorator: GetJobStatus
// @description: Get the status of an ingestion job
// @param jobID: Job ID
// @return *IngestionJob: Job status
// @return error: Any error that occurred
func (e *Engine) GetJobStatus(jobID string) (*IngestionJob, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	job, exists := e.jobs[jobID]
	if !exists {
		return nil, errors.Errorf("job not found: %s", jobID)
	}
	
	return job, nil
}

// @decorator: ListJobs
// @description: List all ingestion jobs
// @return []*IngestionJob: List of jobs
func (e *Engine) ListJobs() []*IngestionJob {
	e.mu.RLock()
	defer e.mu.RUnlock()
	
	jobs := make([]*IngestionJob, 0, len(e.jobs))
	for _, job := range e.jobs {
		jobs = append(jobs, job)
	}
	
	return jobs
}

// @decorator: CancelJob
// @description: Cancel a running ingestion job
// @param jobID: Job ID to cancel
// @return error: Any error that occurred
func (e *Engine) CancelJob(jobID string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	job, exists := e.jobs[jobID]
	if !exists {
		return errors.Errorf("job not found: %s", jobID)
	}
	
	if job.Status != JobStatusRunning && job.Status != JobStatusPending {
		return errors.Errorf("cannot cancel job in status: %s", job.Status)
	}
	
	job.Status = JobStatusCancelled
	endTime := time.Now()
	job.EndTime = &endTime
	
	e.logger.WithField("job_id", jobID).Info("Job cancelled")
	return nil
}

// Helper functions

func generateJobID() string {
	return fmt.Sprintf("job_%d", time.Now().UnixNano())
}

func (e *Engine) failJob(job *IngestionJob, err error) {
	job.Status = JobStatusFailed
	job.Error = err.Error()
	endTime := time.Now()
	job.EndTime = &endTime
	
	e.logger.WithFields(logrus.Fields{
		"job_id": job.ID,
		"error":  err.Error(),
	}).Error("Job failed")
}

func (e *Engine) addToErrorSummary(results *IngestionResults, errorType, errorMsg string) {
	for i := range results.ErrorSummary {
		if results.ErrorSummary[i].ErrorType == errorType {
			results.ErrorSummary[i].Count++
			return
		}
	}
	
	results.ErrorSummary = append(results.ErrorSummary, ErrorSummary{
		ErrorType:       errorType,
		Count:          1,
		SampleError:    errorMsg,
		FirstOccurrence: time.Now(),
	})
}
