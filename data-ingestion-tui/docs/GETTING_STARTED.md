# Getting Started with Data Ingestion TUI 🚀

Welcome to the Data Ingestion TUI! This comprehensive guide will help you get up and running with the most powerful terminal-based data management interface for government and public APIs.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [First Run](#first-run)
- [Configuration](#configuration)
- [API Keys Setup](#api-keys-setup)
- [Basic Usage](#basic-usage)
- [AI Assistant Setup](#ai-assistant-setup)
- [Scheduling Jobs](#scheduling-jobs)
- [Troubleshooting](#troubleshooting)
- [Next Steps](#next-steps)

## Prerequisites

Before you begin, ensure you have the following installed:

### Required
- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **SQLite3** - Usually pre-installed on most systems

### Optional (for development)
- **Make** - For using the Makefile commands
- **Docker** - For containerized deployment
- **Air** - For hot reloading during development

### System Requirements
- **Memory**: 512MB RAM minimum, 1GB recommended
- **Storage**: 100MB for application, additional space for data
- **Network**: Internet connection for API access
- **Terminal**: Modern terminal with Unicode support

## Installation

### Option 1: Build from Source (Recommended)

```bash
# Clone the repository
git clone https://github.com/your-org/data-ingestion-tui.git
cd data-ingestion-tui

# Install dependencies
go mod tidy

# Build the application
make build

# Or build manually
go build -o bin/dit cmd/main.go
```

### Option 2: Using Go Install

```bash
go install github.com/your-org/data-ingestion-tui/cmd@latest
```

### Option 3: Docker

```bash
# Pull the image
docker pull your-org/data-ingestion-tui:latest

# Or build locally
docker build -t data-ingestion-tui .
```

### Option 4: Pre-built Binaries

Download the latest release from the [releases page](https://github.com/your-org/data-ingestion-tui/releases).

## First Run

### Initialize the Application

```bash
# Initialize configuration and database
./bin/dit init

# Or if installed globally
dit init
```

This will:
- Create the configuration directory (`~/.dit/`)
- Initialize the SQLite database
- Create default configuration files
- Set up the data directory structure

### Launch the TUI

```bash
# Start the interactive TUI
./bin/dit

# Or with verbose logging
./bin/dit --verbose --log-level debug
```

### SSH Server Mode

```bash
# Start SSH server for remote access
./bin/dit server --host 0.0.0.0 --port 2222
```

## Configuration

### Configuration File Location

The application looks for configuration in the following order:
1. `--config` flag
2. `DIT_CONFIG_FILE` environment variable
3. `~/.dit.yaml`
4. `./dit.yaml`
5. `/etc/dit/dit.yaml`

### Basic Configuration

Create or edit `~/.dit.yaml`:

```yaml
# Application settings
app_name: "Data Ingestion TUI"
data_dir: "~/.dit"
log_level: "info"
verbose: false

# Database configuration
database:
  driver: "sqlite3"
  url: "~/.dit/dit.db"
  auto_migrate: true

# UI configuration
ui:
  theme: "default"
  refresh_interval: "1s"
  page_size: 20
  enable_mouse: true
  colors:
    primary: "#7C3AED"
    secondary: "#10B981"
    success: "#059669"
    warning: "#D97706"
    error: "#DC2626"

# AI configuration
ai:
  provider: "openrouter"
  base_url: "https://openrouter.ai/api/v1"
  model: "anthropic/claude-3-haiku"
  temperature: 0.7
  max_tokens: 4096
  rag:
    enabled: true
    chunk_size: 1000
    top_k: 5
  mcp:
    enabled: true

# Scheduler configuration
scheduler:
  enabled: true
  max_concurrent: 5
  default_timeout: "30m"
```

### Environment Variables

You can override any configuration using environment variables:

```bash
export DIT_LOG_LEVEL=debug
export DIT_DATA_DIR=/custom/path
export DIT_AI_API_KEY=your-openrouter-key
export DIT_CONGRESS_API_KEY=your-congress-key
```

## API Keys Setup

The application supports multiple APIs. Here's how to set up the most important ones:

### 1. Congress.gov API

```bash
# Get your API key from https://api.congress.gov/sign-up/
dit config set apis.endpoints.congress_gov.api_key "YOUR_API_KEY"

# Or set via environment
export DIT_APIS_ENDPOINTS_CONGRESS_GOV_API_KEY="YOUR_API_KEY"
```

### 2. OpenRouter (AI Assistant)

```bash
# Get your API key from https://openrouter.ai/
dit config set ai.api_key "YOUR_OPENROUTER_KEY"

# Or set via environment
export DIT_AI_API_KEY="YOUR_OPENROUTER_KEY"
```

### 3. FRED (Economic Data)

```bash
# Get your API key from https://fred.stlouisfed.org/docs/api/api_key.html
dit config set apis.endpoints.fred_stlouisfed.api_key "YOUR_FRED_KEY"
```

### 4. SEC.gov (No API key required)

The SEC API doesn't require an API key, but you must include contact information in the User-Agent header:

```bash
dit config set apis.endpoints.sec_gov.headers.User-Agent "DataIngestionTUI/1.0 (your-email@example.com)"
```

### 5. Alpha Vantage (Stock Data)

```bash
# Get your free API key from https://www.alphavantage.co/support/#api-key
dit config set apis.endpoints.alpha_vantage.api_key "YOUR_ALPHA_VANTAGE_KEY"
```

### Secure Key Storage

The application uses your system's keyring for secure storage:

- **macOS**: Keychain
- **Linux**: Secret Service API (GNOME Keyring, KDE Wallet)
- **Windows**: Windows Credential Manager

Keys are automatically encrypted and stored securely.

## Basic Usage

### Navigation

- **Tab** / **Shift+Tab**: Navigate between panes
- **F1-F8**: Direct pane navigation
- **Ctrl+C** / **Q**: Quit application
- **Arrow Keys**: Navigate within panes
- **Enter**: Select/activate items
- **Esc**: Go back/cancel

### Main Panes

1. **📊 Dashboard** (F2): Overview of system status and recent activity
2. **🔗 API Endpoints** (F3): Manage and test API connections
3. **⚡ Ingestion Jobs** (F4): Monitor and manage data ingestion tasks
4. **⏰ Scheduler** (F5): Set up and manage scheduled jobs
5. **🤖 AI Assistant** (F6): Chat with AI for help and automation
6. **⚙️ Configuration** (F7): Manage application settings
7. **📝 Logs** (F8): View application logs and debug information
8. **❓ Help** (F1): Get help and view keyboard shortcuts

### Quick Start Workflow

1. **Set up API keys** (Configuration pane)
2. **Test API connections** (API Endpoints pane)
3. **Create ingestion jobs** (Ingestion Jobs pane)
4. **Schedule regular updates** (Scheduler pane)
5. **Monitor progress** (Dashboard pane)
6. **Get AI assistance** (AI Assistant pane)

## AI Assistant Setup

The AI Assistant is powered by OpenRouter and supports multiple models:

### 1. Get OpenRouter API Key

1. Visit [OpenRouter](https://openrouter.ai/)
2. Sign up for an account
3. Generate an API key
4. Add credits to your account (many models are very affordable)

### 2. Configure the AI Assistant

```bash
# Set your API key
dit config set ai.api_key "YOUR_OPENROUTER_KEY"

# Choose your preferred model
dit config set ai.model "anthropic/claude-3-haiku"  # Fast and affordable
# dit config set ai.model "anthropic/claude-3-sonnet"  # Balanced
# dit config set ai.model "openai/gpt-4-turbo"  # Most capable
```

### 3. Enable RAG (Retrieval-Augmented Generation)

RAG allows the AI to access your ingested data and configuration:

```bash
dit config set ai.rag.enabled true
dit config set ai.rag.chunk_size 1000
dit config set ai.rag.top_k 5
```

### 4. Available AI Features

- **Chat Interface**: Ask questions about your data and configuration
- **Code Generation**: Generate ingestion scripts and SQL queries
- **API Analysis**: Analyze API responses and schemas
- **Troubleshooting**: Get help with errors and issues
- **Best Practices**: Learn data management best practices
- **Automation**: Automate repetitive tasks

### 5. Example AI Interactions

```
You: "How do I ingest data from the Congress.gov API?"

AI: "I can help you set up Congress.gov data ingestion! Here's what you need to do:

1. First, make sure you have an API key from https://api.congress.gov/sign-up/
2. Configure the endpoint in the API Endpoints pane
3. I can generate a custom ingestion script for specific data types

What type of congressional data are you interested in? Bills, members, committees, or something else?"
```

## Scheduling Jobs

### Cron-style Scheduling

The application uses cron-style scheduling:

```
# Format: minute hour day month weekday
# Examples:
0 6 * * *     # Daily at 6 AM
0 8 * * 1     # Weekly on Monday at 8 AM
0 2 1 * *     # Monthly on 1st at 2 AM
*/15 * * * *  # Every 15 minutes
```

### Creating Scheduled Jobs

1. Go to the **Scheduler** pane (F5)
2. Press **N** to create a new job
3. Fill in the job details:
   - **Name**: Descriptive name for the job
   - **Endpoint**: API endpoint to call
   - **Schedule**: Cron expression
   - **Parameters**: API parameters
   - **Enabled**: Whether the job is active

### Pre-configured Schedules

The application comes with several pre-configured schedules in `configs/endpoints.yaml`:

- **Congress Bills**: Daily at 6 AM
- **FRED GDP Data**: Weekly on Monday at 8 AM
- **SEC Filings**: Daily at 8 PM
- **Census Population**: Monthly on 1st at 2 AM

## Troubleshooting

### Common Issues

#### 1. "Database locked" error

```bash
# Stop the application and check for other instances
pkill dit

# Remove lock file if it exists
rm ~/.dit/dit.db-wal ~/.dit/dit.db-shm

# Restart the application
dit
```

#### 2. API key not working

```bash
# Check if the key is set correctly
dit config get apis.endpoints.congress_gov.api_key

# Test the API connection
dit test-api congress_gov

# Check the logs for detailed error messages
dit --log-level debug
```

#### 3. Permission denied errors

```bash
# Check data directory permissions
ls -la ~/.dit/

# Fix permissions if needed
chmod 755 ~/.dit/
chmod 644 ~/.dit/*
```

#### 4. Network connectivity issues

```bash
# Test basic connectivity
curl -I https://api.congress.gov/v3

# Check proxy settings if behind corporate firewall
export HTTP_PROXY=http://proxy.company.com:8080
export HTTPS_PROXY=http://proxy.company.com:8080
```

### Debug Mode

Run with debug logging to see detailed information:

```bash
dit --verbose --log-level debug
```

### Log Files

Logs are stored in:
- **Console**: Real-time logs in the Logs pane
- **File**: `~/.dit/logs/dit.log` (if file logging is enabled)

### Getting Help

1. **In-app Help**: Press F1 or go to the Help pane
2. **AI Assistant**: Ask the AI for help with specific issues
3. **Documentation**: Check the `docs/` directory
4. **Issues**: Report bugs on GitHub

## Next Steps

### 1. Explore the APIs

Start with these beginner-friendly APIs:
- **Congress.gov**: Legislative data
- **FRED**: Economic indicators
- **Census**: Population and demographic data

### 2. Set Up Automation

- Create scheduled jobs for regular data updates
- Use the AI assistant to generate custom scripts
- Set up data quality monitoring

### 3. Advanced Features

- **Custom Endpoints**: Add your own API endpoints
- **Data Transformations**: Create custom data processing rules
- **Notifications**: Set up alerts for job failures
- **SSH Access**: Enable remote access for team collaboration

### 4. Integration

- **Export Data**: Export to CSV, JSON, or databases
- **Webhooks**: Integrate with other systems
- **API**: Use the built-in REST API for programmatic access

### 5. Community

- **Contribute**: Help improve the project
- **Share**: Share your API configurations and scripts
- **Learn**: Join discussions about data management best practices

## Configuration Examples

### Minimal Configuration

```yaml
# ~/.dit.yaml
data_dir: "~/.dit"
log_level: "info"

ai:
  api_key: "your-openrouter-key"
  model: "anthropic/claude-3-haiku"

apis:
  endpoints:
    congress_gov:
      api_key: "your-congress-key"
```

### Production Configuration

```yaml
# ~/.dit.yaml
app_name: "Data Ingestion TUI - Production"
data_dir: "/opt/dit/data"
log_level: "warn"

database:
  url: "/opt/dit/data/dit.db"
  max_open_conns: 50
  max_idle_conns: 10

ui:
  refresh_interval: "5s"
  page_size: 50

ai:
  provider: "openrouter"
  api_key: "your-openrouter-key"
  model: "anthropic/claude-3-sonnet"
  temperature: 0.3
  rag:
    enabled: true
    chunk_size: 1500
    top_k: 10

scheduler:
  enabled: true
  max_concurrent: 10
  persist_jobs: true

security:
  tls_enabled: true
  cert_file: "/opt/dit/certs/server.crt"
  key_file: "/opt/dit/certs/server.key"

server:
  host: "0.0.0.0"
  port: 2222
  max_clients: 50
```

---

🎉 **Congratulations!** You're now ready to start using the Data Ingestion TUI. 

For more detailed information, check out:
- [API Reference](API_REFERENCE.md)
- [Configuration Guide](CONFIGURATION.md)
- [Development Guide](DEVELOPMENT.md)
- [Deployment Guide](DEPLOYMENT.md)

Happy data ingesting! 📊✨
