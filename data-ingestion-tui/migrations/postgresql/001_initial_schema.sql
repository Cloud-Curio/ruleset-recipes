-- Data Ingestion TUI - PostgreSQL Initial Schema Migration
-- 
-- This migration creates the initial database schema for the Data Ingestion TUI
-- application with comprehensive tables for all supported data sources.
-- 
-- Based on actual API schemas from:
-- - Congress.gov API v3 (https://api.congress.gov/v3/)
-- - FRED API (https://api.stlouisfed.org/fred/)
-- - SEC EDGAR API (https://www.sec.gov/edgar/sec-api-documentation)
-- - FBI Crime Data API (https://api.usa.gov/crime/fbi/sapi/)
-- - Census Bureau API (https://www.census.gov/data/developers/data-sets.html)
-- 
-- Author: Codegen AI Assistant
-- Created: 2024
-- License: MIT

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create schemas for organization
CREATE SCHEMA IF NOT EXISTS ingestion;
CREATE SCHEMA IF NOT EXISTS government;
CREATE SCHEMA IF NOT EXISTS economic;
CREATE SCHEMA IF NOT EXISTS financial;
CREATE SCHEMA IF NOT EXISTS security;
CREATE SCHEMA IF NOT EXISTS telemetry;

-- Set search path
SET search_path TO ingestion, government, economic, financial, security, telemetry, public;

-- ============================================================================
-- INGESTION MANAGEMENT TABLES
-- ============================================================================

-- Ingestion jobs tracking
CREATE TABLE ingestion.jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    endpoint VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    job_type VARCHAR(50) NOT NULL,
    config JSONB NOT NULL,
    progress JSONB,
    results JSONB,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT jobs_status_check CHECK (status IN ('pending', 'running', 'completed', 'failed', 'cancelled', 'paused'))
);

-- API endpoints configuration
CREATE TABLE ingestion.api_endpoints (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    base_url VARCHAR(500) NOT NULL,
    api_version VARCHAR(20),
    description TEXT,
    rate_limit INTEGER DEFAULT 60,
    rate_limit_window INTEGER DEFAULT 60, -- seconds
    timeout_seconds INTEGER DEFAULT 30,
    retry_count INTEGER DEFAULT 3,
    retry_delay_seconds INTEGER DEFAULT 5,
    headers JSONB,
    auth_type VARCHAR(50),
    auth_config JSONB,
    schema_config JSONB,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Data quality metrics
CREATE TABLE ingestion.data_quality_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID REFERENCES ingestion.jobs(id) ON DELETE CASCADE,
    endpoint VARCHAR(100) NOT NULL,
    table_name VARCHAR(100) NOT NULL,
    metric_type VARCHAR(50) NOT NULL,
    metric_value DECIMAL(10,4) NOT NULL,
    total_records BIGINT,
    valid_records BIGINT,
    invalid_records BIGINT,
    duplicate_records BIGINT,
    null_records BIGINT,
    details JSONB,
    measured_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT quality_metric_type_check CHECK (metric_type IN ('completeness', 'accuracy', 'consistency', 'validity', 'uniqueness'))
);

-- ============================================================================
-- CONGRESS.GOV API TABLES (Based on actual API v3 schema)
-- ============================================================================

-- Congressional bills (from Congress.gov API)
CREATE TABLE government.congress_bills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    congress INTEGER NOT NULL,
    bill_type VARCHAR(10) NOT NULL,
    bill_number VARCHAR(20) NOT NULL,
    title TEXT,
    short_title TEXT,
    official_title TEXT,
    introduced_date DATE,
    update_date TIMESTAMP WITH TIME ZONE,
    latest_action_date DATE,
    latest_action_text TEXT,
    latest_action_type VARCHAR(100),
    origin_chamber VARCHAR(20),
    origin_chamber_code VARCHAR(5),
    
    -- Sponsor information
    sponsor_bioguide_id VARCHAR(20),
    sponsor_full_name VARCHAR(200),
    sponsor_first_name VARCHAR(100),
    sponsor_last_name VARCHAR(100),
    sponsor_middle_name VARCHAR(100),
    sponsor_suffix VARCHAR(20),
    sponsor_nickname VARCHAR(100),
    sponsor_state VARCHAR(2),
    sponsor_district VARCHAR(10),
    sponsor_party VARCHAR(20),
    sponsor_url VARCHAR(500),
    
    -- Bill status and tracking
    is_private BOOLEAN DEFAULT false,
    constitutional_authority_statement TEXT,
    policy_area VARCHAR(200),
    
    -- URLs and references
    congress_gov_url VARCHAR(500),
    govtrack_url VARCHAR(500),
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT congress_bills_unique UNIQUE (congress, bill_type, bill_number),
    CONSTRAINT congress_bills_congress_check CHECK (congress >= 1 AND congress <= 200),
    CONSTRAINT congress_bills_type_check CHECK (bill_type IN ('hr', 's', 'hjres', 'sjres', 'hconres', 'sconres', 'hres', 'sres'))
);

-- Bill cosponsors
CREATE TABLE government.congress_bill_cosponsors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bill_id UUID REFERENCES government.congress_bills(id) ON DELETE CASCADE,
    bioguide_id VARCHAR(20) NOT NULL,
    full_name VARCHAR(200),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    middle_name VARCHAR(100),
    suffix VARCHAR(20),
    state VARCHAR(2),
    district VARCHAR(10),
    party VARCHAR(20),
    sponsorship_date DATE,
    sponsorship_withdrawn_date DATE,
    is_original_cosponsor BOOLEAN DEFAULT false,
    url VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bill committees
CREATE TABLE government.congress_bill_committees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bill_id UUID REFERENCES government.congress_bills(id) ON DELETE CASCADE,
    committee_code VARCHAR(20) NOT NULL,
    committee_name VARCHAR(300),
    chamber VARCHAR(20),
    committee_type VARCHAR(50),
    subcommittee_code VARCHAR(20),
    subcommittee_name VARCHAR(300),
    activity_name VARCHAR(100),
    activity_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bill subjects/topics
CREATE TABLE government.congress_bill_subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bill_id UUID REFERENCES government.congress_bills(id) ON DELETE CASCADE,
    subject VARCHAR(500) NOT NULL,
    subject_source VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Bill actions
CREATE TABLE government.congress_bill_actions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bill_id UUID REFERENCES government.congress_bills(id) ON DELETE CASCADE,
    action_date DATE NOT NULL,
    action_time TIME,
    action_text TEXT NOT NULL,
    action_type VARCHAR(100),
    action_code VARCHAR(20),
    chamber VARCHAR(20),
    committee_code VARCHAR(20),
    committee_name VARCHAR(300),
    source_system VARCHAR(50),
    source_system_code VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ============================================================================
-- FRED (Federal Reserve Economic Data) TABLES
-- ============================================================================

-- FRED series metadata
CREATE TABLE economic.fred_series (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id VARCHAR(50) NOT NULL UNIQUE,
    title TEXT NOT NULL,
    units VARCHAR(100),
    units_short VARCHAR(50),
    frequency VARCHAR(20),
    frequency_short VARCHAR(10),
    seasonal_adjustment VARCHAR(50),
    seasonal_adjustment_short VARCHAR(20),
    last_updated TIMESTAMP WITH TIME ZONE,
    popularity INTEGER,
    group_popularity INTEGER,
    notes TEXT,
    
    -- Date ranges
    observation_start DATE,
    observation_end DATE,
    realtime_start DATE,
    realtime_end DATE,
    
    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT fred_frequency_check CHECK (frequency IN ('Daily', 'Weekly', 'Biweekly', 'Monthly', 'Quarterly', 'Semiannual', 'Annual'))
);

-- FRED observations (actual data points)
CREATE TABLE economic.fred_observations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id VARCHAR(50) NOT NULL REFERENCES economic.fred_series(series_id) ON DELETE CASCADE,
    observation_date DATE NOT NULL,
    value DECIMAL(20,6),
    value_string VARCHAR(50), -- For non-numeric values like "."
    realtime_start DATE NOT NULL,
    realtime_end DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT fred_observations_unique UNIQUE (series_id, observation_date, realtime_start, realtime_end)
);

-- FRED categories
CREATE TABLE economic.fred_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id INTEGER NOT NULL UNIQUE,
    name VARCHAR(300) NOT NULL,
    parent_id INTEGER,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id)
);

-- FRED series categories relationship
CREATE TABLE economic.fred_series_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    series_id VARCHAR(50) REFERENCES economic.fred_series(series_id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES economic.fred_categories(category_id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT fred_series_categories_unique UNIQUE (series_id, category_id)
);

-- ============================================================================
-- SEC EDGAR TABLES (Based on SEC API schema)
-- ============================================================================

-- Company facts from SEC EDGAR
CREATE TABLE financial.sec_company_facts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cik VARCHAR(20) NOT NULL,
    entity_name VARCHAR(500),
    sic VARCHAR(10),
    sic_description VARCHAR(300),
    insider_transaction_for_owner_exists BOOLEAN,
    insider_transaction_for_issuer_exists BOOLEAN,
    tickers JSONB, -- Array of ticker symbols
    exchanges JSONB, -- Array of exchanges
    ein VARCHAR(20),
    description TEXT,
    website VARCHAR(500),
    investor_website VARCHAR(500),
    category VARCHAR(100),
    fiscal_year_end VARCHAR(10),
    state_of_incorporation VARCHAR(50),
    state_of_incorporation_description VARCHAR(100),
    addresses JSONB,
    phone VARCHAR(50),
    flags TEXT,
    former_names JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT sec_company_facts_cik_unique UNIQUE (cik)
);

-- SEC financial facts/metrics
CREATE TABLE financial.sec_financial_facts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES financial.sec_company_facts(id) ON DELETE CASCADE,
    cik VARCHAR(20) NOT NULL,
    taxonomy VARCHAR(50) NOT NULL, -- us-gaap, dei, invest, etc.
    tag VARCHAR(200) NOT NULL,
    label TEXT,
    description TEXT,
    unit VARCHAR(50),
    end_date DATE,
    value DECIMAL(20,2),
    accession_number VARCHAR(50),
    fiscal_year INTEGER,
    fiscal_period VARCHAR(10),
    form VARCHAR(20),
    filed_date DATE,
    frame VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT sec_financial_facts_unique UNIQUE (cik, taxonomy, tag, end_date, accession_number)
);

-- SEC filings
CREATE TABLE financial.sec_filings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES financial.sec_company_facts(id) ON DELETE CASCADE,
    cik VARCHAR(20) NOT NULL,
    accession_number VARCHAR(50) NOT NULL,
    filing_date DATE,
    report_date DATE,
    acceptance_date_time TIMESTAMP WITH TIME ZONE,
    act VARCHAR(10),
    form VARCHAR(20),
    file_number VARCHAR(50),
    film_number VARCHAR(20),
    items VARCHAR(500),
    size BIGINT,
    is_xbrl BOOLEAN DEFAULT false,
    is_inline_xbrl BOOLEAN DEFAULT false,
    primary_document VARCHAR(200),
    primary_doc_description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT sec_filings_unique UNIQUE (cik, accession_number)
);

-- ============================================================================
-- FBI CRIME DATA TABLES
-- ============================================================================

-- FBI crime statistics by state
CREATE TABLE security.fbi_crime_state (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    state_id INTEGER,
    state_abbr VARCHAR(2) NOT NULL,
    state_name VARCHAR(100) NOT NULL,
    year INTEGER NOT NULL,
    population BIGINT,
    
    -- Violent crimes
    violent_crime BIGINT,
    homicide BIGINT,
    rape BIGINT,
    robbery BIGINT,
    aggravated_assault BIGINT,
    
    -- Property crimes
    property_crime BIGINT,
    burglary BIGINT,
    larceny BIGINT,
    motor_vehicle_theft BIGINT,
    arson BIGINT,
    
    -- Rates per 100,000 population
    violent_crime_rate DECIMAL(10,2),
    homicide_rate DECIMAL(10,2),
    rape_rate DECIMAL(10,2),
    robbery_rate DECIMAL(10,2),
    aggravated_assault_rate DECIMAL(10,2),
    property_crime_rate DECIMAL(10,2),
    burglary_rate DECIMAL(10,2),
    larceny_rate DECIMAL(10,2),
    motor_vehicle_theft_rate DECIMAL(10,2),
    arson_rate DECIMAL(10,2),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT fbi_crime_state_unique UNIQUE (state_abbr, year),
    CONSTRAINT fbi_crime_year_check CHECK (year >= 1960 AND year <= 2030)
);

-- FBI crime statistics by agency
CREATE TABLE security.fbi_crime_agency (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ori VARCHAR(20) NOT NULL, -- Originating Agency Identifier
    agency_name VARCHAR(300),
    agency_type VARCHAR(100),
    state_abbr VARCHAR(2),
    state_name VARCHAR(100),
    division_name VARCHAR(100),
    region_name VARCHAR(100),
    region_desc VARCHAR(200),
    county_name VARCHAR(100),
    msa_name VARCHAR(200),
    population BIGINT,
    population_group VARCHAR(100),
    population_group_code VARCHAR(10),
    year INTEGER NOT NULL,
    
    -- Crime counts (same structure as state table)
    violent_crime BIGINT,
    homicide BIGINT,
    rape BIGINT,
    robbery BIGINT,
    aggravated_assault BIGINT,
    property_crime BIGINT,
    burglary BIGINT,
    larceny BIGINT,
    motor_vehicle_theft BIGINT,
    arson BIGINT,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT fbi_crime_agency_unique UNIQUE (ori, year)
);

-- ============================================================================
-- CENSUS BUREAU TABLES
-- ============================================================================

-- American Community Survey (ACS) data
CREATE TABLE government.census_acs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    geography_type VARCHAR(50) NOT NULL, -- state, county, tract, etc.
    geography_id VARCHAR(50) NOT NULL,
    geography_name VARCHAR(300),
    year INTEGER NOT NULL,
    survey VARCHAR(20) NOT NULL, -- acs1, acs5
    
    -- Demographics
    total_population BIGINT,
    male_population BIGINT,
    female_population BIGINT,
    median_age DECIMAL(5,2),
    
    -- Race and ethnicity
    white_alone BIGINT,
    black_alone BIGINT,
    asian_alone BIGINT,
    hispanic_latino BIGINT,
    
    -- Housing
    total_housing_units BIGINT,
    occupied_housing_units BIGINT,
    vacant_housing_units BIGINT,
    median_home_value DECIMAL(12,2),
    median_rent DECIMAL(8,2),
    
    -- Income
    median_household_income DECIMAL(12,2),
    per_capita_income DECIMAL(12,2),
    poverty_rate DECIMAL(5,2),
    
    -- Education
    high_school_graduate_rate DECIMAL(5,2),
    bachelor_degree_rate DECIMAL(5,2),
    
    -- Employment
    unemployment_rate DECIMAL(5,2),
    labor_force_participation_rate DECIMAL(5,2),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ingestion_job_id UUID REFERENCES ingestion.jobs(id),
    
    CONSTRAINT census_acs_unique UNIQUE (geography_type, geography_id, year, survey)
);

-- ============================================================================
-- TELEMETRY AND MONITORING TABLES
-- ============================================================================

-- Performance metrics
CREATE TABLE telemetry.performance_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_name VARCHAR(100) NOT NULL,
    metric_type VARCHAR(50) NOT NULL, -- counter, gauge, histogram
    metric_value DECIMAL(20,6) NOT NULL,
    labels JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT perf_metric_type_check CHECK (metric_type IN ('counter', 'gauge', 'histogram', 'summary'))
);

-- Benchmark results
CREATE TABLE telemetry.benchmark_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    benchmark_name VARCHAR(200) NOT NULL,
    description TEXT,
    duration_ns BIGINT NOT NULL,
    iterations BIGINT NOT NULL,
    ns_per_op BIGINT,
    mb_per_sec DECIMAL(10,4),
    allocs_per_op BIGINT,
    bytes_per_op BIGINT,
    success_rate DECIMAL(5,2),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- API health checks
CREATE TABLE telemetry.api_health_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    endpoint VARCHAR(100) NOT NULL,
    check_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    response_time_ms INTEGER,
    status_code INTEGER,
    error_message TEXT,
    checked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    CONSTRAINT health_status_check CHECK (status IN ('healthy', 'unhealthy', 'degraded', 'unknown'))
);

-- ============================================================================
-- INDEXES FOR PERFORMANCE
-- ============================================================================

-- Ingestion jobs indexes
CREATE INDEX idx_jobs_status ON ingestion.jobs(status);
CREATE INDEX idx_jobs_endpoint ON ingestion.jobs(endpoint);
CREATE INDEX idx_jobs_created_at ON ingestion.jobs(created_at);
CREATE INDEX idx_jobs_status_endpoint ON ingestion.jobs(status, endpoint);

-- Congress bills indexes
CREATE INDEX idx_congress_bills_congress ON government.congress_bills(congress);
CREATE INDEX idx_congress_bills_type ON government.congress_bills(bill_type);
CREATE INDEX idx_congress_bills_introduced_date ON government.congress_bills(introduced_date);
CREATE INDEX idx_congress_bills_sponsor_state ON government.congress_bills(sponsor_state);
CREATE INDEX idx_congress_bills_policy_area ON government.congress_bills(policy_area);
CREATE INDEX idx_congress_bills_title_gin ON government.congress_bills USING gin(to_tsvector('english', title));

-- FRED observations indexes
CREATE INDEX idx_fred_observations_series_date ON economic.fred_observations(series_id, observation_date);
CREATE INDEX idx_fred_observations_date ON economic.fred_observations(observation_date);
CREATE INDEX idx_fred_observations_realtime ON economic.fred_observations(realtime_start, realtime_end);

-- SEC company facts indexes
CREATE INDEX idx_sec_company_cik ON financial.sec_company_facts(cik);
CREATE INDEX idx_sec_company_name_gin ON financial.sec_company_facts USING gin(to_tsvector('english', entity_name));
CREATE INDEX idx_sec_financial_facts_cik_tag ON financial.sec_financial_facts(cik, tag);
CREATE INDEX idx_sec_financial_facts_end_date ON financial.sec_financial_facts(end_date);

-- FBI crime indexes
CREATE INDEX idx_fbi_crime_state_year ON security.fbi_crime_state(state_abbr, year);
CREATE INDEX idx_fbi_crime_agency_ori_year ON security.fbi_crime_agency(ori, year);
CREATE INDEX idx_fbi_crime_agency_state_year ON security.fbi_crime_agency(state_abbr, year);

-- Census ACS indexes
CREATE INDEX idx_census_acs_geography ON government.census_acs(geography_type, geography_id);
CREATE INDEX idx_census_acs_year ON government.census_acs(year);
CREATE INDEX idx_census_acs_survey ON government.census_acs(survey);

-- Telemetry indexes
CREATE INDEX idx_performance_metrics_name_timestamp ON telemetry.performance_metrics(metric_name, timestamp);
CREATE INDEX idx_benchmark_results_name ON telemetry.benchmark_results(benchmark_name);
CREATE INDEX idx_api_health_endpoint_checked ON telemetry.api_health_checks(endpoint, checked_at);

-- ============================================================================
-- FUNCTIONS AND TRIGGERS
-- ============================================================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers for updated_at
CREATE TRIGGER update_jobs_updated_at BEFORE UPDATE ON ingestion.jobs FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_api_endpoints_updated_at BEFORE UPDATE ON ingestion.api_endpoints FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_congress_bills_updated_at BEFORE UPDATE ON government.congress_bills FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_fred_series_updated_at BEFORE UPDATE ON economic.fred_series FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_sec_company_facts_updated_at BEFORE UPDATE ON financial.sec_company_facts FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate data quality score
CREATE OR REPLACE FUNCTION calculate_data_quality_score(
    p_total_records BIGINT,
    p_valid_records BIGINT,
    p_duplicate_records BIGINT,
    p_null_records BIGINT
) RETURNS DECIMAL(10,4) AS $$
DECLARE
    completeness_score DECIMAL(10,4);
    uniqueness_score DECIMAL(10,4);
    overall_score DECIMAL(10,4);
BEGIN
    -- Avoid division by zero
    IF p_total_records = 0 THEN
        RETURN 0.0;
    END IF;
    
    -- Calculate completeness (non-null records / total records)
    completeness_score := (p_total_records - p_null_records)::DECIMAL / p_total_records * 100;
    
    -- Calculate uniqueness (non-duplicate records / total records)
    uniqueness_score := (p_total_records - p_duplicate_records)::DECIMAL / p_total_records * 100;
    
    -- Overall score is weighted average
    overall_score := (completeness_score * 0.6 + uniqueness_score * 0.4);
    
    RETURN GREATEST(0.0, LEAST(100.0, overall_score));
END;
$$ LANGUAGE plpgsql;

-- ============================================================================
-- INITIAL DATA
-- ============================================================================

-- Insert default API endpoints
INSERT INTO ingestion.api_endpoints (name, base_url, api_version, description, rate_limit, timeout_seconds, auth_type) VALUES
('congress_gov', 'https://api.congress.gov/v3', 'v3', 'Congress.gov API for legislative data', 5000, 30, 'api_key'),
('fred_stlouisfed', 'https://api.stlouisfed.org/fred', 'v1', 'Federal Reserve Economic Data API', 120, 30, 'api_key'),
('sec_edgar', 'https://data.sec.gov', 'v1', 'SEC EDGAR database API', 10, 30, 'none'),
('fbi_crime', 'https://api.usa.gov/crime/fbi/sapi', 'v1', 'FBI Crime Data API', 1000, 30, 'api_key'),
('census_gov', 'https://api.census.gov/data', 'v1', 'U.S. Census Bureau API', 500, 30, 'api_key');

-- Create materialized views for common queries
CREATE MATERIALIZED VIEW government.congress_bills_summary AS
SELECT 
    congress,
    bill_type,
    COUNT(*) as total_bills,
    COUNT(CASE WHEN latest_action_date >= CURRENT_DATE - INTERVAL '30 days' THEN 1 END) as recent_activity,
    sponsor_party,
    COUNT(DISTINCT sponsor_state) as states_represented
FROM government.congress_bills 
GROUP BY congress, bill_type, sponsor_party;

CREATE UNIQUE INDEX idx_congress_bills_summary ON government.congress_bills_summary(congress, bill_type, sponsor_party);

-- Grant permissions
GRANT USAGE ON SCHEMA ingestion, government, economic, financial, security, telemetry TO PUBLIC;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA ingestion, government, economic, financial, security, telemetry TO PUBLIC;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA ingestion, government, economic, financial, security, telemetry TO PUBLIC;

-- Create indexes for full-text search
CREATE INDEX idx_congress_bills_fts ON government.congress_bills USING gin(
    to_tsvector('english', COALESCE(title, '') || ' ' || COALESCE(short_title, '') || ' ' || COALESCE(policy_area, ''))
);

CREATE INDEX idx_sec_company_fts ON financial.sec_company_facts USING gin(
    to_tsvector('english', COALESCE(entity_name, '') || ' ' || COALESCE(description, ''))
);

-- Migration complete
SELECT 'Initial schema migration completed successfully' as status;
