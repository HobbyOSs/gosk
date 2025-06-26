# Makefile for gosk

.PHONY: help all test gen build test-rerun-fails test-ci clean run dep fmt

# Set default goal to 'help'
.DEFAULT_GOAL := help

# --- Help ---
help: ## Show this help message
	@awk 'BEGIN {FS = ":.*##"; printf "Usage: make [target]\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

# Go parameters
GOBUILD=go build
BIN=gosk
NASK=wine nask.exe
KILL_DEAD_CODE = find . -type f -name "*.go" -exec sed -i -E '/^\s*\/\/.*(remove|delete|unnecessary|dead code|no longer needed)/Id' {} +

# Use go run to execute gotestsum, ensuring it uses the version from go.mod
GOTESTSUM = go run gotest.tools/gotestsum

all: build test ## Build and test the project

build: gen ## Build the project
	go build ./...
	$(GOBUILD) -v -o $(BIN) ./cmd/gosk

test: dep ## Run all tests
	$(GOTESTSUM) --format short-verbose -- ./...

test-rerun-fails: dep ## Rerun failed tests
	$(GOTESTSUM) --rerun-fails --format short-verbose --packages="./..."

test-ci: dep ## Run tests for CI
	$(GOTESTSUM) --junitfile report.xml --format dots -- ./...

clean: ## Clean up build artifacts
	go clean
	rm -f report.xml .gotestsum.json

run: build ## Build and run the application
	./$(BIN)

gen: dep ## Generate code from sources
	go generate ./...

fmt: ## Format the source code
	$(KILL_DEAD_CODE)
	go fmt ./...

dep: ## Install dependencies
	go mod download
	go mod tidy