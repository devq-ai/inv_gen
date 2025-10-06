# Invoice TUI - Complete Setup Guide

## Overview

This guide walks you through setting up the Invoice TUI from scratch, including all prerequisites and dependencies.

## Prerequisites

### 1. Go Installation

The TUI requires Go 1.21 or higher.

#### macOS

```bash
# Using Homebrew
brew install go

# Verify installation
go version
```

#### Linux (Ubuntu/Debian)

```bash
# Download and install Go 1.21
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin

# Reload shell
source ~/.bashrc  # or source ~/.zshrc

# Verify installation
go version
```

#### Linux (Other Distributions)

```bash
# Arch Linux
sudo pacman -S go

# Fedora
sudo dnf install golang

# Verify installation
go version
```

### 2. C Compiler (for SQLite CGo)

SQLite requires CGo, which needs a C compiler.

#### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install build-essential

# Fedora
sudo dnf install gcc

# Arch Linux
sudo pacman -S base-devel
```

### 3. Python 3.12+ (for PDF generation)

```bash
# Verify Python installation
python3 --version

# Should be 3.12 or higher
```

### 4. Invoice Database

```bash
# Navigate to parent directory
cd inv_gen

# Create database (if not exists)
python3 create_invoice_db.py

# Verify database created
ls -lh invoices.db
```

## Installation

### Step 1: Navigate to TUI Directory

```bash
cd inv_gen/tui
```

### Step 2: Install Go Dependencies

```bash
# Download all dependencies
go mod download

# Tidy up go.mod and go.sum
go mod tidy

# Verify dependencies
go mod verify
```

Expected output:
```
all modules verified
```

### Step 3: Build the Application

```bash
# Using Make (recommended)
make build

# Or manually
CGO_ENABLED=1 go build -o bin/invoice-tui .
```

Expected output:
```
🔨 Building invoice-tui...
✅ Build complete: bin/invoice-tui
```

### Step 4: Verify Build

```bash
# Check binary exists
ls -lh bin/invoice-tui

# Test run
./bin/invoice-tui --help
```

## Running the TUI

### Default Usage

```bash
# Run from tui directory (uses ../invoices.db)
./bin/invoice-tui
```

### Custom Database Path

```bash
# Specify database location
./bin/invoice-tui /path/to/invoices.db
```

### Using Make Commands

```bash
# Build and run
make run

# Run without building (development)
make run-dev

# Run with specific database
make run-with-db
```

## Verification Checklist

Run through this checklist to ensure everything is set up correctly:

```bash
# 1. Check Go installation
go version
# Expected: go version go1.21.x ...

# 2. Check CGo support
go env CGO_ENABLED
# Expected: 1

# 3. Check C compiler
gcc --version  # or clang --version on macOS
# Expected: version information

# 4. Check database
ls -lh ../invoices.db
# Expected: file exists, ~100KB

# 5. Test database query
sqlite3 ../invoices.db "SELECT COUNT(*) FROM invoices;"
# Expected: 27

# 6. Check Python
python3 --version
# Expected: Python 3.12 or higher

# 7. Test Python CLI
python3 ../invoice_cli.py list
# Expected: invoice list output

# 8. Build TUI
make build
# Expected: ✅ Build complete

# 9. Run TUI
./bin/invoice-tui
# Expected: TUI launches
```

## Troubleshooting

### Issue: "go: command not found"

**Solution:**
```bash
# Verify Go installation
which go

# If not found, ensure PATH is set
export PATH=$PATH:/usr/local/go/bin

# Add to ~/.bashrc or ~/.zshrc for persistence
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

### Issue: "CGo is required"

**Solution:**
```bash
# Enable CGo
export CGO_ENABLED=1

# Verify setting
go env CGO_ENABLED
# Should output: 1

# Install C compiler if needed
# macOS:
xcode-select --install

# Linux:
sudo apt-get install build-essential
```

### Issue: "cannot find package"

**Solution:**
```bash
# Clean module cache
go clean -modcache

# Reinstall dependencies
go mod download
go mod tidy

# Try building again
make build
```

### Issue: "Database file not found"

**Solution:**
```bash
# Check current directory
pwd
# Should be: .../inv_gen/tui

# Check parent directory for database
ls -la ../invoices.db

# If database doesn't exist, create it
cd ..
python3 create_invoice_db.py
cd tui

# Try running again
./bin/invoice-tui
```

### Issue: "invoice_cli.py not found"

**Solution:**
```bash
# Verify Python CLI exists
ls -la ../invoice_cli.py

# Make it executable (if needed)
chmod +x ../invoice_cli.py

# Test it directly
python3 ../invoice_cli.py list
```

### Issue: Build errors with mattn/go-sqlite3

**Solution:**
```bash
# Ensure CGo is enabled
export CGO_ENABLED=1

# Update go-sqlite3
go get -u github.com/mattn/go-sqlite3

# Clean and rebuild
make clean
make build
```

### Issue: "exec format error" on Linux

**Solution:**
```bash
# Rebuild for your architecture
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/invoice-tui .

# Or use make
make clean
make build
```

## Development Setup

For active development with live reload:

### Install Air (live reload)

```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with live reload
make dev
```

### Install golangci-lint (code quality)

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

### Development Workflow

```bash
# 1. Format code
make fmt

# 2. Run linter
make lint

# 3. Run tests
make test

# 4. Build and test
make run-dev
```

## System Installation

To install system-wide:

```bash
# Build optimized release
make release

# Install to /usr/local/bin (requires sudo)
make install

# Run from anywhere
invoice-tui
```

To uninstall:

```bash
make uninstall
```

## Docker Setup (Optional)

If you prefer Docker:

```bash
# Create Dockerfile (not included by default)
cat > Dockerfile << 'EOF'
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o invoice-tui .

FROM alpine:latest
RUN apk add --no-cache sqlite python3 py3-pip
WORKDIR /app
COPY --from=builder /app/invoice-tui .
COPY --from=builder /app/../invoice_cli.py ../
CMD ["./invoice-tui"]
EOF

# Build image
docker build -t invoice-tui .

# Run container
docker run -it --rm -v $(pwd)/../invoices.db:/app/invoices.db invoice-tui
```

## Performance Tuning

### Optimize Build

```bash
# Build with optimization flags
go build -ldflags="-s -w" -o bin/invoice-tui .

# Check binary size
ls -lh bin/invoice-tui
```

### Strip Debug Symbols

```bash
# Strip binary (reduces size)
strip bin/invoice-tui

# Verify size reduction
ls -lh bin/invoice-tui
```

## Environment Variables

Optional environment variables:

```bash
# Set database path
export INVOICE_DB_PATH=/path/to/invoices.db

# Enable debug mode
export INVOICE_TUI_DEBUG=1

# Set log file
export INVOICE_TUI_LOG=/tmp/invoice-tui.log
```

## Quick Start Script

Save this as `setup.sh`:

```bash
#!/bin/bash
set -e

echo "📦 Invoice TUI Setup Script"
echo ""

# Check Go
if ! command -v go &> /dev/null; then
    echo "❌ Go not installed. Please install Go 1.21+"
    exit 1
fi
echo "✅ Go found: $(go version)"

# Check CGo
if [ "$(go env CGO_ENABLED)" != "1" ]; then
    echo "⚠️  CGo not enabled. Enabling..."
    export CGO_ENABLED=1
fi
echo "✅ CGo enabled"

# Check database
if [ ! -f "../invoices.db" ]; then
    echo "❌ Database not found. Run: python3 ../create_invoice_db.py"
    exit 1
fi
echo "✅ Database found"

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Build
echo "🔨 Building..."
CGO_ENABLED=1 go build -o bin/invoice-tui .

echo ""
echo "✅ Setup complete!"
echo ""
echo "Run with: ./bin/invoice-tui"
```

Make it executable and run:

```bash
chmod +x setup.sh
./setup.sh
```

## Next Steps

After successful setup:

1. **Run the TUI**: `./bin/invoice-tui`
2. **Read the README**: Check `README.md` for usage details
3. **Try the features**: Navigate the dashboard, view invoices
4. **Generate PDFs**: Test PDF generation
5. **Explore**: Experiment with filtering and status updates

## Support

If you encounter issues:

1. Check the troubleshooting section above
2. Verify all prerequisites are installed
3. Review error messages carefully
4. Check the main README.md for additional help

For persistent issues:
- **Developer**: Dion Edge
- **Email**: dion@devq.ai
- **Company**: DevQ.ai

---

**Happy invoice managing! 🎨✨**