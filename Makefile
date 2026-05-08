.PHONY: all build test test-verbose vet fmt fmt-check lint tidy clean

PACKAGES := ./...
FMT_FILES := $(shell find . -name "*.go")
BINARY := terraform-provider-chatbotkit

all: vet build test

build:
	go build -o $(BINARY) .

test:
	go test -race -count=1 $(PACKAGES)

test-verbose:
	go test -race -count=1 -v $(PACKAGES)

vet:
	go vet $(PACKAGES)

fmt:
	gofmt -w $(FMT_FILES)

fmt-check:
	@test -z "$$(gofmt -l $(FMT_FILES))" || (echo "Files not formatted:"; gofmt -l $(FMT_FILES); exit 1)

lint: vet fmt-check
	@if command -v golangci-lint > /dev/null 2>&1; then \
		golangci-lint run $(PACKAGES); \
	else \
		echo "golangci-lint not installed — skipping (run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)"; \
	fi

tidy:
	go mod tidy

clean:
	go clean -testcache
	rm -f $(BINARY)
