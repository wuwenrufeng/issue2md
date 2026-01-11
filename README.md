# issue2md

> 将 GitHub Issue/PR/Discussion 转换为格式化的 Markdown 文档

## 简介

**issue2md** 是一个命令行工具，用于快速将 GitHub Issues、Pull Requests 和 Discussions 转换为格式良好的 Markdown 文档。非常适合用于知识管理、文档归档和技术博客编写。

## 功能特性

- ✅ 支持三种资源类型：Issue、Pull Request、Discussion
- ✅ 完整保留讨论内容（标题、正文、所有评论）
- ✅ 自动生成 YAML Frontmatter 元数据
- ✅ 可选的 Reactions 统计（👍❤️😄🎉）
- ✅ 可选的用户链接（`[@username](https://github.com/username)`）
- ✅ 灵活的输出方式（stdout 或文件）
- ✅ GitHub Emoji shortcode 自动转换为 Unicode emoji
- ✅ 通过环境变量安全传入认证信息

## 安装

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

### 使用 Docker（推荐）

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

#### 构建参数

```bash
# 自定义版本信息
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t issue2md:v1.0.0 .
```

#### Docker 镜像特点

- ✅ **多阶段构建** - 镜像体积 < 10MB
- ✅ **依赖缓存** - 快速构建
- ✅ **安全加固** - 非 root 用户运行
- ✅ **静态二进制** - 无需额外依赖
- ✅ **健康检查** - 自动监控服务状态

#### 镜像信息

- **基础镜像**: Alpine 3.19
- **镜像大小**: ~8MB
- **Go 版本**: 1.24.9
- **用户**: appuser (UID 1000)

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

👍 3 ❤️ 1
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

## 环境变量

### GITHUB_TOKEN（可选）

**推荐设置**以提高 API 请求限额：

| 认证状态 | API 限额 |
|---------|---------|
| 无 Token | 60 次/小时 |
| **有 Token** | **5000 次/小时** |

设置方法：

```bash
# 临时设置（当前会话）
export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

# 永久设置（添加到 ~/.bashrc 或 ~/.zshrc）
echo 'export GITHUB_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx' >> ~/.bashrc
source ~/.bashrc
```

> **注意**：不要在命令行中直接使用 `--token` 参数，这会导致 Token 被记录到 Shell 历史中，存在安全风险。

## 命令行参数

```
Usage:
  issue2md [flags] <URL> [output_file]

Arguments:
  URL          GitHub Issue/PR/Discussion 的完整 URL
  output_file  输出文件路径（可选，不提供则输出到 stdout）

Flags:
  -enable-reactions   显示 reactions 统计（如 👍 3 ❤️ 1）
  -enable-user-links  用户名显示为可点击链接
  -version            显示版本信息
  -help               显示此帮助信息

Environment Variables:
  GITHUB_TOKEN  GitHub Personal Access Token（可选）
```

## 输出格式示例

转换后的 Markdown 文件包含：

### YAML Frontmatter

```yaml
---
title: "Fix authentication bug in login flow"
url: "https://github.com/example/myproject/issues/42"
author: "@johndoe"
created_at: "2025-01-04 10:30:00"
status: "open"
---
```

### 正文内容

```markdown
# Fix authentication bug in login flow

**作者**: @johndoe
**创建时间**: 2025-01-04 10:30:00
**状态**: Open

## 描述

用户在使用 OAuth2 登录时，access_token 过期后没有正确刷新。

## 评论

### @alice - 2025-01-04 11:15:00

我也遇到了这个问题，建议在 HTTP 客户端中间件中添加自动重试逻辑。
```

## 构建说明

项目使用 Go 标准库开发，无需额外依赖。

### 构建目标

```bash
# 查看所有可用目标
make help

# 常用目标
make build      # 构建二进制文件
make test       # 运行所有测试
make clean      # 清理构建产物
make install    # 安装到 $GOPATH/bin
make fmt        # 格式化代码
make vet        # 运行 go vet 检查
```

### 交叉编译

```bash
# Linux
GOOS=linux GOARCH=amd64 make build

# macOS
GOOS=darwin GOARCH=amd64 make build

# Windows
GOOS=windows GOARCH=amd64 make build
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

## 开发

项目遵循 [Conventional Commits](https://www.conventionalcommits.org/) 规范。

```bash
# 运行测试
make test

# 代码质量检查
make check

# 生成覆盖率报告
make coverage
```

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 作者

wuwenrufeng

---

**Made with ❤️ by the Go community**
