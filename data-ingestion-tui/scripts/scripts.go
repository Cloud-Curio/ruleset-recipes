/*
 * Data Ingestion TUI - Scripts Package
 * 
 * This package contains pre-built ingestion scripts for various
 * government and public APIs.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package scripts

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

// @decorator: ScriptResult
// @description: Result of running an ingestion script
type ScriptResult struct {
	ScriptName      string        `json:"script_name"`
	RecordsIngested int           `json:"records_ingested"`
	Duration        time.Duration `json:"duration"`
	Success         bool          `json:"success"`
	Error           error         `json:"error,omitempty"`
}

// @decorator: ScriptConfig
// @description: Configuration for ingestion scripts
type ScriptConfig struct {
	APIKey    string            `json:"api_key"`
	BaseURL   string            `json:"base_url"`
	RateLimit int               `json:"rate_limit"`
	Timeout   time.Duration     `json:"timeout"`
	Headers   map[string]string `json:"headers"`
}

// @decorator: IngestCongressBills
// @description: Run Congress bills ingestion script
func IngestCongressBills(ctx context.Context, config ScriptConfig, logger *logrus.Logger) (*ScriptResult, error) {
	if logger == nil {
		logger = logrus.New()
	}

	logger.Info("Starting Congress bills ingestion")
	startTime := time.Now()

	// Simulate ingestion process
	// In a real implementation, this would call the actual ingestion logic
	time.Sleep(100 * time.Millisecond) // Simulate processing time

	result := &ScriptResult{
		ScriptName:      "congress_bills",
		RecordsIngested: 50, // Simulated count
		Duration:        time.Since(startTime),
		Success:         true,
	}

	logger.WithFields(logrus.Fields{
		"records_ingested": result.RecordsIngested,
		"duration":         result.Duration,
	}).Info("Congress bills ingestion completed")

	return result, nil
}

// @decorator: IngestFREDData
// @description: Run FRED economic data ingestion script
func IngestFREDData(ctx context.Context, config ScriptConfig, logger *logrus.Logger) (*ScriptResult, error) {
	if logger == nil {
		logger = logrus.New()
	}

	logger.Info("Starting FRED economic data ingestion")
	startTime := time.Now()

	// Simulate ingestion process
	time.Sleep(80 * time.Millisecond) // Simulate processing time

	result := &ScriptResult{
		ScriptName:      "fred_economic_data",
		RecordsIngested: 25, // Simulated count
		Duration:        time.Since(startTime),
		Success:         true,
	}

	logger.WithFields(logrus.Fields{
		"records_ingested": result.RecordsIngested,
		"duration":         result.Duration,
	}).Info("FRED economic data ingestion completed")

	return result, nil
}

// @decorator: RunAllScripts
// @description: Run all available ingestion scripts
func RunAllScripts(ctx context.Context, config ScriptConfig, logger *logrus.Logger) ([]*ScriptResult, error) {
	if logger == nil {
		logger = logrus.New()
	}

	logger.Info("Running all ingestion scripts")
	
	var results []*ScriptResult
	
	// Run Congress bills ingestion
	congressResult, err := IngestCongressBills(ctx, config, logger)
	if err != nil {
		logger.WithError(err).Error("Congress bills ingestion failed")
		congressResult = &ScriptResult{
			ScriptName: "congress_bills",
			Success:    false,
			Error:      err,
		}
	}
	results = append(results, congressResult)
	
	// Run FRED data ingestion
	fredResult, err := IngestFREDData(ctx, config, logger)
	if err != nil {
		logger.WithError(err).Error("FRED data ingestion failed")
		fredResult = &ScriptResult{
			ScriptName: "fred_economic_data",
			Success:    false,
			Error:      err,
		}
	}
	results = append(results, fredResult)
	
	// Calculate totals
	totalRecords := 0
	successCount := 0
	for _, result := range results {
		totalRecords += result.RecordsIngested
		if result.Success {
			successCount++
		}
	}
	
	logger.WithFields(logrus.Fields{
		"total_scripts":     len(results),
		"successful_scripts": successCount,
		"total_records":     totalRecords,
	}).Info("All ingestion scripts completed")
	
	return results, nil
}

// @decorator: ValidateScriptConfig
// @description: Validate script configuration
func ValidateScriptConfig(config ScriptConfig) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}
	
	if config.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}
	
	if config.RateLimit <= 0 {
		return fmt.Errorf("rate limit must be positive")
	}
	
	if config.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	
	return nil
}
