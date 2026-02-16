APP_NAME := gobackend
BUILD_DIR := ./bin
MAIN_PKG := ./cmd/api
WIRE_PKG := github.com/goforj/wire/cmd/wire@latest
SQLC_PKG := github.com/sqlc-dev/sqlc/cmd/sqlc@latest

WIRE_DIRS := $(shell find . -name "wire.go" -not -path "./vendor/*" -exec dirname {} \;)

.PHONY: all
all: gen wire build

.PHONY: gen
gen:
	@echo "Generating sqlc code..."
	@sqlc generate

.PHONY: wire
wire:
	@echo "Generating wire dependencies..."
	@for dir in $(WIRE_DIRS); do \
		echo "  → $$dir"; \
		cd $$dir && go run $(WIRE_PKG) gen . && cd - > /dev/null; \
	done
	@echo "Wire generation complete!"

.PHONY: wire-dir
wire-dir:
	@if [ -z "$(DIR)" ]; then \
		echo "Usage: make wire-dir DIR=./path/to/wire"; \
		exit 1; \
	fi
	@echo "Generating wire for $(DIR)..."
	@cd $(DIR) && go run $(WIRE_PKG) gen .

.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PKG)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: dev
dev:
	@echo "Starting development server..."
	@GIN_MODE="debug" air

.PHONY: start
start:
	@echo "Starting production server..."
	@GIN_MODE="release" $(BUILD_DIR)/$(APP_NAME)

.PHONY: migrate
migrate:
	@go run ./cmd/migration/main.go $(ARGS)

.PHONY: install
install:
	@echo "Installing dependencies..."
	@go mod download
	@go install $(SQLC_PKG)
	@echo "Dependencies installed!"

.PHONY: tidy
tidy:
	@echo "Tidying modules..."
	@go mod tidy
	@go mod verify

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

.PHONY: clean-all
clean-all: clean
	@echo "Removing generated files..."
	@find . -name "wire_gen.go" -type f -delete
	@echo "All generated files removed!"

.PHONY: check
check: fmt lint test

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make gen           - Generate sqlc code"
	@echo "  make wire          - Generate wire dependencies for all directories"
	@echo "  make wire-dir DIR=./path - Generate wire for specific directory"
	@echo "  make build         - Build application binary"
	@echo "  make test          - Run all tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make dev           - Run development server with hot reload"
	@echo "  make start         - Run production binary"
	@echo "  make migrate       - Run database migrations"
	@echo "  make install       - Install dependencies"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make clean         - Clean build artifacts"
	@echo "  make clean-all     - Clean all including generated files"
	@echo "  make check         - Run fmt, lint, and test"
	@echo "  make all           - Run gen, wire, and build"
	@echo "  make help          - Show this help"
