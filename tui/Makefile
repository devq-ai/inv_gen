# Invoice TUI Makefile
# DevQ.ai Invoice Management System

.PHONY: help build run clean test install deps fmt lint dev docker

# Default target
.DEFAULT_GOAL := help

# Binary name
BINARY_NAME=invoice-tui
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -type f)

# Build flags
LDFLAGS=-ldflags="-s -w"
CGO_ENABLED=1

help: ## Show this help message
	@echo "Invoice TUI - Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

deps: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

build: deps ## Build the TUI application
	@echo "🔨 Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-dev: ## Build without optimization (faster, for development)
	@echo "🔨 Building $(BINARY_NAME) (dev mode)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Dev build complete"

run: build ## Build and run the TUI
	@echo "🚀 Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

run-dev: ## Run without building (using go run)
	@echo "🚀 Running in dev mode..."
	@CGO_ENABLED=$(CGO_ENABLED) go run .

run-with-db: build ## Run with custom database path
	@echo "🚀 Running with database: ../invoices.db"
	@$(BUILD_DIR)/$(BINARY_NAME) ../invoices.db

clean: ## Remove build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "✅ Clean complete"

fmt: ## Format Go code
	@echo "📝 Formatting code..."
	@go fmt ./...
	@echo "✅ Format complete"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

test: ## Run tests
	@echo "🧪 Running tests..."
	@CGO_ENABLED=$(CGO_ENABLED) go test -v ./...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	@CGO_ENABLED=$(CGO_ENABLED) go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

install: build ## Install to system
	@echo "📦 Installing $(BINARY_NAME)..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✅ Installed to /usr/local/bin/$(BINARY_NAME)"

uninstall: ## Uninstall from system
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "✅ Uninstalled"

dev: ## Run with live reload (requires air)
	@echo "🔥 Starting dev server with live reload..."
	@if command -v air >/dev/null 2>&1; then \
		air; \
	else \
		echo "⚠️  air not installed. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

check-db: ## Check if database exists
	@echo "🔍 Checking for database..."
	@if [ -f ../invoices.db ]; then \
		echo "✅ Database found: ../invoices.db"; \
		sqlite3 ../invoices.db "SELECT COUNT(*) || ' invoices in database' FROM invoices;"; \
	else \
		echo "❌ Database not found at ../invoices.db"; \
		echo ""; \
		echo "Create database with:"; \
		echo "  cd .. && python3 create_invoice_db.py"; \
	fi

setup: ## Full setup: deps, database check, build
	@echo "🚀 Setting up Invoice TUI..."
	@$(MAKE) deps
	@$(MAKE) check-db
	@$(MAKE) build
	@echo ""
	@echo "✅ Setup complete! Run with:"
	@echo "   make run"
	@echo ""

release: clean deps test ## Build optimized release binary
	@echo "📦 Building release version..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo "✅ Release build: $(BUILD_DIR)/$(BINARY_NAME)"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t invoice-tui:latest .
	@echo "✅ Docker image built"

docker-run: ## Run in Docker container
	@echo "🐳 Running in Docker..."
	@docker run -it --rm -v $(PWD)/../invoices.db:/app/invoices.db invoice-tui:latest

benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

vet: ## Run go vet
	@echo "🔍 Running go vet..."
	@go vet ./...
	@echo "✅ Vet complete"

mod-update: ## Update all dependencies
	@echo "⬆️  Updating dependencies..."
	@go get -u ./...
	@go mod tidy
	@echo "✅ Dependencies updated"

size: build ## Show binary size
	@echo "📊 Binary size:"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME) | awk '{print "  " $$5 " - " $$9}'
	@echo ""
	@file $(BUILD_DIR)/$(BINARY_NAME)

verify: fmt vet lint test ## Run all verification steps
	@echo ""
	@echo "✅ All verifications passed!"

watch: ## Watch files and rebuild on change
	@echo "👀 Watching for changes..."
	@while true; do \
		inotifywait -q -e modify $(GO_FILES); \
		make build-dev; \
	done

info: ## Show build information
	@echo "📋 Build Information:"
	@echo "  Binary name:  $(BINARY_NAME)"
	@echo "  Build dir:    $(BUILD_DIR)"
	@echo "  Go version:   $$(go version)"
	@echo "  CGO enabled:  $(CGO_ENABLED)"
	@echo "  Go files:     $$(echo $(GO_FILES) | wc -w)"
	@echo "  Dependencies: $$(go list -m all | wc -l)"

.PHONY: all
all: clean deps fmt vet test build ## Run full build pipeline
	@echo ""
	@echo "✅ Full build complete!"
