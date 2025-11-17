/*
 * Data Ingestion TUI - Scripts Tests
 * 
 * Tests for the ingestion scripts functionality.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package scripts

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sirupsen/logrus"
)

// @decorator: TestIngestCongressBills
// @description: Test Congress bills ingestion script
func TestIngestCongressBills(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	config := ScriptConfig{
		APIKey:    "test_key",
		BaseURL:   "https://api.congress.gov",
		RateLimit: 100,
		Timeout:   30 * time.Second,
		Headers:   make(map[string]string),
	}
	
	ctx := context.Background()
	result, err := IngestCongressBills(ctx, config, logger)
	
	require.NoError(t, err)
	require.NotNil(t, result)
	
	assert.Equal(t, "congress_bills", result.ScriptName)
	assert.True(t, result.Success)
	assert.Greater(t, result.RecordsIngested, 0)
	assert.Greater(t, result.Duration, time.Duration(0))
	
	t.Logf("✅ Congress bills ingestion test passed: %d records in %v", 
		result.RecordsIngested, result.Duration)
}

// @decorator: TestIngestFREDData
// @description: Test FRED economic data ingestion script
func TestIngestFREDData(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	config := ScriptConfig{
		APIKey:    "test_key",
		BaseURL:   "https://api.stlouisfed.org",
		RateLimit: 100,
		Timeout:   30 * time.Second,
		Headers:   make(map[string]string),
	}
	
	ctx := context.Background()
	result, err := IngestFREDData(ctx, config, logger)
	
	require.NoError(t, err)
	require.NotNil(t, result)
	
	assert.Equal(t, "fred_economic_data", result.ScriptName)
	assert.True(t, result.Success)
	assert.Greater(t, result.RecordsIngested, 0)
	assert.Greater(t, result.Duration, time.Duration(0))
	
	t.Logf("✅ FRED data ingestion test passed: %d records in %v", 
		result.RecordsIngested, result.Duration)
}

// @decorator: TestRunAllScripts
// @description: Test running all ingestion scripts
func TestRunAllScripts(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	
	config := ScriptConfig{
		APIKey:    "test_key",
		BaseURL:   "https://api.example.com",
		RateLimit: 100,
		Timeout:   30 * time.Second,
		Headers:   make(map[string]string),
	}
	
	ctx := context.Background()
	results, err := RunAllScripts(ctx, config, logger)
	
	require.NoError(t, err)
	require.NotNil(t, results)
	assert.Len(t, results, 2) // Should have Congress and FRED results
	
	totalRecords := 0
	successCount := 0
	
	for _, result := range results {
		assert.NotEmpty(t, result.ScriptName)
		assert.Greater(t, result.Duration, time.Duration(0))
		
		if result.Success {
			successCount++
			totalRecords += result.RecordsIngested
		}
	}
	
	assert.Equal(t, 2, successCount, "All scripts should succeed")
	assert.Greater(t, totalRecords, 0, "Should ingest some records")
	
	t.Logf("✅ All scripts test passed: %d successful scripts, %d total records", 
		successCount, totalRecords)
}

// @decorator: TestValidateScriptConfig
// @description: Test script configuration validation
func TestValidateScriptConfig(t *testing.T) {
	// Test valid config
	validConfig := ScriptConfig{
		APIKey:    "test_key",
		BaseURL:   "https://api.example.com",
		RateLimit: 100,
		Timeout:   30 * time.Second,
		Headers:   make(map[string]string),
	}
	
	err := ValidateScriptConfig(validConfig)
	assert.NoError(t, err, "Valid config should pass validation")
	
	// Test invalid configs
	testCases := []struct {
		name   string
		config ScriptConfig
		hasError bool
	}{
		{
			name: "missing API key",
			config: ScriptConfig{
				BaseURL:   "https://api.example.com",
				RateLimit: 100,
				Timeout:   30 * time.Second,
			},
			hasError: true,
		},
		{
			name: "missing base URL",
			config: ScriptConfig{
				APIKey:    "test_key",
				RateLimit: 100,
				Timeout:   30 * time.Second,
			},
			hasError: true,
		},
		{
			name: "invalid rate limit",
			config: ScriptConfig{
				APIKey:    "test_key",
				BaseURL:   "https://api.example.com",
				RateLimit: 0,
				Timeout:   30 * time.Second,
			},
			hasError: true,
		},
		{
			name: "invalid timeout",
			config: ScriptConfig{
				APIKey:    "test_key",
				BaseURL:   "https://api.example.com",
				RateLimit: 100,
				Timeout:   0,
			},
			hasError: true,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateScriptConfig(tc.config)
			if tc.hasError {
				assert.Error(t, err, "Config should be invalid: %s", tc.name)
			} else {
				assert.NoError(t, err, "Config should be valid: %s", tc.name)
			}
		})
	}
	
	t.Log("✅ Script configuration validation test passed")
}
