BINARY_NAME := mcp-server
PKG := github.com/brduru/mcp-defect-dojo/cmd/mcp-server
BUILD_DIR := ./bin

# Version information
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS := -ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

.PHONY: all build test clean help examples version release

all: build

help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

version: ## Show version information
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

build: ## Build the MCP server binary
	@echo "Building $(BINARY_NAME)..."
	@go mod tidy
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(PKG)
	@echo "✅ Built $(BUILD_DIR)/$(BINARY_NAME)"

release: ## Build binaries for all platforms
	@echo "Building release binaries..."
	@mkdir -p $(BUILD_DIR)/release
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-amd64 $(PKG)
	@GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-linux-arm64 $(PKG)
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-amd64 $(PKG)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-darwin-arm64 $(PKG)
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/release/$(BINARY_NAME)-windows-amd64.exe $(PKG)
	@cd $(BUILD_DIR)/release && sha256sum * > checksums.txt
	@echo "✅ Release binaries built in $(BUILD_DIR)/release"

test: ## Run Go tests
	@echo "Running tests..."
	@go test ./...

test-verbose: ## Run Go tests with verbose output
	@echo "Running tests (verbose)..."
	@go test -v ./...

test-race: ## Run Go tests with race detector
	@echo "Running tests with race detector..."
	@go test -race ./...

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: ## Run golangci-lint (requires golangci-lint to be installed)
	@echo "Running golangci-lint..."
	@golangci-lint run

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "✅ Cleaned"

.DEFAULT_GOAL := help
