# Data Ingestion TUI - Makefile
# 
# This Makefile provides convenient commands for building, testing, and
# developing the Data Ingestion TUI application.
# 
# Author: Codegen AI Assistant
# Created: 2024
# License: MIT

# Application information
APP_NAME := dit
VERSION := 1.0.0
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Build configuration
GO_VERSION := 1.21
BINARY_NAME := $(APP_NAME)
BINARY_PATH := ./bin/$(BINARY_NAME)
MAIN_PATH := ./cmd/main.go

# Directories
BUILD_DIR := ./bin
DIST_DIR := ./dist
DOCS_DIR := ./docs
SCRIPTS_DIR := ./scripts

# Go build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME) -X main.gitCommit=$(GIT_COMMIT)"
BUILD_FLAGS := -v $(LDFLAGS)

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
NC := \033[0m # No Color

# Default target
.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "$(CYAN)Data Ingestion TUI - Build System$(NC)"
	@echo "$(YELLOW)Version: $(VERSION)$(NC)"
	@echo ""
	@echo "$(GREEN)Available targets:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(BLUE)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
.PHONY: build
build: clean ## Build the application
	@echo "$(GREEN)Building $(APP_NAME) v$(VERSION)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build $(BUILD_FLAGS) -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "$(GREEN)✅ Build complete: $(BINARY_PATH)$(NC)"

.PHONY: build-all
build-all: clean ## Build for all platforms
	@echo "$(GREEN)Building $(APP_NAME) for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	
	# Linux AMD64
	@echo "$(BLUE)Building for Linux AMD64...$(NC)"
	@GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	
	# Linux ARM64
	@echo "$(BLUE)Building for Linux ARM64...$(NC)"
	@GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	
	# macOS AMD64
	@echo "$(BLUE)Building for macOS AMD64...$(NC)"
	@GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	
	# macOS ARM64
	@echo "$(BLUE)Building for macOS ARM64...$(NC)"
	@GOOS=darwin GOARCH=arm64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	
	# Windows AMD64
	@echo "$(BLUE)Building for Windows AMD64...$(NC)"
	@GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	
	@echo "$(GREEN)✅ Multi-platform build complete!$(NC)"
	@ls -la $(DIST_DIR)/

.PHONY: install
install: build ## Install the application to $GOPATH/bin
	@echo "$(GREEN)Installing $(APP_NAME)...$(NC)"
	@go install $(BUILD_FLAGS) $(MAIN_PATH)
	@echo "$(GREEN)✅ Installation complete!$(NC)"

# Development targets
.PHONY: dev
dev: ## Run the application in development mode
	@echo "$(GREEN)Starting $(APP_NAME) in development mode...$(NC)"
	@go run $(MAIN_PATH) --verbose --log-level debug

.PHONY: dev-server
dev-server: ## Run the SSH server in development mode
	@echo "$(GREEN)Starting $(APP_NAME) SSH server in development mode...$(NC)"
	@go run $(MAIN_PATH) server --verbose --log-level debug --port 2223

.PHONY: watch
watch: ## Watch for changes and rebuild automatically
	@echo "$(GREEN)Watching for changes...$(NC)"
	@which air > /dev/null || (echo "$(RED)Installing air...$(NC)" && go install github.com/cosmtrek/air@latest)
	@air -c .air.toml

# Testing targets
.PHONY: test
test: ## Run all tests
	@echo "$(GREEN)Running tests...$(NC)"
	@go test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✅ Tests complete!$(NC)"

.PHONY: test-coverage
test-coverage: test ## Run tests with coverage report
	@echo "$(GREEN)Generating coverage report...$(NC)"
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Coverage report generated: coverage.html$(NC)"

.PHONY: test-integration
test-integration: ## Run integration tests
	@echo "$(GREEN)Running integration tests...$(NC)"
	@go test -v -tags=integration ./...

.PHONY: benchmark
benchmark: ## Run benchmarks
	@echo "$(GREEN)Running benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

# Code quality targets
.PHONY: lint
lint: ## Run linters
	@echo "$(GREEN)Running linters...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)Installing golangci-lint...$(NC)" && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...
	@echo "$(GREEN)✅ Linting complete!$(NC)"

.PHONY: fmt
fmt: ## Format code
	@echo "$(GREEN)Formatting code...$(NC)"
	@go fmt ./...
	@goimports -w .
	@echo "$(GREEN)✅ Code formatted!$(NC)"

.PHONY: vet
vet: ## Run go vet
	@echo "$(GREEN)Running go vet...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Vet complete!$(NC)"

.PHONY: security
security: ## Run security checks
	@echo "$(GREEN)Running security checks...$(NC)"
	@which gosec > /dev/null || (echo "$(RED)Installing gosec...$(NC)" && go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest)
	@gosec ./...
	@echo "$(GREEN)✅ Security checks complete!$(NC)"

# Dependency management
.PHONY: deps
deps: ## Download and tidy dependencies
	@echo "$(GREEN)Managing dependencies...$(NC)"
	@go mod download
	@go mod tidy
	@go mod verify
	@echo "$(GREEN)✅ Dependencies updated!$(NC)"

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "$(GREEN)Updating dependencies...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✅ Dependencies updated!$(NC)"

.PHONY: deps-graph
deps-graph: ## Generate dependency graph
	@echo "$(GREEN)Generating dependency graph...$(NC)"
	@go mod graph | modgraphviz | dot -Tpng -o deps.png
	@echo "$(GREEN)✅ Dependency graph generated: deps.png$(NC)"

# Database targets
.PHONY: db-migrate
db-migrate: ## Run database migrations
	@echo "$(GREEN)Running database migrations...$(NC)"
	@go run $(MAIN_PATH) init
	@echo "$(GREEN)✅ Database migrations complete!$(NC)"

.PHONY: db-reset
db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "$(RED)⚠️  WARNING: This will delete all data!$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		echo "$(YELLOW)Resetting database...$(NC)"; \
		rm -f ~/.dit/dit.db; \
		$(MAKE) db-migrate; \
		echo "$(GREEN)✅ Database reset complete!$(NC)"; \
	else \
		echo "$(BLUE)Database reset cancelled.$(NC)"; \
	fi

# Documentation targets
.PHONY: docs
docs: ## Generate documentation
	@echo "$(GREEN)Generating documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)
	@which godoc > /dev/null || (echo "$(RED)Installing godoc...$(NC)" && go install golang.org/x/tools/cmd/godoc@latest)
	@godoc -http=:6060 &
	@echo "$(GREEN)✅ Documentation server started at http://localhost:6060$(NC)"

.PHONY: docs-generate
docs-generate: ## Generate static documentation
	@echo "$(GREEN)Generating static documentation...$(NC)"
	@mkdir -p $(DOCS_DIR)
	@go doc -all ./... > $(DOCS_DIR)/api.txt
	@echo "$(GREEN)✅ Documentation generated in $(DOCS_DIR)/$(NC)"

# Docker targets
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "$(GREEN)Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .
	@echo "$(GREEN)✅ Docker image built!$(NC)"

.PHONY: docker-run
docker-run: docker-build ## Run Docker container
	@echo "$(GREEN)Running Docker container...$(NC)"
	@docker run -it --rm -p 2222:2222 $(APP_NAME):latest

# Release targets
.PHONY: release
release: clean test lint build-all ## Create a release
	@echo "$(GREEN)Creating release v$(VERSION)...$(NC)"
	@mkdir -p $(DIST_DIR)
	@cd $(DIST_DIR) && \
	for binary in $(BINARY_NAME)-*; do \
		if [[ $$binary == *.exe ]]; then \
			zip $${binary%.exe}.zip $$binary; \
		else \
			tar -czf $$binary.tar.gz $$binary; \
		fi; \
	done
	@echo "$(GREEN)✅ Release v$(VERSION) created in $(DIST_DIR)/$(NC)"

.PHONY: release-notes
release-notes: ## Generate release notes
	@echo "$(GREEN)Generating release notes...$(NC)"
	@echo "# Release v$(VERSION)" > RELEASE_NOTES.md
	@echo "" >> RELEASE_NOTES.md
	@echo "**Build Information:**" >> RELEASE_NOTES.md
	@echo "- Version: $(VERSION)" >> RELEASE_NOTES.md
	@echo "- Build Time: $(BUILD_TIME)" >> RELEASE_NOTES.md
	@echo "- Git Commit: $(GIT_COMMIT)" >> RELEASE_NOTES.md
	@echo "" >> RELEASE_NOTES.md
	@echo "**Changes:**" >> RELEASE_NOTES.md
	@git log --oneline --since="1 week ago" >> RELEASE_NOTES.md || echo "- Initial release" >> RELEASE_NOTES.md
	@echo "$(GREEN)✅ Release notes generated: RELEASE_NOTES.md$(NC)"

# Cleanup targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -rf $(DIST_DIR)
	@rm -f coverage.out coverage.html
	@rm -f deps.png
	@echo "$(GREEN)✅ Cleanup complete!$(NC)"

.PHONY: clean-all
clean-all: clean ## Clean everything including dependencies
	@echo "$(GREEN)Cleaning everything...$(NC)"
	@go clean -modcache
	@rm -rf vendor/
	@echo "$(GREEN)✅ Deep cleanup complete!$(NC)"

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "$(CYAN)$(APP_NAME) v$(VERSION)$(NC)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"
	@echo "Go Version: $(shell go version)"

.PHONY: env
env: ## Show environment information
	@echo "$(CYAN)Environment Information:$(NC)"
	@echo "GOOS: $(shell go env GOOS)"
	@echo "GOARCH: $(shell go env GOARCH)"
	@echo "GOVERSION: $(shell go env GOVERSION)"
	@echo "GOPATH: $(shell go env GOPATH)"
	@echo "GOROOT: $(shell go env GOROOT)"
	@echo "PWD: $(shell pwd)"

.PHONY: check
check: deps vet lint test ## Run all checks (deps, vet, lint, test)
	@echo "$(GREEN)✅ All checks passed!$(NC)"

# Development setup
.PHONY: setup
setup: ## Set up development environment
	@echo "$(GREEN)Setting up development environment...$(NC)"
	@go mod download
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/godoc@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "$(GREEN)✅ Development environment setup complete!$(NC)"

# Air configuration for hot reloading
.air.toml:
	@echo "$(GREEN)Creating Air configuration...$(NC)"
	@cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
EOF
	@echo "$(GREEN)✅ Air configuration created!$(NC)"
