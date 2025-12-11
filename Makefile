.PHONY: build install test clean fmt lint

# Default target
default: build

# Build the provider binary
build:
	go build -o terraform-provider-chatbotkit

# Install dependencies
install:
	go mod download
	go mod tidy

# Run tests
test:
	go test ./... -v -timeout=30m

# Clean build artifacts
clean:
	rm -f terraform-provider-chatbotkit
	rm -rf dist/

# Format code
fmt:
	go fmt ./...

# Run linter (requires golangci-lint)
lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

# Validate API synchronization
validate-api:
	go run tools/validate-api-sync.go

# Run all checks
check: fmt lint validate-api test

# Help
help:
	@echo "Available targets:"
	@echo "  build    - Build the provider binary"
	@echo "  install  - Install/update Go dependencies"
	@echo "  test     - Run tests"
	@echo "  clean    - Remove build artifacts"
	@echo "  fmt      - Format Go code"
	@echo "  lint     - Run linter (requires golangci-lint)"
	@echo "  check    - Run fmt, lint, and test"
