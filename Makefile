.PHONY: help build test lint clean install-deps run analyze generate

# Variables
APP_NAME=tf-iamgen
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=build/bin
GO=go
GOFLAGS=-v

help: ## Show this help message
	@echo "$(APP_NAME) v$(VERSION)"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

build: ## Build the application
	$(GO) build $(GOFLAGS) -ldflags "-X main.Version=$(VERSION)" -o $(BUILD_DIR)/$(APP_NAME) .

install-deps: ## Install Go dependencies
	$(GO) mod tidy
	$(GO) mod download
	$(GO) mod verify

test: ## Run all tests
	$(GO) test -v -cover ./...

test-unit: ## Run unit tests only
	$(GO) test -v -cover -short ./tests/unit/...

lint: ## Run code quality checks
	$(GO) fmt ./...
	$(GO) vet ./...
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "⚠️  golangci-lint not installed, skipping"

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)
	$(GO) clean -cache -testcache

run: build ## Build and run the application
	$(BUILD_DIR)/$(APP_NAME)

analyze: build ## Run: analyze current directory
	$(BUILD_DIR)/$(APP_NAME) analyze .

generate: build ## Run: generate IAM policy
	$(BUILD_DIR)/$(APP_NAME) generate --output generated-policy.json

dev: ## Watch for changes and rebuild (requires entr)
	@command -v entr >/dev/null 2>&1 || (echo "⚠️  entr not installed. Install with: brew install entr"; exit 1)
	find . -name "*.go" | entr -r $(MAKE) build

fmt: ## Format all Go code
	$(GO) fmt ./...
	@echo "✨ Code formatted"

mod-tidy: ## Tidy Go modules
	$(GO) mod tidy
	@echo "✨ Modules tidied"

all: clean install-deps lint test build ## Run all targets

.DEFAULT_GOAL := help
