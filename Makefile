# Variables
APP_NAME = movie-database
BUILD_DIR = bin
GO_FILES = $(shell find . -type f -name '*.go')

# Default Target
.PHONY: all
all: run

# Build the application
.PHONY: build
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./main.go
	@echo "Build complete! Executable is at $(BUILD_DIR)/$(APP_NAME)"

# Run the application (used locally)
.PHONY: run
run: build
	@echo "Running $(APP_NAME)..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Cleanup complete!"

# Lint the code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Format the code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Vendor dependencies
.PHONY: vendor
vendor:
	@echo "Tidying and vendorizing dependencies..."
	@go mod tidy
	@go mod vendor
