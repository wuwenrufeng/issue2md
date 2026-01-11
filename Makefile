# issue2md Makefile
# 将 GitHub Issue/PR/Discussion 转换为 Markdown

# 变量定义
BINARY_NAME=issue2md
BUILD_DIR=bin
CMD_DIR=cmd/issue2md
GO=go
GOFLAGS=-v

# Lint 工具配置
GOLANGCI_LINT=golangci-lint
GOLANGCI_LINT_VERSION=v1.60
LINT_TIMEOUT=5m

# 版本信息（默认值，构建时可覆盖）
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildDate=$(BUILD_DATE)"

# .PHONY 声明伪目标
.PHONY: all build test clean run help install fmt check vet lint lint-fix install-lint docker-build docker-test docker-run

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

## check: 运行代码质量检查（fmt + vet + lint）
check: fmt vet lint

## lint: 运行 golangci-lint 静态代码检查
lint:
	@echo "Running golangci-lint..."
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint is not installed. Run 'make install-lint' to install it."; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run --timeout=$(LINT_TIMEOUT) --verbose ./...

## lint-fix: 运行 golangci-lint 并自动修复问题
lint-fix:
	@echo "Running golangci-lint with fixes..."
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		echo "golangci-lint is not installed. Run 'make install-lint' to install it."; \
		exit 1; \
	fi
	$(GOLANGCI_LINT) run --timeout=$(LINT_TIMEOUT) --fix --verbose ./...

## install-lint: 安装 golangci-lint 工具
install-lint:
	@echo "Installing golangci-lint..."
	@if ! command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then \
		if [ "$$(uname -s)" = "Darwin" ]; then \
			brew install golangci/tap/golangci-lint; \
		else \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
		fi \
	else \
		echo "golangci-lint is already installed: $$(which $(GOLANGCI_LINT))"; \
	fi
	@$(GOLANGCI_LINT) version

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

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "Building Docker image..."
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown") \
		-t issue2md:$(VERSION) \
		-t issue2md:latest \
		.
	@echo "Docker image built successfully"

## docker-test: 测试 Docker 镜像
docker-test: docker-build
	@echo "Testing Docker image..."
	@docker run --rm issue2md:$(VERSION) --version
	@echo "Docker image test passed"

## docker-run: 使用 Docker 运行 issue2md
docker-run: docker-build
	@echo "Running issue2md in Docker..."
	@if [ -z "$(URL)" ]; then \
		echo "Error: Please provide URL parameter"; \
		echo "Usage: make docker-run URL=https://github.com/owner/repo/issues/1"; \
		exit 1; \
	fi
	docker run --rm -v $(PWD)/output:/app/output issue2md:$(VERSION) $(URL)

## docker-clean: 清理 Docker 镜像
docker-clean:
	@echo "Cleaning Docker images..."
	docker rmi issue2md:latest issue2md:$(VERSION) 2>/dev/null || true
	@echo "Docker images cleaned"
