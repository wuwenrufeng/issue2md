# issue2md API 接口概览

本文档概述 `internal/` 下各包对外暴露的主要接口，作为后续开发的参考。

**版本**: 1.0
**状态**: Draft
**最后更新**: 2025-01-04

---

## 📦 包结构总览

```
internal/
├── parser/       # URL 解析与类型识别
├── github/       # GitHub API 交互
├── converter/    # 数据转换为 Markdown
├── config/       # 配置管理
└── cli/          # CLI 业务逻辑
```

---

## 1. internal/parser

**职责**: 解析 GitHub URL，识别资源类型（Issue/PR/Discussion）并提取关键信息

### 核心类型

```go
// ResourceType 表示 GitHub 资源的类型
type ResourceType int

const (
    ResourceTypeUnknown ResourceType = iota
    ResourceTypeIssue
    ResourceTypePullRequest
    ResourceTypeDiscussion
)

// Resource 表示解析后的 GitHub 资源
type Resource struct {
    Type       ResourceType
    Owner      string    // 仓库所有者
    Repo       string    // 仓库名称
    Number     int       // Issue/PR/Discussion 编号
    OriginalURL string   // 原始 URL
}

// 实现 Stringer 接口
func (rt ResourceType) String() string
```

### 核心接口

```go
// ParseURL 解析 GitHub URL 并返回 Resource
//
// 支持的 URL 格式：
//   - Issue: https://github.com/{owner}/{repo}/issues/{number}
//   - PR:    https://github.com/{owner}/{repo}/pull/{number}
//   - Discussion: https://github.com/{owner}/{repo}/discussions/{number}
//
// 返回错误：
//   - ErrInvalidURLFormat: URL 格式无效
//   - ErrUnsupportedResourceType: 不支持的资源类型
func ParseURL(url string) (*Resource, error)
```

### 错误定义

```go
var (
    ErrInvalidURLFormat       = errors.New("invalid GitHub URL format")
    ErrUnsupportedResourceType = errors.New("unsupported resource type")
)
```

---

## 2. internal/github

**职责**: 与 GitHub REST API v3 交互，获取 Issue/PR/Discussion 数据

### 核心数据结构

```go
// User 表示 GitHub 用户
type User struct {
    Login   string // 用户名
    HTMLURL string // 个人主页 URL
}

// Reaction 表示评论的 reaction
type Reaction struct {
    Content string // "+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"
    Count   int
}

// Comment 表示通用评论（适用于 Issue、PR、Discussion）
type Comment struct {
    ID        int64
    User      User
    CreatedAt time.Time
    Body      string
    Reactions []Reaction
    Deleted   bool // 标记是否已删除（从 API 返回的 null 字段判断）
}

// Issue 表示 GitHub Issue
type Issue struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed"
    Body      string
    Comments  []Comment
}

// PullRequest 表示 GitHub Pull Request
type PullRequest struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed", "merged"
    Body      string
    Comments  []Comment // 仅包含 Review 评论
}

// Discussion 表示 GitHub Discussion
type Discussion struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed"
    Body      string
    Comments  []Comment // 包含主楼评论和所有回复，按时间排序
}
```

### 核心接口

```go
// Client GitHub API 客户端
type Client struct {
    token       string
    baseURL     string
    httpClient  *http.Client
}

// NewClient 创建新的 GitHub API 客户端
// token: Personal Access Token（可为空，但受 API 限流）
func NewClient(token string) *Client

// FetchIssue 获取指定 Issue 及其评论
// owner: 仓库所有者
// repo: 仓库名称
// number: Issue 编号
//
// 返回错误：
//   - ErrResourceNotFound: Issue 不存在
//   - ErrAPIRateLimit: API 限流
//   - ErrNetwork: 网络错误
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)

// FetchPullRequest 获取指定 PR 及其 Review 评论
// 不包含 diff 和提交历史
//
// 返回错误：同 FetchIssue
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)

// FetchDiscussion 获取指定 Discussion 及其所有评论（含回复）
// 评论已按时间正序排列，回复平铺展示
//
// 返回错误：同 FetchIssue
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

### API 端点常量

```go
const (
    BaseURL          = "https://api.github.com"
    IssueEndpoint    = "/repos/%s/%s/issues/%d"
    CommentsEndpoint = "/repos/%s/%s/issues/%d/comments"
    PREndpoint       = "/repos/%s/%s/pulls/%d"
    PRCommentsEndpoint = "/repos/%s/%s/pulls/%d/comments"
    DiscussionEndpoint = "/repos/%s/%s/discussions/%d"
    DiscussionCommentsEndpoint = "/repos/%s/%s/discussions/%d/comments"
)
```

### 错误定义

```go
var (
    ErrResourceNotFound = errors.New("resource not found")
    ErrAPIRateLimit     = errors.New("API rate limit exceeded")
    ErrNetwork          = errors.New("network error")
)

// APIError 表示 GitHub API 返回的错误
type APIError struct {
    Message  string
    Status   int // HTTP 状态码
    Response string // 响应体
}

func (e *APIError) Error() string
```

---

## 3. internal/converter

**职责**: 将 GitHub 数据结构转换为格式化的 Markdown 文本

### 核心类型

```go
// Converter Markdown 转换器
type Converter struct {
    enableReactions    bool // 是否启用 reactions
    enableUserLinks    bool // 是否启用用户链接
}

// NewConverter 创建新的 Converter
// options: 可选配置（通过函数式选项模式）
func NewConverter(options ...Option) *Converter
```

### 配置选项

```go
// Option 配置选项类型
type Option func(*Converter)

// WithReactions 启用 reactions
func WithReactions(enable bool) Option

// WithUserLinks 启用用户链接
func WithUserLinks(enable bool) Option
```

### 核心接口

```go
// ConvertIssue 将 Issue 转换为 Markdown
// 返回格式：
//   ---
//   title: "..."
//   url: "..."
//   author: "@..."
//   created_at: "YYYY-MM-DD HH:MM:SS"
//   status: "..."
//   ---
//
//   # {标题}
//   ...
func (c *Converter) ConvertIssue(issue *github.Issue) (string, error)

// ConvertPullRequest 将 PR 转换为 Markdown
func (c *Converter) ConvertPullRequest(pr *github.PullRequest) (string, error)

// ConvertDiscussion 将 Discussion 转换为 Markdown
func (c *Converter) ConvertDiscussion(discussion *github.Discussion) (string, error)
```

### 内部辅助方法（不对外暴露）

```go
// formatYAMLFrontmatter 生成 YAML Frontmatter
func (c *Converter) formatYAMLFrontmatter(...) string

// formatUser 格式化用户名（根据 enableUserLinks）
func (c *Converter) formatUser(user github.User) string

// formatTimestamp 格式化时间戳
func (c *Converter) formatTimestamp(t time.Time) string

// formatReactions 格式化 reactions（根据 enableReactions）
func (c *Converter) formatReactions(reactions []github.Reaction) string

// convertEmojiShortcode 将 GitHub emoji shortcode 转换为 Unicode emoji
func (c *Converter) convertEmojiShortcode(body string) string
```

---

## 4. internal/config

**职责**: 管理配置（从命令行参数和环境变量加载）

### 核心类型

```go
// Config 应用配置
type Config struct {
    // 输入
    URL         string // GitHub URL

    // 输出
    OutputFile  string // 输出文件路径（空字符串表示 stdout）

    // 功能开关
    EnableReactions    bool
    EnableUserLinks    bool

    // 认证
    Token       string // GitHub Token（从环境变量读取）
}
```

### 核心接口

```go
// LoadFromFlags 从命令行参数加载配置
// argv: 命令行参数（通常是 os.Args[1:]）
// stdout, stderr: 输出流（用于 --help 和 --version）
//
// 返回：
//   - config: 加载的配置
//   - exitCode: 如果需要立即退出（如 --help），返回退出码；否则返回 -1
func LoadFromFlags(argv []string, stdout, stderr io.Writer) (config *Config, exitCode int)
```

### 内部实现细节

```go
// 使用标准库 flag 包
// Flag 定义：
//   -enable-reactions: 启用 reactions
//   -enable-user-links: 启用用户链接
//   --version: 输出版本号
//   --help: 输出帮助信息
//   -h (同 --help)

// 位置参数：
//   [output_file]: 输出文件路径
```

### 版本信息

```go
var (
    Version   = "dev" // 在构建时通过 -ldflags 注入
    BuildDate = "unknown"
)
```

---

## 5. internal/cli

**职责**: CLI 业务逻辑编排（协调整个流程）

### 核心接口

```go
// Run 运行 CLI 应用
// argv: 命令行参数
// stdout: 标准输出
// stderr: 标准错误
//
// 返回退出码：
//   - 0: 成功
//   - 1: 任何错误
//
// 流程：
//   1. 加载配置 (config.LoadFromFlags)
//   2. 解析 URL (parser.ParseURL)
//   3. 创建 GitHub 客户端 (github.NewClient)
//   4. 获取数据 (FetchIssue/FetchPullRequest/FetchDiscussion)
//   5. 转换为 Markdown (converter.Convert*)
//   6. 输出到文件或 stdout
func Run(argv []string, stdout, stderr io.Writer) int
```

### 内部实现步骤

```go
func Run(argv []string, stdout, stderr io.Writer) int {
    // 1. 加载配置
    cfg, exitCode := config.LoadFromFlags(argv, stdout, stderr)
    if exitCode != -1 {
        return exitCode
    }

    // 2. 解析 URL
    resource, err := parser.ParseURL(cfg.URL)
    if err != nil {
        fmt.Fprintf(stderr, "Error: %v\n", err)
        return 1
    }

    // 3. 创建 GitHub 客户端
    client := github.NewClient(cfg.Token)

    // 4. 获取数据（根据资源类型分发）
    var markdown string
    switch resource.Type {
    case parser.ResourceTypeIssue:
        issue, err := client.FetchIssue(resource.Owner, resource.Repo, resource.Number)
        if err != nil {
            fmt.Fprintf(stderr, "Error fetching issue: %v\n", err)
            return 1
        }
        conv := converter.NewConverter(
            converter.WithReactions(cfg.EnableReactions),
            converter.WithUserLinks(cfg.EnableUserLinks),
        )
        markdown, err = conv.ConvertIssue(issue)
        // ... 处理错误

    case parser.ResourceTypePullRequest:
        // ... 类似逻辑

    case parser.ResourceTypeDiscussion:
        // ... 类似逻辑

    default:
        fmt.Fprintf(stderr, "Error: unsupported resource type\n")
        return 1
    }

    // 5. 输出
    if cfg.OutputFile != "" {
        if err := os.WriteFile(cfg.OutputFile, []byte(markdown), 0644); err != nil {
            fmt.Fprintf(stderr, "Error writing file: %v\n", err)
            return 1
        }
    } else {
        fmt.Fprint(stdout, markdown)
    }

    return 0
}
```

---

## 6. cmd/issue2md

**职责**: 应用入口点（minimal main 函数）

```go
package main

import (
    "os"
    "github.com/issue2md/internal/cli"
)

func main() {
    os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
```

---

## 🔄 包依赖关系

```
cmd/issue2md (main)
  └── internal/cli (业务逻辑编排)
        ├── internal/config (配置管理)
        ├── internal/parser (URL 解析)
        ├── internal/github (API 客户端)
        └── internal/converter (Markdown 转换)
```

**依赖原则**：
- ✅ 单向依赖：上层依赖下层
- ✅ `github` 包是最底层（仅依赖标准库）
- ✅ `converter` 依赖 `github` 包的数据结构
- ✅ `cli` 协调所有包
- ❌ 下层不依赖上层

---

## 🧪 测试策略

### 单元测试

每个包的测试文件：

```
internal/parser/parser_test.go       # 表格驱动测试 URL 解析
internal/github/github_test.go       # 集成测试（使用真实 GitHub API）
internal/converter/converter_test.go # 表格驱动测试 Markdown 生成
internal/config/config_test.go       # 表格驱动测试配置加载
internal/cli/cli_test.go             # 端到端测试（使用 fake stdout/stderr）
```

### 测试数据

创建测试仓库或在测试中使用真实的 public Issues（如 go/go.go1issues）。

---

## 📝 后续优化（非 MVP）

1. **性能优化**
   - 添加并发获取评论
   - 实现缓存机制

2. **功能扩展**
   - 支持私有仓库
   - 支持批量处理
   - 支持自定义模板

3. **代码质量**
   - 添加 Context 支持（超时控制）
   - 添加请求重试逻辑
   - 改进错误信息（国际化）

---

**文档结束**
