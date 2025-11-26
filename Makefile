BINARY_NAME=slimjson
CMD_PATH=./cmd/slimjson
BUILD_DIR=bin

# Colors
COLOR_RESET=\033[0m
COLOR_BOLD=\033[1m
COLOR_GREEN=\033[32m
COLOR_BLUE=\033[34m
COLOR_YELLOW=\033[33m
COLOR_CYAN=\033[36m

.PHONY: all build test lint clean docker-build docker-run podman-build podman-run

all: lint test build

build:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)🔨 Building...$(COLOR_RESET)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "$(COLOR_GREEN)✅ Build complete: $(BUILD_DIR)/$(BINARY_NAME)$(COLOR_RESET)"

test:
	@echo "$(COLOR_BOLD)$(COLOR_YELLOW)🧪 Running tests...$(COLOR_RESET)"
	@go test -v ./...
	@echo "$(COLOR_GREEN)✅ Tests passed$(COLOR_RESET)"

lint:
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)🔍 Linting...$(COLOR_RESET)"
	@golangci-lint run
	@echo "$(COLOR_GREEN)✅ Linting passed$(COLOR_RESET)"

clean:
	@echo "$(COLOR_BOLD)$(COLOR_YELLOW)🧹 Cleaning...$(COLOR_RESET)"
	@rm -rf $(BUILD_DIR)
	@echo "$(COLOR_GREEN)✅ Clean complete$(COLOR_RESET)"

docker-build:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)🐳 Building Docker image...$(COLOR_RESET)"
	@docker build -t slimjson:latest .
	@echo "$(COLOR_GREEN)✅ Docker image built$(COLOR_RESET)"

docker-run:
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)🐳 Running Docker container...$(COLOR_RESET)"
	@docker run -i --rm slimjson:latest

podman-build:
	@echo "$(COLOR_BOLD)$(COLOR_BLUE)📦 Building Podman image...$(COLOR_RESET)"
	@podman build -t slimjson:latest .
	@echo "$(COLOR_GREEN)✅ Podman image built$(COLOR_RESET)"

podman-run:
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)📦 Running Podman container...$(COLOR_RESET)"
	@podman run -i --rm slimjson:latest
