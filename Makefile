BIN_NAME        := dockernotifier
BUILD_DIR       := ./bin
CMD_DIR         := ./cmd/docker-notifier
VERSION         := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS         := -ldflags "-X main.version=$(VERSION)"

GOOS            ?= $(shell go env GOOS)
GOARCH          ?= $(shell go env GOARCH)
GO_FILES        := $(shell find . -name "*.go" -type f)

# Docker
APP_NAME        := $(BIN_NAME)
DOCKER_IMAGE    := $(DOCKER_USERNAME)/$(APP_NAME)
DOCKER_TAG      := latest
DOCKER_PLATFORM := linux/amd64,linux/arm64

# GITHUB
GITHUB_USERNAME := andvarfolomeev

# GHCR
GHCR_IMAGE      := ghcr.io/$(GITHUB_USERNAME)/$(APP_NAME)

.PHONY: all build clean test lint deps install run \
        docker-build docker-run docker-publish docker-ghcr-publish

all: clean build docker-build

build:
	@echo "🔨 Building $(BIN_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BIN_NAME) $(CMD_DIR)
	@echo "✅ Build complete: $(BUILD_DIR)/$(BIN_NAME)"

clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

test:
	@echo "🧪 Running tests..."
	@go test -v ./...

lint:
	@echo "🔍 Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not installed. Run 'make deps' first."; \
		exit 1; \
	fi

deps:
	@echo "📦 Installing dependencies..."
	@go get -v ./...
	@go install @latest

install: build
	@echo "📥 Installing $(BIN_NAME)..."
	@cp $(BUILD_DIR)/$(BIN_NAME) $(GOPATH)/bin/

run: build
	@echo "🚀 Running $(BIN_NAME)..."
	@$(BUILD_DIR)/$(BIN_NAME) $(ARGS)

docker-build:
	@echo "🐳 Building Docker image..."
	@docker buildx build \
	  --platform $(DOCKER_PLATFORM) \
	  -t $(DOCKER_IMAGE):$(DOCKER_TAG) \
	  --load .
	@echo "✅ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-publish:
	@echo "📤 Publishing Docker image to Docker Hub..."
	@docker buildx build \
	  --platform $(DOCKER_PLATFORM) \
	  -t $(DOCKER_IMAGE):$(DOCKER_TAG) \
	  --push .

docker-push:
	@echo "📤 Publishing Docker image to GHCR..."
	@docker buildx build \
	  --platform $(DOCKER_PLATFORM) \
	  -t $(GHCR_IMAGE):$(DOCKER_TAG) \
	  --push .

docker-run:
	@echo "🐳 Running Docker container..."
	@docker run --rm -v /var/run/docker.sock:/var/run/docker.sock $(DOCKER_IMAGE):$(DOCKER_TAG) $(ARGS)
