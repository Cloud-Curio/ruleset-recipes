/*
 * Data Ingestion TUI - Congress.gov Bills Ingestion Script
 * 
 * Pre-configured script for bulk ingestion of congressional bills data
 * from the Congress.gov API. Includes data validation, transformation,
 * and comprehensive error handling.
 * 
 * Features:
 * - Bulk ingestion with pagination
 * - Data validation and transformation
 * - Duplicate detection and handling
 * - Progress tracking and monitoring
 * - Comprehensive error handling
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
	"data-ingestion-tui/internal/models"
)

// @decorator: CongressBillsIngestion
// @description: Pre-configured ingestion script for Congress.gov bills
// @version: 1.0.0
// @author: Codegen AI Assistant

// CongressBillsIngestion handles bulk ingestion of congressional bills
type CongressBillsIngestion struct {
	config *CongressBillsConfig
	logger *logrus.Logger
	client *http.Client
}

// @decorator: CongressBillsConfig
// @description: Configuration for Congress bills ingestion
type CongressBillsConfig struct {
	APIKey      string            `json:"api_key"`
	BaseURL     string            `json:"base_url"`
	Congress    int               `json:"congress"`     // e.g., 118 for 118th Congress
	BillType    string            `json:"bill_type"`    // hr, s, hjres, sjres, hconres, sconres, hres, sres
	BatchSize   int               `json:"batch_size"`
	MaxWorkers  int               `json:"max_workers"`
	RateLimit   int               `json:"rate_limit"`
	Headers     map[string]string `json:"headers"`
	TableName   string            `json:"table_name"`
}

// @decorator: CongressBill
// @description: Represents a congressional bill from Congress.gov API
type CongressBill struct {
	Number              string                 `json:"number" db:"number"`
	Title               string                 `json:"title" db:"title"`
	Type                string                 `json:"type" db:"type"`
	Congress            int                    `json:"congress" db:"congress"`
	IntroducedDate      *time.Time             `json:"introducedDate" db:"introduced_date"`
	UpdateDate          *time.Time             `json:"updateDate" db:"update_date"`
	OriginChamber       string                 `json:"originChamber" db:"origin_chamber"`
	LatestAction        *CongressBillAction    `json:"latestAction" db:"-"`
	LatestActionText    string                 `json:"-" db:"latest_action_text"`
	LatestActionDate    *time.Time             `json:"-" db:"latest_action_date"`
	Sponsors            []CongressBillSponsor  `json:"sponsors" db:"-"`
	PrimarySponsorName  string                 `json:"-" db:"primary_sponsor_name"`
	PrimarySponsorState string                 `json:"-" db:"primary_sponsor_state"`
	PrimarySponsorParty string                 `json:"-" db:"primary_sponsor_party"`
	Cosponsors          []CongressBillSponsor  `json:"cosponsors" db:"-"`
	CosponsorCount      int                    `json:"-" db:"cosponsor_count"`
	Committees          []CongressBillCommittee `json:"committees" db:"-"`
	CommitteeNames      string                 `json:"-" db:"committee_names"`
	PolicyArea          string                 `json:"policyArea" db:"policy_area"`
	Subjects            []string               `json:"subjects" db:"-"`
	SubjectsList        string                 `json:"-" db:"subjects_list"`
	ConstitutionalAuthority string             `json:"constitutionalAuthorityStatementText" db:"constitutional_authority"`
	URL                 string                 `json:"url" db:"url"`
	CreatedAt           time.Time              `json:"-" db:"created_at"`
	UpdatedAt           time.Time              `json:"-" db:"updated_at"`
}

// @decorator: CongressBillAction
// @description: Represents an action taken on a congressional bill
type CongressBillAction struct {
	ActionDate *time.Time `json:"actionDate"`
	Text       string     `json:"text"`
	Type       string     `json:"type"`
}

// @decorator: CongressBillSponsor
// @description: Represents a sponsor or cosponsor of a congressional bill
type CongressBillSponsor struct {
	BioguideID   string `json:"bioguideId"`
	FullName     string `json:"fullName"`
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	State        string `json:"state"`
	District     string `json:"district"`
	Party        string `json:"party"`
	IsByRequest  string `json:"isByRequest"`
}

// @decorator: CongressBillCommittee
// @description: Represents a committee associated with a congressional bill
type CongressBillCommittee struct {
	SystemCode string `json:"systemCode"`
	Name       string `json:"name"`
	Chamber    string `json:"chamber"`
	Type       string `json:"type"`
}

// @decorator: NewCongressBillsIngestion
// @description: Create a new Congress bills ingestion script
// @param config: Ingestion configuration
// @param logger: Logger instance
// @return *CongressBillsIngestion: New ingestion script
// @return error: Any error that occurred during creation
func NewCongressBillsIngestion(config *CongressBillsConfig, logger *logrus.Logger) (*CongressBillsIngestion, error) {
	if config == nil {
		return nil, errors.New("configuration cannot be nil")
	}
	
	if config.APIKey == "" {
		return nil, errors.New("API key is required")
	}
	
	if config.BaseURL == "" {
		config.BaseURL = "https://api.congress.gov/v3"
	}
	
	if config.BatchSize == 0 {
		config.BatchSize = 250 // Congress.gov API default
	}
	
	if config.MaxWorkers == 0 {
		config.MaxWorkers = 3
	}
	
	if config.RateLimit == 0 {
		config.RateLimit = 1 // 1 request per second to be respectful
	}
	
	if config.TableName == "" {
		config.TableName = "congress_bills"
	}
	
	// Set default headers
	if config.Headers == nil {
		config.Headers = make(map[string]string)
	}
	config.Headers["X-API-Key"] = config.APIKey
	config.Headers["User-Agent"] = "DataIngestionTUI/1.0"
	config.Headers["Accept"] = "application/json"
	
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	return &CongressBillsIngestion{
		config: config,
		logger: logger,
		client: client,
	}, nil
}

// @decorator: GetJobConfig
// @description: Get the ingestion job configuration
// @return *ingestion.JobConfig: Job configuration
func (c *CongressBillsIngestion) GetJobConfig() *ingestion.JobConfig {
	return &ingestion.JobConfig{
		Endpoint:      "congress_gov",
		APIKey:        c.config.APIKey,
		BaseURL:       c.config.BaseURL,
		Headers:       c.config.Headers,
		BatchSize:     c.config.BatchSize,
		MaxWorkers:    c.config.MaxWorkers,
		RateLimit:     c.config.RateLimit,
		Timeout:       30 * time.Second,
		DataFormat:    "json",
		Validation:    true,
		Transform:     true,
		Deduplicate:   true,
		TableName:     c.config.TableName,
		CreateTable:   true,
		TruncateFirst: false,
		MaxRetries:    3,
		RetryDelay:    5 * time.Second,
		Parameters: map[string]interface{}{
			"congress":   c.config.Congress,
			"bill_type":  c.config.BillType,
		},
	}
}

// @decorator: FetchBatch
// @description: Fetch a batch of bills from Congress.gov API
// @param ctx: Context for cancellation
// @param offset: Offset for pagination
// @param limit: Number of records to fetch
// @return []byte: Raw JSON data
// @return bool: Whether there are more records
// @return error: Any error that occurred during fetch
func (c *CongressBillsIngestion) FetchBatch(ctx context.Context, offset, limit int) ([]byte, bool, error) {
	// Build URL with parameters
	url := fmt.Sprintf("%s/bill", c.config.BaseURL)
	
	// Add query parameters
	params := []string{
		fmt.Sprintf("format=json"),
		fmt.Sprintf("offset=%d", offset),
		fmt.Sprintf("limit=%d", limit),
	}
	
	if c.config.Congress > 0 {
		params = append(params, fmt.Sprintf("congress=%d", c.config.Congress))
	}
	
	if c.config.BillType != "" {
		params = append(params, fmt.Sprintf("type=%s", c.config.BillType))
	}
	
	url += "?" + strings.Join(params, "&")
	
	c.logger.WithFields(logrus.Fields{
		"url":    url,
		"offset": offset,
		"limit":  limit,
	}).Debug("Fetching batch from Congress.gov API")
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to create request")
	}
	
	// Add headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}
	
	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()
	
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, false, errors.Errorf("API returned status %d", resp.StatusCode)
	}
	
	// Parse response to check for more data
	var response struct {
		Bills []json.RawMessage `json:"bills"`
		Pagination struct {
			Count int `json:"count"`
			Next  string `json:"next"`
		} `json:"pagination"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, false, errors.Wrap(err, "failed to decode response")
	}
	
	// Convert bills to JSON array
	billsJSON, err := json.Marshal(response.Bills)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to marshal bills")
	}
	
	// Check if there are more records
	hasMore := response.Pagination.Next != ""
	
	c.logger.WithFields(logrus.Fields{
		"bills_count": len(response.Bills),
		"has_more":    hasMore,
	}).Debug("Fetched batch successfully")
	
	return billsJSON, hasMore, nil
}

// @decorator: ValidateRecord
// @description: Validate a congressional bill record
// @param record: Raw record data
// @return error: Validation error if any
func (c *CongressBillsIngestion) ValidateRecord(record map[string]interface{}) error {
	// Check required fields
	requiredFields := []string{"number", "type", "congress"}
	
	for _, field := range requiredFields {
		if _, exists := record[field]; !exists {
			return errors.Errorf("missing required field: %s", field)
		}
	}
	
	// Validate congress number
	if congress, ok := record["congress"].(float64); ok {
		if congress < 1 || congress > 200 {
			return errors.Errorf("invalid congress number: %v", congress)
		}
	}
	
	// Validate bill type
	if billType, ok := record["type"].(string); ok {
		validTypes := []string{"hr", "s", "hjres", "sjres", "hconres", "sconres", "hres", "sres"}
		valid := false
		for _, validType := range validTypes {
			if strings.EqualFold(billType, validType) {
				valid = true
				break
			}
		}
		if !valid {
			return errors.Errorf("invalid bill type: %s", billType)
		}
	}
	
	return nil
}

// @decorator: TransformRecord
// @description: Transform a raw bill record into normalized format
// @param record: Raw record data
// @return map[string]interface{}: Transformed record
// @return error: Transformation error if any
func (c *CongressBillsIngestion) TransformRecord(record map[string]interface{}) (map[string]interface{}, error) {
	bill := &CongressBill{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	// Basic fields
	if number, ok := record["number"].(string); ok {
		bill.Number = number
	}
	
	if title, ok := record["title"].(string); ok {
		bill.Title = title
	}
	
	if billType, ok := record["type"].(string); ok {
		bill.Type = strings.ToLower(billType)
	}
	
	if congress, ok := record["congress"].(float64); ok {
		bill.Congress = int(congress)
	}
	
	if originChamber, ok := record["originChamber"].(string); ok {
		bill.OriginChamber = originChamber
	}
	
	if url, ok := record["url"].(string); ok {
		bill.URL = url
	}
	
	if policyArea, ok := record["policyArea"].(string); ok {
		bill.PolicyArea = policyArea
	}
	
	if constitutionalAuthority, ok := record["constitutionalAuthorityStatementText"].(string); ok {
		bill.ConstitutionalAuthority = constitutionalAuthority
	}
	
	// Parse dates
	if introducedDateStr, ok := record["introducedDate"].(string); ok {
		if introducedDate, err := time.Parse("2006-01-02", introducedDateStr); err == nil {
			bill.IntroducedDate = &introducedDate
		}
	}
	
	if updateDateStr, ok := record["updateDate"].(string); ok {
		if updateDate, err := time.Parse("2006-01-02T15:04:05Z", updateDateStr); err == nil {
			bill.UpdateDate = &updateDate
		}
	}
	
	// Parse latest action
	if latestActionData, ok := record["latestAction"].(map[string]interface{}); ok {
		if actionText, ok := latestActionData["text"].(string); ok {
			bill.LatestActionText = actionText
		}
		
		if actionDateStr, ok := latestActionData["actionDate"].(string); ok {
			if actionDate, err := time.Parse("2006-01-02", actionDateStr); err == nil {
				bill.LatestActionDate = &actionDate
			}
		}
	}
	
	// Parse sponsors
	if sponsorsData, ok := record["sponsors"].([]interface{}); ok && len(sponsorsData) > 0 {
		if primarySponsor, ok := sponsorsData[0].(map[string]interface{}); ok {
			if fullName, ok := primarySponsor["fullName"].(string); ok {
				bill.PrimarySponsorName = fullName
			}
			if state, ok := primarySponsor["state"].(string); ok {
				bill.PrimarySponsorState = state
			}
			if party, ok := primarySponsor["party"].(string); ok {
				bill.PrimarySponsorParty = party
			}
		}
	}
	
	// Parse cosponsors
	if cosponsorsData, ok := record["cosponsors"].([]interface{}); ok {
		bill.CosponsorCount = len(cosponsorsData)
	}
	
	// Parse committees
	if committeesData, ok := record["committees"].([]interface{}); ok {
		var committeeNames []string
		for _, committeeData := range committeesData {
			if committee, ok := committeeData.(map[string]interface{}); ok {
				if name, ok := committee["name"].(string); ok {
					committeeNames = append(committeeNames, name)
				}
			}
		}
		bill.CommitteeNames = strings.Join(committeeNames, "; ")
	}
	
	// Parse subjects
	if subjectsData, ok := record["subjects"].([]interface{}); ok {
		var subjects []string
		for _, subjectData := range subjectsData {
			if subject, ok := subjectData.(string); ok {
				subjects = append(subjects, subject)
			}
		}
		bill.SubjectsList = strings.Join(subjects, "; ")
	}
	
	// Convert to map for database storage
	result := make(map[string]interface{})
	result["number"] = bill.Number
	result["title"] = bill.Title
	result["type"] = bill.Type
	result["congress"] = bill.Congress
	result["introduced_date"] = bill.IntroducedDate
	result["update_date"] = bill.UpdateDate
	result["origin_chamber"] = bill.OriginChamber
	result["latest_action_text"] = bill.LatestActionText
	result["latest_action_date"] = bill.LatestActionDate
	result["primary_sponsor_name"] = bill.PrimarySponsorName
	result["primary_sponsor_state"] = bill.PrimarySponsorState
	result["primary_sponsor_party"] = bill.PrimarySponsorParty
	result["cosponsor_count"] = bill.CosponsorCount
	result["committee_names"] = bill.CommitteeNames
	result["policy_area"] = bill.PolicyArea
	result["subjects_list"] = bill.SubjectsList
	result["constitutional_authority"] = bill.ConstitutionalAuthority
	result["url"] = bill.URL
	result["created_at"] = bill.CreatedAt
	result["updated_at"] = bill.UpdatedAt
	
	return result, nil
}

// @decorator: GetTableSchema
// @description: Get the database table schema for congressional bills
// @return string: SQL CREATE TABLE statement
func (c *CongressBillsIngestion) GetTableSchema() string {
	return `
	CREATE TABLE IF NOT EXISTS congress_bills (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number TEXT NOT NULL,
		title TEXT,
		type TEXT NOT NULL,
		congress INTEGER NOT NULL,
		introduced_date DATETIME,
		update_date DATETIME,
		origin_chamber TEXT,
		latest_action_text TEXT,
		latest_action_date DATETIME,
		primary_sponsor_name TEXT,
		primary_sponsor_state TEXT,
		primary_sponsor_party TEXT,
		cosponsor_count INTEGER DEFAULT 0,
		committee_names TEXT,
		policy_area TEXT,
		subjects_list TEXT,
		constitutional_authority TEXT,
		url TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		UNIQUE(number, type, congress)
	);
	
	CREATE INDEX IF NOT EXISTS idx_congress_bills_congress ON congress_bills(congress);
	CREATE INDEX IF NOT EXISTS idx_congress_bills_type ON congress_bills(type);
	CREATE INDEX IF NOT EXISTS idx_congress_bills_introduced_date ON congress_bills(introduced_date);
	CREATE INDEX IF NOT EXISTS idx_congress_bills_sponsor ON congress_bills(primary_sponsor_name);
	CREATE INDEX IF NOT EXISTS idx_congress_bills_policy_area ON congress_bills(policy_area);
	`
}

// @decorator: GetTotalRecordCount
// @description: Get the total number of records available
// @param ctx: Context for cancellation
// @return int64: Total record count
// @return error: Any error that occurred
func (c *CongressBillsIngestion) GetTotalRecordCount(ctx context.Context) (int64, error) {
	// Build URL for count query
	url := fmt.Sprintf("%s/bill", c.config.BaseURL)
	
	params := []string{
		"format=json",
		"limit=1", // We only need the count, not the data
	}
	
	if c.config.Congress > 0 {
		params = append(params, fmt.Sprintf("congress=%d", c.config.Congress))
	}
	
	if c.config.BillType != "" {
		params = append(params, fmt.Sprintf("type=%s", c.config.BillType))
	}
	
	url += "?" + strings.Join(params, "&")
	
	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create request")
	}
	
	// Add headers
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}
	
	// Make request
	resp, err := c.client.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "failed to make request")
	}
	defer resp.Body.Close()
	
	// Parse response
	var response struct {
		Pagination struct {
			Count int `json:"count"`
		} `json:"pagination"`
	}
	
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return 0, errors.Wrap(err, "failed to decode response")
	}
	
	return int64(response.Pagination.Count), nil
}

// @decorator: GetDuplicateCheckFields
// @description: Get fields to use for duplicate checking
// @return []string: Field names for duplicate checking
func (c *CongressBillsIngestion) GetDuplicateCheckFields() []string {
	return []string{"number", "type", "congress"}
}
