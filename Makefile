BINARY_NAME=slimjson
CMD_PATH=./cmd/slimjson
BUILD_DIR=bin

.PHONY: all build test lint clean

all: lint test build

build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

test:
	@echo "Running tests..."
	go test -v ./...

lint:
	@echo "Linting..."
	golangci-lint run

clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)
