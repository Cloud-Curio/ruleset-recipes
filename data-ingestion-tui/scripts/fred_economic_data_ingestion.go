/*
 * Data Ingestion TUI - FRED Economic Data Ingestion Script
 * 
 * Pre-configured script for bulk ingestion of economic data from the
 * Federal Reserve Economic Data (FRED) API. Includes data validation,
 * transformation, and comprehensive error handling.
 * 
 * Features:
 * - Multiple economic series ingestion
 * - Time series data handling
 * - Data validation and transformation
 * - Duplicate detection and handling
 * - Progress tracking and monitoring
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package scripts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"data-ingestion-tui/internal/ingestion"
)

// @decorator: FREDEconomicDataIngestion
// @description: Pre-configured ingestion script for FRED economic data
// @version: 1.0.0
// @author: Codegen AI Assistant

// FREDEconomicDataIngestion handles bulk ingestion of FRED economic data
type FREDEconomicDataIngestion struct {
	config *FREDEconomicDataConfig
	logger *logrus.Logger
	client *http.Client
}

// @decorator: FREDEconomicDataConfig
// @description: Configuration for FRED economic data ingestion
type FREDEconomicDataConfig struct {
	APIKey      string            `json:"api_key"`
	BaseURL     string            `json:"base_url"`
	SeriesID    string            `json:"series_id"`    // e.g., "GDP", "UNRATE", "CPIAUCSL"
	StartDate   string            `json:"start_date"`   // YYYY-MM-DD format
	EndDate     string            `json:"end_date"`     // YYYY-MM-DD format
	Frequency   string            `json:"frequency"`    // d, w, bw, m, q, sa, a
	BatchSize   int               `json:"batch_size"`
	MaxWorkers  int               `json:"max_workers"`
	RateLimit   int               `json:"rate_limit"`
	Headers     map[string]string `json:"headers"`
	TableName   string            `json:"table_name"`
}

// @decorator: FREDObservation
// @description: Represents an economic observation from FRED API
type FREDObservation struct {
	SeriesID      string     `json:"series_id" db:"series_id"`
	Date          time.Time  `json:"date" db:"date"`
	Value         *float64   `json:"value" db:"value"`
	ValueStr      string     `json:"value_str" db:"value_str"`
	RealtimeStart time.Time  `json:"realtime_start" db:"realtime_start"`
	RealtimeEnd   time.Time  `json:"realtime_end" db:"realtime_end"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	UpdatedAt     time.Time  `json:"-" db:"updated_at"`
}

// @decorator: NewFREDEconomicDataIngestion
// @description: Create a new FRED economic data ingestion script
// @param config: Ingestion configuration
// @param logger: Logger instance
// @return *FREDEconomicDataIngestion: New ingestion script
// @return error: Any error that occurred during creation
func NewFREDEconomicDataIngestion(config *FREDEconomicDataConfig, logger *logrus.Logger) (*FREDEconomicDataIngestion, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}
	
	if config.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	
	if config.SeriesID == "" {
		return nil, errors.New("series ID is required")
	}
	
	if config.BaseURL == "" {
		config.BaseURL = "https://api.stlouisfed.org/fred"
	}
	
	if config.BatchSize == 0 {
		config.BatchSize = 1000 // FRED API can handle large batches
	}
	
	if config.MaxWorkers == 0 {
		config.MaxWorkers = 2
	}
	
	if config.RateLimit == 0 {
		config.RateLimit = 2 // 2 requests per second (120 per minute limit)
	}
	
	if config.TableName == "" {
		config.TableName = "fred_observations"
	}
	
	// Set default headers
	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}
	config.Headers["User-Agent"] = "DataIngestionTUI/1.0"
	config.Headers["Accept"] = "application/json"
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	return &FREDEconomicDataIngestion{
		config: config,
		logger: logger,
		client: client,
	}, nil
}

// @decorator: GetJobConfig
// @description: Get the ingestion job configuration
// @return *ingestion.JobConfig: Job configuration
func (f *FREDEconomicDataIngestion) GetJobConfig() *ingestion.JobConfig {
	return &ingestion.JobConfig{
		Endpoint:      "fred_stlouisfed",
		APIKey:        f.config.APIKey,
		BaseURL:       f.config.BaseURL,
		Headers:       f.config.Headers,
		BatchSize:     f.config.BatchSize,
		MaxWorkers:    f.config.MaxWorkers,
		RateLimit:     f.config.RateLimit,
		Timeout:       30 * time.Second,
		DataFormat:    "json",
		Validation:    true,
		Transform:     true,
		Deduplicate:   true,
		TableName:     f.config.TableName,
		CreateTable:   true,
		TruncateFirst: false,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		Parameters: map[string]interface{}{
			"series_id":   f.config.SeriesID,
			"start_date":  f.config.StartDate,
			"end_date":    f.config.EndDate,
			"frequency":   f.config.Frequency,
		},
	}
}

// @decorator: FetchBatch
// @description: Fetch a batch of observations from FRED API
// @param ctx: Context for cancellation
// @param offset: Offset for pagination
// @param limit: Number of records to fetch
// @return []byte: Raw JSON data
// @return bool: Whether there are more records
// @return error: Any error that occurred during fetch
func (f *FREDEconomicDataIngestion) FetchBatch(ctx context.Context, offset, limit int) ([]byte, bool, error) {
	// Build URL with parameters
	url := fmt.Sprintf("%s/series/observations", f.config.BaseURL)
	
	// Add query parameters
	params := []string{
		fmt.Sprintf("series_id=%s", f.config.SeriesID),
		fmt.Sprintf("api_key=%s", f.config.APIKey),
		"file_type=json",
		fmt.Sprintf("offset=%d", offset),
		fmt.Sprintf("limit=%d", limit),
	}
	
	if f.config.StartDate != "" {
		params = append(params, fmt.Sprintf("observation_start=%s", f.config.StartDate))
	}
	
	if f.config.EndDate != "" {
		params = append(params, fmt.Sprintf("observation_end=%s", f.config.EndDate))
	}
	
	if f.config.Frequency != "" {
		params = append(params, fmt.Sprintf("frequency=%s", f.config.Frequency))
	}
	
	url += "?" + strings.Join(params, "&")
	
	f.logger.WithFields(logrus.Fields{
		"url":       url,
		"series_id": f.config.SeriesID,
		"offset":    offset,
		"limit":     limit,
	}).Debug("Fetching batch from FRED API")
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to create request")
	}
	
	// Add headers
	for key, value := range f.config.Headers {
		req.Header.Set(key, value)
	}
	
	// Make request
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, false, errors.Errorf("API returned status %d", resp.StatusCode)
	}
	
	// Parse response
	var response struct {
		Observations []json.RawMessage `json:"observations"`
		Count        int               `json:"count"`
		Offset       int               `json:"offset"`
		Limit        int               `json:"limit"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, false, errors.Wrap(err, "failed to decode response")
	}
	
	// Convert observations to JSON array
	observationsJSON, err := json.Marshal(response.Observations)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to marshal observations")
	}
	
	// Check if there are more records
	hasMore := len(response.Observations) == limit && (response.Offset+limit) < response.Count
	
	f.logger.WithFields(logrus.Fields{
		"observations_count": len(response.Observations),
		"total_count":        response.Count,
		"has_more":          hasMore,
	}).Debug("Fetched batch successfully")
	
	return observationsJSON, hasMore, nil
}

// @decorator: ValidateRecord
// @description: Validate a FRED observation record
// @param record: Raw record data
// @return error: Validation error if any
func (f *FREDEconomicDataIngestion) ValidateRecord(record map[string]interface{}) error {
	// Check required fields
	requiredFields := []string{"date", "realtime_start", "realtime_end"}
	
	for _, field := range requiredFields {
		if _, exists := record[field]; !exists {
			return errors.Errorf("missing required field: %s", field)
		}
	}
	
	// Validate date format
	if dateStr, ok := record["date"].(string); ok {
		if _, err := time.Parse("2006-01-02", dateStr); err != nil {
			return errors.Errorf("invalid date format: %s", dateStr)
		}
	}
	
	// Validate realtime_start format
	if realtimeStartStr, ok := record["realtime_start"].(string); ok {
		if _, err := time.Parse("2006-01-02", realtimeStartStr); err != nil {
			return errors.Errorf("invalid realtime_start format: %s", realtimeStartStr)
		}
	}
	
	// Validate realtime_end format
	if realtimeEndStr, ok := record["realtime_end"].(string); ok {
		if _, err := time.Parse("2006-01-02", realtimeEndStr); err != nil {
			return errors.Errorf("invalid realtime_end format: %s", realtimeEndStr)
		}
	}
	
	// Validate value (can be "." for missing data)
	if valueStr, ok := record["value"].(string); ok {
		if valueStr != "." {
			if _, err := strconv.ParseFloat(valueStr, 64); err != nil {
				return errors.Errorf("invalid value format: %s", valueStr)
			}
		}
	}
	
	return nil
}

// @decorator: TransformRecord
// @description: Transform a raw observation record into normalized format
// @param record: Raw record data
// @return map[string]interface{}: Transformed record
// @return error: Transformation error if any
func (f *FREDEconomicDataIngestion) TransformRecord(record map[string]interface{}) (map[string]interface{}, error) {
	observation := &FREDObservation{
		SeriesID:  f.config.SeriesID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Parse date
	if dateStr, ok := record["date"].(string); ok {
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			observation.Date = date
		} else {
			return nil, errors.Wrap(err, "failed to parse date")
		}
	}
	
	// Parse realtime_start
	if realtimeStartStr, ok := record["realtime_start"].(string); ok {
		if realtimeStart, err := time.Parse("2006-01-02", realtimeStartStr); err == nil {
			observation.RealtimeStart = realtimeStart
		} else {
			return nil, errors.Wrap(err, "failed to parse realtime_start")
		}
	}
	
	// Parse realtime_end
	if realtimeEndStr, ok := record["realtime_end"].(string); ok {
		if realtimeEnd, err := time.Parse("2006-01-02", realtimeEndStr); err == nil {
			observation.RealtimeEnd = realtimeEnd
		} else {
			return nil, errors.Wrap(err, "failed to parse realtime_end")
		}
	}
	
	// Parse value
	if valueStr, ok := record["value"].(string); ok {
		observation.ValueStr = valueStr
		if valueStr != "." {
			if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
				observation.Value = &value
			}
		}
	}
	
	// Convert to map for database storage
	result := make(map[string]interface{})
	result["series_id"] = observation.SeriesID
	result["date"] = observation.Date
	result["value"] = observation.Value
	result["value_str"] = observation.ValueStr
	result["realtime_start"] = observation.RealtimeStart
	result["realtime_end"] = observation.RealtimeEnd
	result["created_at"] = observation.CreatedAt
	result["updated_at"] = observation.UpdatedAt
	
	return result, nil
}

// @decorator: GetTableSchema
// @description: Get the database table schema for FRED observations
// @return string: SQL CREATE TABLE statement
func (f *FREDEconomicDataIngestion) GetTableSchema() string {
	return `
	CREATE TABLE IF NOT EXISTS fred_observations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		series_id TEXT NOT NULL,
		date DATE NOT NULL,
		value REAL,
		value_str TEXT NOT NULL,
		realtime_start DATE NOT NULL,
		realtime_end DATE NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(series_id, date, realtime_start, realtime_end)
	);
	
	CREATE INDEX IF NOT EXISTS idx_fred_observations_series_id ON fred_observations(series_id);
	CREATE INDEX IF NOT EXISTS idx_fred_observations_date ON fred_observations(date);
	CREATE INDEX IF NOT EXISTS idx_fred_observations_series_date ON fred_observations(series_id, date);
	CREATE INDEX IF NOT EXISTS idx_fred_observations_realtime ON fred_observations(realtime_start, realtime_end);
	`
}

// @decorator: GetTotalRecordCount
// @description: Get the total number of records available
// @param ctx: Context for cancellation
// @return int64: Total record count
// @return error: Any error that occurred
func (f *FREDEconomicDataIngestion) GetTotalRecordCount(ctx context.Context) (int64, error) {
	// Build URL for count query
	url := fmt.Sprintf("%s/series/observations", f.config.BaseURL)
	
	params := []string{
		fmt.Sprintf("series_id=%s", f.config.SeriesID),
		fmt.Sprintf("api_key=%s", f.config.APIKey),
		"file_type=json",
		"limit=1", // We only need the count, not the data
	}
	
	if f.config.StartDate != "" {
		params = append(params, fmt.Sprintf("observation_start=%s", f.config.StartDate))
	}
	
	if f.config.EndDate != "" {
		params = append(params, fmt.Sprintf("observation_end=%s", f.config.EndDate))
	}
	
	if f.config.Frequency != "" {
		params = append(params, fmt.Sprintf("frequency=%s", f.config.Frequency))
	}
	
	url += "?" + strings.Join(params, "&")
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create request")
	}
	
	// Add headers
	for key, value := range f.config.Headers {
		req.Header.Set(key, value)
	}
	
	// Make request
	resp, err := f.client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()
	
	// Parse response
	var response struct {
		Count int `json:"count"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return 0, errors.Wrap(err, "failed to decode response")
	}
	
	return int64(response.Count), nil
}

// @decorator: GetDuplicateCheckFields
// @description: Get fields to use for duplicate checking
// @return []string: Field names for duplicate checking
func (f *FREDEconomicDataIngestion) GetDuplicateCheckFields() []string {
	return []string{"series_id", "date", "realtime_start", "realtime_end"}
}
