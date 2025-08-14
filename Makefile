# vssh Makefile

# Variables
BINARY_NAME=vssh
VERSION?=dev
COMMIT?=$(shell git rev-parse --short HEAD)
DATE?=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Build directory
BUILD_DIR=dist

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Default target
.PHONY: all
all: clean test build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

# Build for all platforms
.PHONY: build-all
build-all: clean
	mkdir -p $(BUILD_DIR)
	
	# Linux builds
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-amd64 .
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-linux-arm64 .
	
	# macOS builds
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-darwin-arm64 .
	
	# Windows builds
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-windows-amd64.exe .
	GOOS=windows GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-windows-arm64.exe .

# Run tests
.PHONY: test
test:
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run linting
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w .

# Run security scan
.PHONY: security
security:
	gosec ./...
	govulncheck ./...

# Install development tools
.PHONY: install-tools
install-tools:
	$(GOGET) golang.org/x/tools/cmd/goimports@latest
	$(GOGET) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	$(GOGET) golang.org/x/vuln/cmd/govulncheck@latest

# Install the binary locally
.PHONY: install
install: build
	sudo mv $(BINARY_NAME) /usr/local/bin/

# Uninstall the binary
.PHONY: uninstall
uninstall:
	sudo rm -f /usr/local/bin/$(BINARY_NAME)

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Development build with debug info
.PHONY: build-dev
build-dev:
	$(GOBUILD) -gcflags="all=-N -l" $(LDFLAGS) -o $(BINARY_NAME)-debug .

# Create a release (requires VERSION to be set)
.PHONY: release
release:
	@if [ "$(VERSION)" = "dev" ]; then \
		echo "Error: VERSION must be set for release (e.g., make release VERSION=v1.0.0)"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Release $(VERSION) created. GitHub Actions will build and publish the release."

# Generate checksums for release binaries
.PHONY: checksums
checksums:
	@if [ ! -d "$(BUILD_DIR)" ]; then \
		echo "Error: Build directory not found. Run 'make build-all' first."; \
		exit 1; \
	fi
	cd $(BUILD_DIR) && sha256sum * > checksums.txt
	@echo "Checksums generated in $(BUILD_DIR)/checksums.txt"

# Verify checksums
.PHONY: verify-checksums
verify-checksums:
	@if [ ! -f "$(BUILD_DIR)/checksums.txt" ]; then \
		echo "Error: Checksums file not found. Run 'make checksums' first."; \
		exit 1; \
	fi
	cd $(BUILD_DIR) && sha256sum -c checksums.txt

# Docker build (optional)
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

# Show help
.PHONY: help
help:
	@echo "vssh Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  all              Clean, test, and build"
	@echo "  build            Build binary for current platform"
	@echo "  build-all        Build binaries for all platforms"
	@echo "  build-dev        Build with debug information"
	@echo "  test             Run tests"
	@echo "  test-coverage    Run tests with coverage report"
	@echo "  clean            Clean build artifacts"
	@echo "  deps             Download and tidy dependencies"
	@echo "  lint             Run linting"
	@echo "  fmt              Format code"
	@echo "  security         Run security scans"
	@echo "  install-tools    Install development tools"
	@echo "  install          Install binary locally"
	@echo "  uninstall        Uninstall binary"
	@echo "  run              Build and run the application"
	@echo "  release          Create a release (requires VERSION)"
	@echo "  checksums        Generate checksums for release binaries"
	@echo "  verify-checksums Verify checksums"
	@echo "  docker-build     Build Docker image"
	@echo "  help             Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION          Version to build (default: dev)"
	@echo "  COMMIT           Git commit hash (auto-detected)"
	@echo "  DATE             Build date (auto-generated)"
	@echo ""
	@echo "Examples:"
	@echo "  make build                    # Build for current platform"
	@echo "  make build-all VERSION=v1.0.0 # Build all platforms with version"
	@echo "  make release VERSION=v1.0.0   # Create a release"
	@echo "  make test                     # Run tests"
