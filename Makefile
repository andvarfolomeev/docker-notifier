BIN_NAME := dockernotifier
BUILD_DIR := ./bin
CMD_DIR := ./cmd/docker-notifier
GO_FILES := $(shell find . -name "*.go" -type f)
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: all build clean test lint deps install run docker-build docker-run

all: clean build docker-build

build:
	@echo "Building $(BIN_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN_NAME) $(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BIN_NAME)"

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Run 'make deps' first."; \
		exit 1; \
	fi

deps:
	@echo "Installing dependencies..."
	@go get -v ./...
	@go install @latest

install: build
	@echo "Installing $(BIN_NAME)..."
	@cp $(BUILD_DIR)/$(BIN_NAME) $(GOPATH)/bin/

run: build
	@echo "Running $(BIN_NAME)..."
	@$(BUILD_DIR)/$(BIN_NAME) $(ARGS)

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(BIN_NAME):$(VERSION) .
	@echo "Docker image built: $(BIN_NAME):$(VERSION)"

docker-run:
	@echo "Running Docker container..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock $(BIN_NAME):$(VERSION) $(ARGS)
