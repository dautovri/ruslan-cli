.PHONY: build install test clean release

BINARY_NAME=ruslan-cli
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build the binary
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) .

# Install to $GOPATH/bin
install:
	go install $(LDFLAGS) .

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run tests with race detector
test-race:
	go test -v -race ./...

# Run tests with coverage and race detector (CI mode)
test-ci:
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Check test coverage percentage
coverage:
	@go test -coverprofile=coverage.out ./... > /dev/null
	@go tool cover -func=coverage.out | grep total | awk '{print "Total Coverage: " $$3}'

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -f coverage.out

# Build for multiple platforms
build-all:
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64 .
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64 .

# Run locally
run:
	go run $(LDFLAGS) . $(ARGS)

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Lint with autofix
lint-fix:
	golangci-lint run --fix

# Verify code (format, lint, test)
verify: fmt lint test

# CI verification (used in GitHub Actions)
ci: deps fmt lint test-ci

# Update dependencies
deps:
	go mod download
	go mod tidy

# Generate completions
completions:
	mkdir -p completions
	go run . completion bash > completions/$(BINARY_NAME).bash
	go run . completion zsh > completions/$(BINARY_NAME).zsh
	go run . completion fish > completions/$(BINARY_NAME).fish

# Help
help:
	@echo "Available targets:"
	@echo "  build         - Build the binary"
	@echo "  install       - Install to GOPATH/bin"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage"
	@echo "  test-race     - Run tests with race detector"
	@echo "  test-ci       - Run tests in CI mode (race + coverage)"
	@echo "  bench         - Run benchmarks"
	@echo "  coverage      - Show test coverage percentage"
	@echo "  clean         - Clean build artifacts"
	@echo "  build-all     - Build for multiple platforms"
	@echo "  run           - Run locally (use ARGS=... for arguments)"
	@echo "  fmt           - Format code"
	@echo "  lint          - Lint code"
	@echo "  lint-fix      - Lint with autofix"
	@echo "  verify        - Verify code (fmt + lint + test)"
	@echo "  ci            - CI verification (deps + fmt + lint + test-ci)"
	@echo "  deps          - Update dependencies"
	@echo "  completions   - Generate shell completions"
