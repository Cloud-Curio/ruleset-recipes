/*
 * Data Ingestion TUI - Go Module Definition
 * 
 * A comprehensive TUI application for managing data ingestion from various
 * government and public APIs using the Bubble Tea ecosystem.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

module data-ingestion-tui

go 1.21

require (
	// Bubble Tea ecosystem
	github.com/charmbracelet/bubbletea v0.25.0
	github.com/charmbracelet/lipgloss v0.9.1
	github.com/charmbracelet/wish v1.2.0
	github.com/charmbracelet/harmonica v0.2.0
	github.com/charmbracelet/huh v0.2.3
	github.com/charmbracelet/log v0.3.1
	
	// Database and storage
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.18
	github.com/golang-migrate/migrate/v4 v4.16.2
	
	// HTTP and API clients
	github.com/go-resty/resty/v2 v2.10.0
	github.com/gorilla/websocket v1.5.1
	
	// Configuration and secrets management
	github.com/spf13/viper v1.17.0
	github.com/zalando/go-keyring v0.2.3
	
	// Scheduling and cron
	github.com/robfig/cron/v3 v3.0.1
	
	// AI and OpenRouter integration
	github.com/sashabaranov/go-openai v1.17.9
	
	// Utilities
	github.com/google/uuid v1.4.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
	
	// JSON and data processing
	github.com/tidwall/gjson v1.17.0
	github.com/tidwall/sjson v1.2.5
	
	// File system operations
	github.com/fsnotify/fsnotify v1.7.0
	
	// Encryption and security
	golang.org/x/crypto v0.15.0
	
	// Time and date utilities
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
)
