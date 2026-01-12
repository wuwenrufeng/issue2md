# issue2md

> 将 GitHub Issue/PR/Discussion 转换为格式化的 Markdown 文档

## 项目简介

**issue2md** 是一个用 Go 语言编写的命令行工具，用于快速将 GitHub Issues、Pull Requests 和 Discussions 转换为格式良好的 Markdown 文档。非常适合用于知识管理、文档归档和技术博客编写。

项目遵循 Go 语言的"少即是多"哲学，仅使用标准库实现，无任何第三方依赖。

## 核心特性

- 支持三种 GitHub 资源类型：Issue、Pull Request、Discussion
- 完整保留讨论内容（标题、正文、所有评论）
- 自动生成 YAML Frontmatter 元数据
- 可选的 Reactions 统计（[emoji] [数量]）
- 可选的用户链接（`[@username](https://github.com/username)`）
- 灵活的输出方式（stdout 或文件）
- GitHub Emoji shortcode 自动转换为 Unicode emoji
- 通过环境变量安全传入认证信息

## 安装指南

### 从源码构建

```bash
# 克隆仓库
git clone https://github.com/wuwenrufeng/issue2md.git
cd issue2md

# 构建二进制文件
make build

# （可选）安装到 $GOPATH/bin
make install
```

### 下载预编译二进制

前往 [Releases](https://github.com/wuwenrufeng/issue2md/releases) 页面下载适合你操作系统的二进制文件。

### 使用 Docker

#### 快速开始

```bash
# 构建镜像
docker build -t issue2md:latest .

# 转换 Issue 并输出到终端
docker run --rm issue2md:latest https://github.com/golang/go/issues/1

# 转换 Issue 并保存到文件（挂载卷）
docker run --rm -v $(pwd)/output:/app/output \
  issue2md:latest https://github.com/golang/go/issues/1 /app/output/golang-issue-1.md
```

#### 使用 Docker Compose

```bash
# 启动服务并执行转换
docker-compose run issue2md https://github.com/golang/go/issues/1

# 转换并保存到文件
docker-compose run issue2md \
  https://github.com/golang/go/issues/1 /app/output/golang-issue-1.md

# 使用环境变量（私有仓库）
GITHUB_TOKEN=ghp_xxx docker-compose run issue2md \
  https://github.com/owner/private-repo/issues/1
```

#### Docker 镜像特点

- **多阶段构建** - 镜像体积 < 10MB
- **依赖缓存** - 快速构建
- **安全加固** - 非 root 用户运行
- **静态二进制** - 无需额外依赖
- **健康检查** - 自动监控服务状态

#### 镜像信息

| 项目 | 说明 |
|------|------|
| 基础镜像 | Alpine 3.19 |
| 镜像大小 | ~8MB |
| Go 版本 | 1.24.9 |
| 运行用户 | appuser (UID 1000) |

## 使用方法

### 基本用法

```bash
# 转换 Issue 并输出到终端
issue2md https://github.com/golang/go/issues/1

# 转换 Issue 并保存到文件
issue2md https://github.com/golang/go/issues/1 golang-issue-1.md

# 转换 Pull Request
issue2md https://github.com/golang/go/pull/2

# 转换 Discussion
issue2md https://github.com/github/community/discussions/12345
```

### 启用 Reactions 统计

```bash
issue2md -enable-reactions https://github.com/owner/repo/issues/123
```

输出示例：
```markdown
### @alice - 2025-01-04 11:15:00

我也遇到了这个问题，建议在 HTTP 客户端中间件中添加自动重试逻辑。

 3 1
```

### 启用用户链接

```bash
issue2md -enable-user-links https://github.com/owner/repo/issues/123
```

用户名将显示为可点击链接：
```markdown
**作者**: [@username](https://github.com/username)
```

### 组合多个选项

```bash
issue2md -enable-reactions -enable-user-links \
  https://github.com/owner/repo/issues/123 output.md
```

## 命令行参数

### 语法格式

```bash
issue2md [flags] <URL> [output_file]
```

### 参数说明

| 参数 | 类型 | 必需 | 说明 |
|------|------|------|------|
| `URL` | string | 是 | GitHub Issue/PR/Discussion 的完整 URL |
| `output_file` | string | 否 | 输出文件路径，不提供则输出到 stdout |

### 支持的 URL 格式

| 资源类型 | URL 模式 | 示例 |
|----------|----------|------|
| Issue | `https://github.com/{owner}/{repo}/issues/{number}` | `https://github.com/golang/go/issues/123` |
| PR | `https://github.com/{owner}/{repo}/pull/{number}` | `https://github.com/golang/go/pull/456` |
| Discussion | `https://github.com/{owner}/{repo}/discussions/{number}` | `https://github.com/github/community/discussions/789` |

### 命令行选项

| 选项 | 说明 |
|------|------|
| `-enable-reactions` | 显示 reactions 统计（如  3 1） |
| `-enable-user-links` | 用户名显示为可点击链接 |
| `-version` | 显示版本信息 |
| `-help` | 显示帮助信息 |

### 环境变量

| 变量 | 说明 |
|------|------|
| `GITHUB_TOKEN` | GitHub Personal Access Token（可选） |

#### API 限额说明

| 认证状态 | API 限额 |
|---------|---------|
| 无 Token | 60 次/小时 |
| 有 Token | 5000 次/小时 |

设置方法：
```bash
# 临时设置（当前会话）
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# 永久设置（添加到 ~/.bashrc 或 ~/.zshrc）
echo 'export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' >> ~/.bashrc
source ~/.bashrc
```

> **注意**：不要在命令行中直接使用 `--token` 参数，这会导致 Token 被记录到 Shell 历史中，存在安全风险。

## 输出格式

### 文件结构

```markdown
---
title: "{标题}"
url: "{原始URL}"
author: "{@username}"
created_at: "YYYY-MM-DD HH:MM:SS"
status: "{open|closed|merged}"
---

# {标题}

**作者**: @username
**创建时间**: 2025-01-04 10:30:00
**状态**: Open

{主楼正文内容}

## 评论

### @username1 - 2025-01-04 13:00:00

{评论内容}

### @username2 - 2025-01-04 14:30:00

{评论内容}
 3 1
```

### YAML Frontmatter 字段

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `title` | string | Issue/PR/Discussion 标题 | `"Fix authentication bug"` |
| `url` | string | 原始 GitHub URL | `"https://github.com/owner/repo/issues/123"` |
| `author` | string | 作者用户名（带 @ 前缀） | `"@johndoe"` |
| `created_at` | string | 创建时间（本地化格式） | `"2025-01-04 10:30:00"` |
| `status` | string | 当前状态 | `"open"` / `"closed"` / `"merged"` |

### 支持的 Reactions 类型

| Unicode Emoji | Shortcode |
|---------------|-----------|
|  | `:thumbsup:` / `:+1:` |
|  | `:thumbsdown:` / `:-1:` |
|  | `:smile:` / `:laugh:` |
|  | `:tada:` / `:hooray:` |
|  | `:confused:` |
|  | `:heart:` |
|  | `:rocket:` |
|  | `:eyes:` |

## 构建方法

### Makefile 目标

```bash
# 查看所有可用目标
make help

# 常用目标
make build          # 构建二进制文件到 bin/ 目录
make test           # 运行所有测试（含竞态检测）
make test-short     # 运行快速测试（跳过集成测试）
make clean          # 清理构建产物
make install        # 安装到 $GOPATH/bin
make fmt            # 格式化代码
make vet            # 运行 go vet 检查
make lint           # 运行 golangci-lint 静态代码检查
make lint-fix       # 运行 golangci-lint 并自动修复问题
make coverage       # 生成测试覆盖率报告
make check          # 运行代码质量检查（fmt + vet + lint）
```

### Docker 相关目标

```bash
make docker-build   # 构建 Docker 镜像
make docker-test    # 测试 Docker 镜像
make docker-run     # 使用 Docker 运行 issue2md（需要 URL 参数）
make docker-clean   # 清理 Docker 镜像

# 示例：使用 Docker 运行
make docker-run URL=https://github.com/golang/go/issues/1
```

### 交叉编译

```bash
# Linux amd64
GOOS=linux GOARCH=amd64 make build

# macOS amd64
GOOS=darwin GOARCH=amd64 make build

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 make build

# Windows amd64
GOOS=windows GOARCH=amd64 make build
```

### 使用自定义版本信息构建

```bash
# 使用 Makefile
make build VERSION=v1.0.0 BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# 或直接使用 go build
go build \
  -ldflags="-X main.Version=v1.0.0 -X main.BuildDate=2025-01-04T10:00:00Z" \
  -o issue2md \
  ./cmd/issue2md/main.go
```

## 常见问题

### Q: 为什么会出现 "API rate limit exceeded" 错误？

**A**: 未设置 `GITHUB_TOKEN` 时，GitHub API 限流为 60 次/小时。解决方法：
1. 设置 `GITHUB_TOKEN` 环境变量（提升至 5000 次/小时）
2. 等待限流重置（每小时重置一次）

### Q: 支持私有仓库吗？

**A**: 支持。设置 `GITHUB_TOKEN` 后即可访问你有权限的私有仓库。

### Q: 输出的 Markdown 可以直接使用吗？

**A**: 可以。输出的 Markdown 符合标准格式，可直接用于：
- 静态站点生成器（Hugo, Jekyll, Hexo）
- Wiki 系统
- 技术博客
- 文档归档

### Q: 为什么不包含 PR 的代码 diff？

**A**: issue2md 专注于**讨论内容**（Issue/PR 描述、Review 评论），代码变更通常不是归档的重点。如果需要 diff，建议直接使用 GitHub 的 Patch 功能。

### Q: 支持哪些 URL 格式？

**A**: 目前仅支持完整的 GitHub URL 格式。不支持简化格式（如 `owner/repo#123`）。

## 开发

项目遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

### 运行测试

```bash
# 运行所有测试
make test

# 运行快速测试
make test-short

# 生成覆盖率报告
make coverage
```

### 代码质量检查

```bash
# 完整代码质量检查
make check

# 或分别运行
make fmt
make vet
make lint
```

### 项目结构

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go          # 程序入口
├── internal/
│   ├── cli/                 # CLI 逻辑
│   ├── config/              # 配置加载
│   ├── converter/           # Markdown 转换器
│   ├── github/              # GitHub API 客户端
│   └── parser/              # URL 解析器
├── specs/                   # 功能规格文档
├── Makefile                 # 构建脚本
├── Dockerfile               # Docker 镜像定义
└── go.mod                   # Go 模块定义
```

## 许可证

[MIT License](LICENSE)

## 作者

wuwenrufeng

---

**Made with Go**
