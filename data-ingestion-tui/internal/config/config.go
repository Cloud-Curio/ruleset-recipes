/*
 * Data Ingestion TUI - Configuration Management
 * 
 * This package handles all configuration management for the application,
 * including loading from files, environment variables, and defaults.
 * 
 * Features:
 * - Hierarchical configuration loading
 * - Environment variable support
 * - Validation and defaults
 * - Hot-reloading capabilities
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

// @decorator: Config
// @description: Main configuration structure for the application
// @version: 1.0.0
// @author: Codegen AI Assistant

// Config represents the main application configuration
type Config struct {
	// Application settings
	AppName    string `mapstructure:"app_name" yaml:"app_name"`
	Version    string `mapstructure:"version" yaml:"version"`
	DataDir    string `mapstructure:"data_dir" yaml:"data_dir"`
	ConfigFile string `mapstructure:"config_file" yaml:"config_file"`
	
	// Logging configuration
	LogLevel string `mapstructure:"log_level" yaml:"log_level"`
	Verbose  bool   `mapstructure:"verbose" yaml:"verbose"`
	
	// Database configuration
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	
	// UI configuration
	UI UIConfig `mapstructure:"ui" yaml:"ui"`
	
	// API configuration
	APIs APIConfig `mapstructure:"apis" yaml:"apis"`
	
	// AI configuration
	AI AIConfig `mapstructure:"ai" yaml:"ai"`
	
	// Scheduler configuration
	Scheduler SchedulerConfig `mapstructure:"scheduler" yaml:"scheduler"`
	
	// Security configuration
	Security SecurityConfig `mapstructure:"security" yaml:"security"`
	
	// Server configuration (for SSH server)
	Server ServerConfig `mapstructure:"server" yaml:"server"`
}

// @decorator: DatabaseConfig
// @description: Database connection and migration settings
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" yaml:"driver"`
	URL             string        `mapstructure:"url" yaml:"url"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	MigrationsPath  string        `mapstructure:"migrations_path" yaml:"migrations_path"`
	AutoMigrate     bool          `mapstructure:"auto_migrate" yaml:"auto_migrate"`
}

// @decorator: UIConfig
// @description: Terminal UI appearance and behavior settings
type UIConfig struct {
	Theme           string        `mapstructure:"theme" yaml:"theme"`
	RefreshInterval time.Duration `mapstructure:"refresh_interval" yaml:"refresh_interval"`
	PageSize        int           `mapstructure:"page_size" yaml:"page_size"`
	EnableMouse     bool          `mapstructure:"enable_mouse" yaml:"enable_mouse"`
	EnableAltScreen bool          `mapstructure:"enable_alt_screen" yaml:"enable_alt_screen"`
	Colors          ColorConfig   `mapstructure:"colors" yaml:"colors"`
}

// @decorator: ColorConfig
// @description: Color scheme configuration for the TUI
type ColorConfig struct {
	Primary     string `mapstructure:"primary" yaml:"primary"`
	Secondary   string `mapstructure:"secondary" yaml:"secondary"`
	Success     string `mapstructure:"success" yaml:"success"`
	Warning     string `mapstructure:"warning" yaml:"warning"`
	Error       string `mapstructure:"error" yaml:"error"`
	Background  string `mapstructure:"background" yaml:"background"`
	Foreground  string `mapstructure:"foreground" yaml:"foreground"`
	Border      string `mapstructure:"border" yaml:"border"`
}

// @decorator: APIConfig
// @description: Configuration for various API endpoints and clients
type APIConfig struct {
	Timeout         time.Duration           `mapstructure:"timeout" yaml:"timeout"`
	RetryAttempts   int                     `mapstructure:"retry_attempts" yaml:"retry_attempts"`
	RetryDelay      time.Duration           `mapstructure:"retry_delay" yaml:"retry_delay"`
	RateLimiting    RateLimitConfig         `mapstructure:"rate_limiting" yaml:"rate_limiting"`
	Endpoints       map[string]EndpointConfig `mapstructure:"endpoints" yaml:"endpoints"`
}

// @decorator: RateLimitConfig
// @description: Rate limiting configuration for API calls
type RateLimitConfig struct {
	Enabled     bool          `mapstructure:"enabled" yaml:"enabled"`
	RequestsPerSecond int     `mapstructure:"requests_per_second" yaml:"requests_per_second"`
	BurstSize   int           `mapstructure:"burst_size" yaml:"burst_size"`
	Timeout     time.Duration `mapstructure:"timeout" yaml:"timeout"`
}

// @decorator: EndpointConfig
// @description: Configuration for individual API endpoints
type EndpointConfig struct {
	BaseURL     string            `mapstructure:"base_url" yaml:"base_url"`
	APIKey      string            `mapstructure:"api_key" yaml:"api_key"`
	Headers     map[string]string `mapstructure:"headers" yaml:"headers"`
	Timeout     time.Duration     `mapstructure:"timeout" yaml:"timeout"`
	Enabled     bool              `mapstructure:"enabled" yaml:"enabled"`
	RateLimit   int               `mapstructure:"rate_limit" yaml:"rate_limit"`
	Description string            `mapstructure:"description" yaml:"description"`
}

// @decorator: AIConfig
// @description: AI integration and OpenRouter configuration
type AIConfig struct {
	Provider        string            `mapstructure:"provider" yaml:"provider"`
	APIKey          string            `mapstructure:"api_key" yaml:"api_key"`
	BaseURL         string            `mapstructure:"base_url" yaml:"base_url"`
	Model           string            `mapstructure:"model" yaml:"model"`
	Temperature     float32           `mapstructure:"temperature" yaml:"temperature"`
	MaxTokens       int               `mapstructure:"max_tokens" yaml:"max_tokens"`
	Timeout         time.Duration     `mapstructure:"timeout" yaml:"timeout"`
	RAG             RAGConfig         `mapstructure:"rag" yaml:"rag"`
	MCP             MCPConfig         `mapstructure:"mcp" yaml:"mcp"`
	Tools           []string          `mapstructure:"tools" yaml:"tools"`
	CustomModels    map[string]string `mapstructure:"custom_models" yaml:"custom_models"`
}

// @decorator: RAGConfig
// @description: Retrieval-Augmented Generation configuration
type RAGConfig struct {
	Enabled         bool   `mapstructure:"enabled" yaml:"enabled"`
	VectorDB        string `mapstructure:"vector_db" yaml:"vector_db"`
	EmbeddingModel  string `mapstructure:"embedding_model" yaml:"embedding_model"`
	ChunkSize       int    `mapstructure:"chunk_size" yaml:"chunk_size"`
	ChunkOverlap    int    `mapstructure:"chunk_overlap" yaml:"chunk_overlap"`
	TopK            int    `mapstructure:"top_k" yaml:"top_k"`
	IndexPath       string `mapstructure:"index_path" yaml:"index_path"`
}

// @decorator: MCPConfig
// @description: Model Context Protocol server configuration
type MCPConfig struct {
	Enabled bool                    `mapstructure:"enabled" yaml:"enabled"`
	Servers map[string]MCPServer    `mapstructure:"servers" yaml:"servers"`
}

// @decorator: MCPServer
// @description: Individual MCP server configuration
type MCPServer struct {
	Command     []string          `mapstructure:"command" yaml:"command"`
	Args        []string          `mapstructure:"args" yaml:"args"`
	Env         map[string]string `mapstructure:"env" yaml:"env"`
	Enabled     bool              `mapstructure:"enabled" yaml:"enabled"`
	Description string            `mapstructure:"description" yaml:"description"`
}

// @decorator: SchedulerConfig
// @description: Job scheduling and cron configuration
type SchedulerConfig struct {
	Enabled         bool          `mapstructure:"enabled" yaml:"enabled"`
	MaxConcurrent   int           `mapstructure:"max_concurrent" yaml:"max_concurrent"`
	DefaultTimeout  time.Duration `mapstructure:"default_timeout" yaml:"default_timeout"`
	RetryAttempts   int           `mapstructure:"retry_attempts" yaml:"retry_attempts"`
	RetryDelay      time.Duration `mapstructure:"retry_delay" yaml:"retry_delay"`
	PersistJobs     bool          `mapstructure:"persist_jobs" yaml:"persist_jobs"`
}

// @decorator: SecurityConfig
// @description: Security and encryption settings
type SecurityConfig struct {
	EncryptionKey   string `mapstructure:"encryption_key" yaml:"encryption_key"`
	KeyringService  string `mapstructure:"keyring_service" yaml:"keyring_service"`
	KeyringUser     string `mapstructure:"keyring_user" yaml:"keyring_user"`
	TLSEnabled      bool   `mapstructure:"tls_enabled" yaml:"tls_enabled"`
	CertFile        string `mapstructure:"cert_file" yaml:"cert_file"`
	KeyFile         string `mapstructure:"key_file" yaml:"key_file"`
}

// @decorator: ServerConfig
// @description: SSH server configuration for remote access
type ServerConfig struct {
	Host        string        `mapstructure:"host" yaml:"host"`
	Port        int           `mapstructure:"port" yaml:"port"`
	KeyPath     string        `mapstructure:"key_path" yaml:"key_path"`
	MaxClients  int           `mapstructure:"max_clients" yaml:"max_clients"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout" yaml:"idle_timeout"`
	Banner      string        `mapstructure:"banner" yaml:"banner"`
}

// @decorator: Load
// @description: Load configuration from file, environment, and defaults
// @param configFile: Path to configuration file (optional)
// @param dataDir: Data directory path (optional)
// @return *Config: Loaded configuration
// @return error: Any error that occurred during loading
func Load(configFile, dataDir string) (*Config, error) {
	// Initialize Viper
	v := viper.New()
	
	// Set defaults
	setDefaults(v, dataDir)
	
	// Set up configuration file
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Look for config in standard locations
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get user home directory")
		}
		
		v.SetConfigName(".dit")
		v.SetConfigType("yaml")
		v.AddConfigPath(home)
		v.AddConfigPath(".")
		v.AddConfigPath("/etc/dit/")
	}
	
	// Set up environment variables
	v.SetEnvPrefix("DIT")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	
	// Read configuration file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, errors.Wrap(err, "failed to read configuration file")
		}
		// Config file not found is OK, we'll use defaults
	}
	
	// Unmarshal configuration
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal configuration")
	}
	
	// Set the config file path
	if v.ConfigFileUsed() != "" {
		cfg.ConfigFile = v.ConfigFileUsed()
	} else {
		home, _ := os.UserHomeDir()
		cfg.ConfigFile = filepath.Join(home, ".dit.yaml")
	}
	
	// Validate configuration
	if err := validate(&cfg); err != nil {
		return nil, errors.Wrap(err, "configuration validation failed")
	}
	
	return &cfg, nil
}

// @decorator: setDefaults
// @description: Set default configuration values
// @param v: Viper instance
// @param dataDir: Data directory path
func setDefaults(v *viper.Viper, dataDir string) {
	// Determine data directory
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".dit")
	}
	
	// Application defaults
	v.SetDefault("app_name", "Data Ingestion TUI")
	v.SetDefault("version", "1.0.0")
	v.SetDefault("data_dir", dataDir)
	v.SetDefault("log_level", "info")
	v.SetDefault("verbose", false)
	
	// Database defaults
	v.SetDefault("database.driver", "sqlite3")
	v.SetDefault("database.url", filepath.Join(dataDir, "dit.db"))
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", "5m")
	v.SetDefault("database.migrations_path", "migrations")
	v.SetDefault("database.auto_migrate", true)
	
	// UI defaults
	v.SetDefault("ui.theme", "default")
	v.SetDefault("ui.refresh_interval", "1s")
	v.SetDefault("ui.page_size", 20)
	v.SetDefault("ui.enable_mouse", true)
	v.SetDefault("ui.enable_alt_screen", true)
	
	// Color defaults
	v.SetDefault("ui.colors.primary", "#7C3AED")
	v.SetDefault("ui.colors.secondary", "#10B981")
	v.SetDefault("ui.colors.success", "#059669")
	v.SetDefault("ui.colors.warning", "#D97706")
	v.SetDefault("ui.colors.error", "#DC2626")
	v.SetDefault("ui.colors.background", "#1F2937")
	v.SetDefault("ui.colors.foreground", "#F9FAFB")
	v.SetDefault("ui.colors.border", "#6B7280")
	
	// API defaults
	v.SetDefault("apis.timeout", "30s")
	v.SetDefault("apis.retry_attempts", 3)
	v.SetDefault("apis.retry_delay", "1s")
	v.SetDefault("apis.rate_limiting.enabled", true)
	v.SetDefault("apis.rate_limiting.requests_per_second", 10)
	v.SetDefault("apis.rate_limiting.burst_size", 20)
	v.SetDefault("apis.rate_limiting.timeout", "1m")
	
	// AI defaults
	v.SetDefault("ai.provider", "openrouter")
	v.SetDefault("ai.base_url", "https://openrouter.ai/api/v1")
	v.SetDefault("ai.model", "anthropic/claude-3-haiku")
	v.SetDefault("ai.temperature", 0.7)
	v.SetDefault("ai.max_tokens", 4096)
	v.SetDefault("ai.timeout", "60s")
	
	// RAG defaults
	v.SetDefault("ai.rag.enabled", true)
	v.SetDefault("ai.rag.vector_db", "sqlite")
	v.SetDefault("ai.rag.embedding_model", "text-embedding-ada-002")
	v.SetDefault("ai.rag.chunk_size", 1000)
	v.SetDefault("ai.rag.chunk_overlap", 200)
	v.SetDefault("ai.rag.top_k", 5)
	v.SetDefault("ai.rag.index_path", filepath.Join(dataDir, "embeddings"))
	
	// MCP defaults
	v.SetDefault("ai.mcp.enabled", true)
	
	// Scheduler defaults
	v.SetDefault("scheduler.enabled", true)
	v.SetDefault("scheduler.max_concurrent", 5)
	v.SetDefault("scheduler.default_timeout", "30m")
	v.SetDefault("scheduler.retry_attempts", 3)
	v.SetDefault("scheduler.retry_delay", "5m")
	v.SetDefault("scheduler.persist_jobs", true)
	
	// Security defaults
	v.SetDefault("security.keyring_service", "dit")
	v.SetDefault("security.keyring_user", "default")
	v.SetDefault("security.tls_enabled", false)
	
	// Server defaults
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 2222)
	v.SetDefault("server.max_clients", 10)
	v.SetDefault("server.idle_timeout", "30m")
	v.SetDefault("server.banner", "Welcome to Data Ingestion TUI")
}

// @decorator: validate
// @description: Validate configuration values
// @param cfg: Configuration to validate
// @return error: Validation error if any
func validate(cfg *Config) error {
	// Validate data directory
	if cfg.DataDir == "" {
		return errors.New("data directory cannot be empty")
	}
	
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		return errors.Wrap(err, "failed to create data directory")
	}
	
	// Validate database configuration
	if cfg.Database.Driver == "" {
		return errors.New("database driver cannot be empty")
	}
	
	if cfg.Database.URL == "" {
		return errors.New("database URL cannot be empty")
	}
	
	// Validate UI configuration
	if cfg.UI.PageSize <= 0 {
		return errors.New("UI page size must be positive")
	}
	
	// Validate API configuration
	if cfg.APIs.Timeout <= 0 {
		return errors.New("API timeout must be positive")
	}
	
	// Validate AI configuration
	if cfg.AI.MaxTokens <= 0 {
		return errors.New("AI max tokens must be positive")
	}
	
	if cfg.AI.Temperature < 0 || cfg.AI.Temperature > 2 {
		return errors.New("AI temperature must be between 0 and 2")
	}
	
	// Validate scheduler configuration
	if cfg.Scheduler.MaxConcurrent <= 0 {
		return errors.New("scheduler max concurrent jobs must be positive")
	}
	
	// Validate server configuration
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return errors.New("server port must be between 1 and 65535")
	}
	
	return nil
}

// @decorator: Save
// @description: Save configuration to file
// @param cfg: Configuration to save
// @return error: Any error that occurred during saving
func Save(cfg *Config) error {
	v := viper.New()
	v.SetConfigFile(cfg.ConfigFile)
	v.SetConfigType("yaml")
	
	// Convert config to map for viper
	configMap := make(map[string]interface{})
	
	// This is a simplified version - in a real implementation,
	// you'd use reflection or a more sophisticated method
	configMap["app_name"] = cfg.AppName
	configMap["version"] = cfg.Version
	configMap["data_dir"] = cfg.DataDir
	configMap["log_level"] = cfg.LogLevel
	configMap["verbose"] = cfg.Verbose
	
	// Set all values in viper
	for key, value := range configMap {
		v.Set(key, value)
	}
	
	// Write configuration file
	if err := v.WriteConfig(); err != nil {
		return errors.Wrap(err, "failed to write configuration file")
	}
	
	return nil
}

// @decorator: GetAPIEndpoints
// @description: Get list of configured API endpoints
// @param cfg: Configuration instance
// @return []string: List of endpoint names
func GetAPIEndpoints(cfg *Config) []string {
	endpoints := make([]string, 0, len(cfg.APIs.Endpoints))
	for name := range cfg.APIs.Endpoints {
		endpoints = append(endpoints, name)
	}
	return endpoints
}

// @decorator: GetEnabledEndpoints
// @description: Get list of enabled API endpoints
// @param cfg: Configuration instance
// @return []string: List of enabled endpoint names
func GetEnabledEndpoints(cfg *Config) []string {
	endpoints := make([]string, 0)
	for name, endpoint := range cfg.APIs.Endpoints {
		if endpoint.Enabled {
			endpoints = append(endpoints, name)
		}
	}
	return endpoints
}
