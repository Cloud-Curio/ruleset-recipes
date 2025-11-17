# Data Ingestion TUI 🚀

A comprehensive Terminal User Interface (TUI) application built with Go and the Bubble Tea ecosystem for managing data ingestion from various government and public APIs.

## Features ✨

### 🏛️ Data Sources
- **Government APIs**: Congress.gov, GovInfo.gov, OpenStates.org
- **Economic Data**: FRED (Federal Reserve Economic Data), Census Data
- **Financial Data**: SEC Filings, Political Investment Trades
- **Crime & Security**: FBI Crime Data, Declassified Documents
- **Housing**: APC Housing Survey Data
- **Market Data**: Free Stock Price and Economic Data

### 🤖 AI Integration
- **OpenRouter SDK** integration with multiple AI models
- **RAG (Retrieval-Augmented Generation)** capabilities
- **MCP (Model Context Protocol)** server support
- **Chat Interface** with AI assistance
- Support for popular AI coding tools (Cline, Continue.dev, etc.)

### 🔧 Management Tools
- **API Key Management** with secure keyring storage
- **SQL Migration Management** with version control
- **Script Scheduling** with cron-like functionality
- **Endpoint Schema Management** (Postman-like functionality)
- **Reverse Engineering** tools for API discovery
- **Debug Interface** for troubleshooting ingestion scripts

### 🎨 User Interface
- Built with **Bubble Tea** framework
- Styled with **Lip Gloss** for beautiful terminal UI
- **Wish** for SSH server capabilities
- **Harmonica** for animations
- **Huh** for interactive forms

## Architecture 🏗️

```
data-ingestion-tui/
├── cmd/                    # CLI entry points
├── internal/
│   ├── app/               # Application core
│   ├── config/            # Configuration management
│   ├── database/          # Database operations
│   ├── api/               # API client implementations
│   ├── ai/                # AI integration and RAG
│   ├── scheduler/         # Job scheduling
│   ├── ui/                # TUI components
│   ├── models/            # Data models
│   └── utils/             # Utility functions
├── pkg/                   # Public packages
│   ├── keyring/           # Secure key management
│   ├── migration/         # Database migrations
│   └── ingestion/         # Data ingestion engines
├── scripts/               # Ingestion scripts
├── migrations/            # SQL migration files
├── configs/               # Configuration files
└── docs/                  # Documentation
```

## Quick Start 🚀

1. **Clone and Build**
   ```bash
   git clone <repository-url>
   cd data-ingestion-tui
   go mod tidy
   go build -o dit cmd/main.go
   ```

2. **Initialize Configuration**
   ```bash
   ./dit init
   ```

3. **Launch TUI**
   ```bash
   ./dit
   ```

## Configuration 📝

The application uses a hierarchical configuration system:
- Environment variables
- Configuration files (YAML/JSON)
- Command-line flags
- Interactive prompts

## API Keys 🔐

API keys are securely stored using the system keyring:
- **macOS**: Keychain
- **Linux**: Secret Service API
- **Windows**: Windows Credential Manager

## Development 🛠️

### Prerequisites
- Go 1.21+
- SQLite3
- Git

### Building from Source
```bash
make build
```

### Running Tests
```bash
make test
```

### Development Mode
```bash
make dev
```

## Contributing 🤝

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License 📄

MIT License - see [LICENSE](LICENSE) file for details.

## Support 💬

- 📖 Documentation: [docs/](docs/)
- 🐛 Issues: GitHub Issues
- 💡 Feature Requests: GitHub Discussions

---

Built with ❤️ using the [Charm](https://charm.sh/) ecosystem.
