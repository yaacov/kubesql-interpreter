# kubesql Parser Makefile

.PHONY: build build-static sha256 test test-coverage clean fmt vet golangci-lint lint deps install-golangci-lint

# Build variables
BINARY_NAME := kubesql
BUILD_DIR := ./bin
CMD_DIR := ./cmd/kubesql
PKG_DIR := ./pkg/kubesql
COVERAGE_FILE := coverage.out
COVERAGE_HTML := coverage.html

# Default target
build: ## Build the kubesql command-line tool
	@echo "Building kubesql..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)

build-static: # Build a static binary for kubesql
	@echo "Building static kubesql binary..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w -extldflags=-static" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(CMD_DIR)
	@echo "Static binary built: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

sha256: build-static ## Generate SHA256 signature for the static binary
	@echo "Generating SHA256 signature..."
	@cd $(BUILD_DIR) && sha256sum $(BINARY_NAME)-linux-amd64 > $(BINARY_NAME)-linux-amd64.sha256
	@echo "SHA256 signature generated: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256"
	@cat $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64.sha256

test: ## Run all tests
	@echo "Running tests..."
	@go test -v $(PKG_DIR)/

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=$(COVERAGE_FILE) $(PKG_DIR)/
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)/
	@rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)

fmt: ## Format Go code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

lint: fmt vet golangci-lint ## Run linting tools

deps: ## Install development dependencies
	@echo "Installing development dependencies..."
	@go mod download
	@go mod tidy

golangci-lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	@$(shell go env GOPATH)/bin/golangci-lint run

install-golangci-lint: ## Install golangci-lint
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
