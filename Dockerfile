# ===========================================
# issue2md 多阶段构建 Dockerfile
# ===========================================
# 生产级最佳实践：
# 1. 多阶段构建 - 最小化镜像体积
# 2. 依赖缓存 - 加速构建
# 3. 安全加固 - 非 root 用户运行
# 4. 精简镜像 - 仅包含必需文件
# ===========================================

# ===========================================
# Stage 1: Builder - 编译阶段
# ===========================================
FROM golang:1.24.9-alpine AS builder

# 设置工作目录
WORKDIR /build

# 安装构建依赖（如果需要 CGO）
# apk add --no-cache git ca-certificates gcc musl-dev

# ===========================================
# 依赖缓存优化 - 关键步骤
# ===========================================
# 先只复制 go.mod 和 go.sum，利用 Docker 层缓存
# 这样当依赖未变化时，不会重新下载依赖
COPY go.mod go.sum ./

# 下载依赖（这一层会被缓存，除非 go.mod/go.sum 变化）
RUN go mod download && \
    go mod verify

# ===========================================
# 复制源代码
# ===========================================
# 复制整个源代码（包括 vendor 如果存在）
COPY . .

# ===========================================
# 构建应用
# ===========================================
# 设置构建参数
ARG VERSION="dev"
ARG BUILD_DATE="unknown"
ARG GIT_COMMIT="unknown"

# 构建参数配置
# - CGO_ENABLED=0: 禁用 CGO，生成纯静态二进制文件
# - GOOS=linux: 目标操作系统
# - GOARCH=amd64: 目标架构（可根据需要改为 arm64 等）
# - -ldflags: 注入版本信息和减小二进制体积
#   -s: 去除符号表
#   -w: 去除 DWARF 调试信息
#   -X: 注入变量
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -v \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildDate=${BUILD_DATE} -X main.GitCommit=${GIT_COMMIT}" \
    -o issue2md \
    ./cmd/issue2md/main.go

# ===========================================
# Stage 2: Final - 运行阶段
# ===========================================
# 使用 alpine 镜像（极度精简，~5MB）
FROM alpine:3.19

# 安装运行时必需的 CA 证书（用于 HTTPS 请求）
RUN apk add --no-cache ca-certificates tzdata && \
    # 更新 CA 证书
    update-ca-certificates

# ===========================================
# 安全加固 - 创建非 root 用户
# ===========================================
# 创建一个专用用户来运行应用
# UID 65534 通常是 nobody 用户，但我们创建自己的用户更安全
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

# 设置工作目录
WORKDIR /app

# ===========================================
# 从 builder 阶段复制二进制文件
# ===========================================
COPY --from=builder --chown=appuser:appuser /build/issue2md /app/issue2md

# ===========================================
# 切换到非 root 用户
# ===========================================
USER appuser

# 设置默认环境变量
ENV PATH="/app:${PATH}"

# 健康检查（可选）
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD issue2md --version || exit 1

# ===========================================
# 入口点
# ===========================================
ENTRYPOINT ["/app/issue2md"]

# 默认命令（显示帮助信息）
CMD ["--help"]
