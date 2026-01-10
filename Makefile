# issue2md Makefile
# 将 GitHub Issue/PR/Discussion 转换为 Markdown

# 变量定义
BINARY_NAME=issue2md
BUILD_DIR=bin
CMD_DIR=cmd/issue2md
GO=go
GOFLAGS=-v

# 版本信息（默认值，构建时可覆盖）
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)"

# .PHONY 声明伪目标
.PHONY: all build test clean run help install fmt check vet

# 默认目标
all: build

## build: 构建二进制文件
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_DIR)/main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## test: 运行所有测试
test:
	@echo "Running tests..."
	$(GO) test -v -race -cover ./...

## test-short: 运行快速测试（跳过集成测试）
test-short:
	@echo "Running short tests..."
	$(GO) test -v -short ./...

## clean: 清理构建产物
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

## run: 构建并运行（需要提供 URL 参数）
run: build
	@echo "Running $(BINARY_NAME)..."
	@if [ -z "$(URL)" ]; then \
		echo "Error: Please provide URL parameter"; \
		echo "Usage: make run URL=https://github.com/owner/repo/issues/1"; \
		exit 1; \
	fi
	$(BUILD_DIR)/$(BINARY_NAME) $(URL)

## install: 安装到 $GOPATH/bin
install:
	@echo "Installing $(BINARY_NAME)..."
	$(GO) install $(LDFLAGS) $(CMD_DIR)/main.go
	@echo "Install complete"

## fmt: 格式化代码
fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

## vet: 运行 go vet 检查
vet:
	@echo "Running go vet..."
	$(GO) vet ./...

## check: 运行代码质量检查（fmt + vet）
check: fmt vet

## coverage: 生成测试覆盖率报告
coverage:
	@echo "Generating coverage report..."
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## help: 显示帮助信息
help:
	@echo "issue2md Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'
