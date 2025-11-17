/*
 * Data Ingestion TUI - Main Entry Point
 * 
 * This file serves as the main entry point for the Data Ingestion TUI application.
 * It initializes the application, sets up configuration, and launches the TUI.
 * 
 * Features:
 * - Command-line interface with Cobra
 * - Configuration management
 * - Graceful shutdown handling
 * - Logging setup
 * 
 * Author: Codegen AI Assistant
 * Created: 2024
 * License: MIT
 */

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"

	"data-ingestion-tui/internal/app"
	"data-ingestion-tui/internal/config"
	"data-ingestion-tui/internal/database"
	"data-ingestion-tui/internal/ui"
	"data-ingestion-tui/internal/utils"
)

// @decorator: main
// @description: Application entry point with comprehensive error handling and graceful shutdown
// @version: 1.0.0
// @author: Codegen AI Assistant

var (
	// Application version information
	version   = "1.0.0"
	buildTime = "unknown"
	gitCommit = "unknown"
	
	// Global configuration
	cfg *config.Config
	
	// Logger instance
	logger *logrus.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dit",
	Short: "Data Ingestion TUI - A comprehensive data management interface",
	Long: `Data Ingestion TUI (dit) is a powerful terminal user interface for managing
data ingestion from various government and public APIs.

Features:
- Multi-source data ingestion (Congress.gov, FRED, SEC, etc.)
- AI-powered assistance with OpenRouter integration
- Secure API key management
- SQL migration management
- Job scheduling and monitoring
- Interactive debugging tools`,
	Version: fmt.Sprintf("%s (built: %s, commit: %s)", version, buildTime, gitCommit),
	RunE:    runTUI,
}

// initCmd represents the init command for first-time setup
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the application configuration",
	Long: `Initialize the Data Ingestion TUI application by creating necessary
configuration files, database schema, and setting up the environment.`,
	RunE: runInit,
}

// serverCmd represents the server command for SSH access
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the SSH server for remote access",
	Long: `Start an SSH server that allows remote access to the TUI interface.
This enables team collaboration and remote management capabilities.`,
	RunE: runServer,
}

// @decorator: init
// @description: Initialize command-line flags and configuration
func init() {
	// Add persistent flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.dit.yaml)")
	rootCmd.PersistentFlags().StringP("log-level", "l", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringP("data-dir", "d", "", "data directory (default is $HOME/.dit)")
	
	// Add subcommands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(serverCmd)
	
	// Server-specific flags
	serverCmd.Flags().StringP("host", "H", "localhost", "server host")
	serverCmd.Flags().IntP("port", "p", 2222, "server port")
	serverCmd.Flags().StringP("key-path", "k", "", "SSH private key path")
}

// @decorator: main
// @description: Main function with comprehensive error handling and graceful shutdown
func main() {
	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	// Handle interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-sigChan
		logger.Info("Received interrupt signal, shutting down gracefully...")
		cancel()
	}()
	
	// Execute the root command
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// @decorator: runTUI
// @description: Main TUI application runner with error handling and cleanup
// @param cmd: Cobra command instance
// @param args: Command arguments
// @return error: Any error that occurred during execution
func runTUI(cmd *cobra.Command, args []string) error {
	var err error
	
	// Initialize configuration
	if cfg, err = initializeConfig(cmd); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}
	
	// Initialize logger
	if logger, err = initializeLogger(cfg); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	
	logger.Info("Starting Data Ingestion TUI...")
	
	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			logger.WithError(closeErr).Error("Failed to close database connection")
		}
	}()
	
	// Initialize application
	application, err := app.New(cfg, db, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize application")
		return fmt.Errorf("application initialization failed: %w", err)
	}
	
	// Initialize and run TUI
	tui, err := ui.New(application, cfg.UI)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize TUI")
		return fmt.Errorf("TUI initialization failed: %w", err)
	}
	
	logger.Info("Launching TUI interface...")
	
	// Run the TUI with context for graceful shutdown
	if err := tui.Run(cmd.Context()); err != nil {
		logger.WithError(err).Error("TUI execution failed")
		return fmt.Errorf("TUI execution failed: %w", err)
	}
	
	logger.Info("Data Ingestion TUI shutdown complete")
	return nil
}

// @decorator: runInit
// @description: Initialize application configuration and setup
// @param cmd: Cobra command instance
// @param args: Command arguments
// @return error: Any error that occurred during initialization
func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("🚀 Initializing Data Ingestion TUI...")
	
	// Initialize configuration
	cfg, err := initializeConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}
	
	// Initialize logger
	logger, err := initializeLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	
	// Create necessary directories
	if err := utils.CreateDirectories(cfg.DataDir); err != nil {
		logger.WithError(err).Error("Failed to create directories")
		return fmt.Errorf("directory creation failed: %w", err)
	}
	
	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer db.Close()
	
	// Run database migrations
	if err := database.RunMigrations(cfg.Database.URL); err != nil {
		logger.WithError(err).Error("Failed to run database migrations")
		return fmt.Errorf("database migration failed: %w", err)
	}
	
	fmt.Println("✅ Initialization complete!")
	fmt.Println("📝 Configuration saved to:", cfg.ConfigFile)
	fmt.Println("🗄️  Database initialized at:", cfg.Database.URL)
	fmt.Println("📁 Data directory:", cfg.DataDir)
	fmt.Println("\n🎉 You can now run 'dit' to start the TUI!")
	
	return nil
}

// @decorator: runServer
// @description: Start SSH server for remote TUI access
// @param cmd: Cobra command instance
// @param args: Command arguments
// @return error: Any error that occurred during server startup
func runServer(cmd *cobra.Command, args []string) error {
	// Get server configuration from flags
	host, _ := cmd.Flags().GetString("host")
	port, _ := cmd.Flags().GetInt("port")
	keyPath, _ := cmd.Flags().GetString("key-path")
	
	fmt.Printf("🌐 Starting SSH server on %s:%d...\n", host, port)
	
	// Initialize configuration
	cfg, err := initializeConfig(cmd)
	if err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}
	
	// Initialize logger
	logger, err := initializeLogger(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	
	// Initialize database
	db, err := database.Initialize(cfg.Database)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		return fmt.Errorf("database initialization failed: %w", err)
	}
	defer db.Close()
	
	// Initialize application
	application, err := app.New(cfg, db, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize application")
		return fmt.Errorf("application initialization failed: %w", err)
	}
	
	// Start SSH server
	server, err := ui.NewSSHServer(application, host, port, keyPath)
	if err != nil {
		logger.WithError(err).Error("Failed to create SSH server")
		return fmt.Errorf("SSH server creation failed: %w", err)
	}
	
	logger.Infof("SSH server listening on %s:%d", host, port)
	fmt.Printf("✅ SSH server ready! Connect with: ssh -p %d %s\n", port, host)
	
	// Run server with context for graceful shutdown
	return server.Run(cmd.Context())
}

// @decorator: initializeConfig
// @description: Initialize application configuration from various sources
// @param cmd: Cobra command instance
// @return *config.Config: Initialized configuration
// @return error: Any error that occurred during configuration initialization
func initializeConfig(cmd *cobra.Command) (*config.Config, error) {
	configFile, _ := cmd.Flags().GetString("config")
	dataDir, _ := cmd.Flags().GetString("data-dir")
	logLevel, _ := cmd.Flags().GetString("log-level")
	verbose, _ := cmd.Flags().GetBool("verbose")
	
	cfg, err := config.Load(configFile, dataDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Override with command-line flags
	if logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if verbose {
		cfg.Verbose = verbose
	}
	
	return cfg, nil
}

// @decorator: initializeLogger
// @description: Initialize structured logger with configuration
// @param cfg: Application configuration
// @return *logrus.Logger: Configured logger instance
// @return error: Any error that occurred during logger initialization
func initializeLogger(cfg *config.Config) (*logrus.Logger, error) {
	logger := logrus.New()
	
	// Set log level
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		return nil, fmt.Errorf("invalid log level '%s': %w", cfg.LogLevel, err)
	}
	logger.SetLevel(level)
	
	// Set formatter
	if cfg.Verbose {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
	} else {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	
	// Set output
	logger.SetOutput(os.Stdout)
	
	return logger, nil
}
