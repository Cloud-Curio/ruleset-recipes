/*
 * Data Ingestion TUI - Database Management
 * 
 * This package handles database connections and operations for the
 * data ingestion TUI, supporting both SQLite and PostgreSQL.
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// @decorator: Database
// @description: Database connection wrapper
type Database struct {
	DB     *sqlx.DB
	Driver string
	logger *logrus.Logger
}

// @decorator: Config
// @description: Database configuration
type Config struct {
	Driver   string
	URL      string
	MaxConns int
	MaxIdle  int
	Timeout  time.Duration
}

// @decorator: NewDatabase
// @description: Create a new database connection
func NewDatabase(config Config, logger *logrus.Logger) (*Database, error) {
	if logger == nil {
		logger = logrus.New()
	}

	db, err := sqlx.Connect(config.Driver, config.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxConns)
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetConnMaxLifetime(config.Timeout)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"driver":    config.Driver,
		"max_conns": config.MaxConns,
		"max_idle":  config.MaxIdle,
	}).Info("Database connection established")

	return &Database{
		DB:     db,
		Driver: config.Driver,
		logger: logger,
	}, nil
}

// @decorator: Close
// @description: Close database connection
func (d *Database) Close() error {
	if d.DB != nil {
		d.logger.Info("Closing database connection")
		return d.DB.Close()
	}
	return nil
}

// @decorator: Ping
// @description: Test database connection
func (d *Database) Ping() error {
	return d.DB.Ping()
}

// @decorator: Exec
// @description: Execute a query without returning rows
func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}

// @decorator: Query
// @description: Execute a query that returns rows
func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

// @decorator: QueryRow
// @description: Execute a query that returns a single row
func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

// @decorator: Select
// @description: Execute a query and scan results into a slice
func (d *Database) Select(dest interface{}, query string, args ...interface{}) error {
	return d.DB.Select(dest, query, args...)
}

// @decorator: Get
// @description: Execute a query and scan result into a single struct
func (d *Database) Get(dest interface{}, query string, args ...interface{}) error {
	return d.DB.Get(dest, query, args...)
}

// @decorator: NamedExec
// @description: Execute a named query
func (d *Database) NamedExec(query string, arg interface{}) (sql.Result, error) {
	return d.DB.NamedExec(query, arg)
}

// @decorator: NamedQuery
// @description: Execute a named query that returns rows
func (d *Database) NamedQuery(query string, arg interface{}) (*sqlx.Rows, error) {
	return d.DB.NamedQuery(query, arg)
}

// @decorator: BeginTx
// @description: Begin a database transaction
func (d *Database) BeginTx() (*sqlx.Tx, error) {
	return d.DB.Beginx()
}

// @decorator: IsHealthy
// @description: Check if database is healthy
func (d *Database) IsHealthy() bool {
	return d.Ping() == nil
}

// @decorator: GetStats
// @description: Get database connection statistics
func (d *Database) GetStats() sql.DBStats {
	return d.DB.Stats()
}
