# issue2md 技术实现方案

## 📋 文档信息

- **版本**: 1.0
- **创建日期**: 2025-01-04
- **状态**: Draft
- **作者**: 首席架构师
- **审查状态**: 待审查

---

## 1. 技术上下文总结

### 1.1 项目定位

issue2md 是一个**命令行工具**，用于将 GitHub Issue/PR/Discussion 的内容转换为格式化的 Markdown 文件，主要服务于个人开发者的知识归档需求。

### 1.2 技术栈选型

| 类别                  | 技术选型                | 版本/说明 | 选型理由                                    |
| --------------------- | ----------------------- | --------- | ------------------------------------------- |
| **编程语言**          | Go                      | >= 1.24.9 | 静态类型、高性能、跨平台编译                |
| **Web 框架**          | `net/http`              | 标准库    | 遵循"简单性原则"，不引入过度封装            |
| **GitHub API 客户端** | `google/go-github`      | v68+      | GitHub 官方库，处理认证、限流、REST/GraphQL |
| **Markdown 生成**     | 标准库 `fmt`, `strings` | -         | 逻辑简单，无需第三方库                      |
| **命令行参数**        | `flag`                  | 标准库    | Go 惯用的 CLI 参数解析                      |
| **数据存储**          | 无（实时 API 获取）     | -         | MVP 阶段无持久化需求                        |
| **测试框架**          | `testing`               | 标准库    | 表格驱动测试，TDD 流程                      |

### 1.3 架构风格

- **分层架构**：CLI 层 → 业务逻辑层 → 数据访问层
- **依赖方向**：单向依赖，上层依赖下层，下层不依赖上层
- **包设计原则**：高内聚、低耦合、单一职责

---

## 2. "合宪性"审查

### 2.1 对照 constitution.md 逐条审查

#### ✅ 第一条：简单性原则 (Simplicity First)

| 原则                                  | 审查结果 | 证据                                              |
| ------------------------------------- | -------- | ------------------------------------------------- |
| 1.1 YAGNI (只实现 spec.md 要求的功能) | ✅ 通过  | 技术方案仅覆盖 MVP 功能：Issue/PR/Discussion 转换 |
| 1.2 标准库优先                        | ✅ 通过  | Web 框架用`net/http`，Markdown 用标准库字符串操作 |
| 1.3 反过度工程                        | ✅ 通过  | 不引入 ORM、不引入消息队列、不实现复杂设计模式    |

**⚠️ 偏离说明**：

- 使用 `google/go-github` 而非纯手工 HTTP 调用
- **理由**：GitHub API 涉及认证（Bearer Token）、限流（X-RateLimit-Remaining）、分页、GraphQL 查询等复杂逻辑，手工实现容易出错且维护成本高。使用官方库符合"实用主义"，且该库本身轻量、仅依赖标准库。

#### ✅ 第二条：测试先行铁律 (Test-First Imperative)

| 原则                              | 审查结果 | 实施方案                                           |
| --------------------------------- | -------- | -------------------------------------------------- |
| 2.1 TDD 循环 (Red-Green-Refactor) | ✅ 通过  | 每个功能先写测试用例，再实现                       |
| 2.2 表格驱动测试                  | ✅ 通过  | `parser_test.go`, `converter_test.go` 使用表格驱动 |
| 2.3 拒绝 Mocks                    | ✅ 通过  | `github` 包使用真实 GitHub API（测试仓库）         |

**实施细节**：

- `internal/parser`：表格驱动测试 URL 解析逻辑
- `internal/github`：集成测试，使用真实的 public Issue（如 `golang/go#1`）
- `internal/converter`：表格驱动测试 Markdown 生成
- `internal/cli`：端到端测试，使用 fake stdout/stderr

#### ✅ 第三条：明确性原则 (Clarity and Explicitness)

| 原则                     | 审查结果 | 实施方案                                                   |
| ------------------------ | -------- | ---------------------------------------------------------- |
| 3.1 错误处理（显式处理） | ✅ 通过  | 所有错误必须检查，使用`fmt.Errorf("context: %w", err)`包装 |
| 3.2 无全局变量           | ✅ 通过  | 所有依赖通过函数参数或结构体成员注入                       |

**代码示例**：

```go
// ✅ 正确：错误包装
if err != nil {
    return fmt.Errorf("failed to fetch issue: %w", err)
}

// ❌ 错误：全局变量
var globalClient *github.Client  // 禁止！

// ✅ 正确：依赖注入
func FetchIssue(client *github.Client, owner, repo string, number int) (*Issue, error) {
    // ...
}
```

### 2.2 合规性总结

| 宪法原则   | 符合性  | 关键控制点                                |
| ---------- | ------- | ----------------------------------------- |
| 简单性原则 | ✅ 符合 | 仅实现 MVP 功能，使用轻量级依赖           |
| 测试先行   | ✅ 符合 | TDD 流程，表格驱动测试，真实 API 集成测试 |
| 明确性原则 | ✅ 符合 | 显式错误处理，依赖注入，无全局变量        |

**结论**：本技术方案**完全符合**项目宪法，仅有 1 个合理偏离（使用`google/go-github`），已记录理由。

---

## 3. 项目结构细化

### 3.1 目录树

```
issue2md/
├── cmd/
│   └── issue2md/
│       └── main.go                 # 入口点（minimal，仅调用 cli.Run）
│
├── internal/
│   ├── parser/                     # URL解析器
│   │   ├── parser.go
│   │   ├── types.go                # Resource, ResourceType
│   │   └── parser_test.go
│   │
│   ├── github/                     # GitHub API客户端封装
│   │   ├── client.go               # Client 结构体及方法
│   │   ├── types.go                # Issue, PR, Discussion, Comment 等数据结构
│   │   └── github_test.go          # 集成测试
│   │
│   ├── converter/                  # Markdown生成器
│   │   ├── converter.go            # Converter 结构体及方法
│   │   └── converter_test.go
│   │
│   ├── config/                     # 配置管理
│   │   ├── config.go               # Config 结构体
│   │   ├── loader.go               # LoadFromFlags 函数
│   │   └── config_test.go
│   │
│   └── cli/                        # CLI业务逻辑编排
│       ├── cli.go                  # Run 函数
│       ├── version.go              # Version 变量
│       └── cli_test.go
│
├── pkg/                            # （未来预留）可被外部使用的库
│
├── go.mod
├── go.sum
├── Makefile                        # 构建脚本
├── LICENSE
├── README.md
└── CLAUDE.md

```

### 3.2 包职责与依赖关系

#### 📦 internal/parser

**职责**：

- 解析 GitHub URL（Issue/PR/Discussion）
- 识别资源类型并提取关键信息（owner, repo, number）
- 验证 URL 格式

**对外接口**：

```go
func ParseURL(url string) (*Resource, error)
```

**依赖**：

- 仅依赖 Go 标准库（`net/url`, `regexp`, `strconv`）

**被依赖**：

- `internal/cli`

---

#### 📦 internal/github

**职责**：

- 封装 `google/go-github` 库
- 获取 Issue/PR/Discussion 数据及评论
- 处理 API 错误（404、限流等）
- 转换为内部数据结构（`Issue`, `PullRequest`, `Discussion`）

**对外接口**：

```go
type Client struct {
    client *github.Client
    ctx    context.Context
}

func NewClient(token string) *Client
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

**依赖**：

- `google/go-github`（外部依赖）
- Go 标准库（`context`, `time`）

**被依赖**：

- `internal/cli`
- `internal/converter`（使用其数据结构）

---

#### 📦 internal/converter

**职责**：

- 将 `github.Issue`/`PullRequest`/`Discussion` 转换为 Markdown
- 生成 YAML Frontmatter
- 格式化用户名、时间戳、reactions
- 转换 emoji shortcode

**对外接口**：

```go
type Converter struct {
    enableReactions bool
    enableUserLinks bool
}

func NewConverter(options ...Option) *Converter
func (c *Converter) ConvertIssue(issue *github.Issue) (string, error)
func (c *Converter) ConvertPullRequest(pr *github.PullRequest) (string, error)
func (c *Converter) ConvertDiscussion(discussion *github.Discussion) (string, error)
```

**依赖**：

- `internal/github`（使用其数据结构）
- Go 标准库（`fmt`, `strings`, `time`）

**被依赖**：

- `internal/cli`

---

#### 📦 internal/config

**职责**：

- 解析命令行参数（`flag` 包）
- 读取环境变量（`GITHUB_TOKEN`）
- 生成帮助和版本信息

**对外接口**：

```go
type Config struct {
    URL              string
    OutputFile       string
    EnableReactions  bool
    EnableUserLinks  bool
    Token            string
}

func LoadFromFlags(argv []string, stdout, stderr io.Writer) (*Config, int)
```

**依赖**：

- Go 标准库（`flag`, `os`, `io`）

**被依赖**：

- `internal/cli`

---

#### 📦 internal/cli

**职责**：

- 编排整个业务流程（配置加载 → URL 解析 → 数据获取 → Markdown 生成 → 输出）
- 错误处理和用户反馈
- 退出码管理

**对外接口**：

```go
func Run(argv []string, stdout, stderr io.Writer) int
```

**依赖**：

- `internal/config`
- `internal/parser`
- `internal/github`
- `internal/converter`
- Go 标准库（`fmt`, `os`）

**被依赖**：

- `cmd/issue2md`

---

### 3.3 依赖关系图

```
┌─────────────────────────────────────────┐
│           cmd/issue2md/main.go          │
│         (minimal entry point)           │
└─────────────────┬───────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────┐
│         internal/cli (编排层)           │
│  协调所有包，处理错误，管理退出码         │
└─────┬───────┬────────┬──────────┬──────┘
      │       │        │          │
      ▼       ▼        ▼          ▼
┌─────────┐ ┌──────┐ ┌──────┐ ┌──────────┐
│ config  │ │parser│ │github│ │converter │
└─────────┘ └──────┘ └──────┘ └──────────┘
                        │
                        ▼
                ┌─────────────────┐
                │ google/go-github│
                └─────────────────┘
```

**关键原则**：

- ✅ 单向依赖：无循环依赖
- ✅ 分层清晰：cli 在顶层，github 在底层
- ✅ 低耦合：包之间仅通过明确接口通信

---

## 4. 核心数据结构

### 4.1 internal/parser

```go
// ResourceType GitHub资源类型
type ResourceType int

const (
    ResourceTypeUnknown ResourceType = iota
    ResourceTypeIssue
    ResourceTypePullRequest
    ResourceTypeDiscussion
)

// Resource 解析后的GitHub资源
type Resource struct {
    Type        ResourceType
    Owner       string
    Repo        string
    Number      int
    OriginalURL string
}

// 实现 Stringer 接口
func (rt ResourceType) String() string {
    switch rt {
    case ResourceTypeIssue:
        return "issue"
    case ResourceTypePullRequest:
        return "pull_request"
    case ResourceTypeDiscussion:
        return "discussion"
    default:
        return "unknown"
    }
}
```

---

### 4.2 internal/github

```go
// User GitHub用户
type User struct {
    Login   string
    HTMLURL string
}

// Reaction 评论的reaction
type Reaction struct {
    Content string // "+1", "-1", "laugh", "hooray", "confused", "heart", "rocket", "eyes"
    Count   int
}

// Comment 通用评论（适用于Issue、PR、Discussion）
type Comment struct {
    ID        int64
    User      User
    CreatedAt time.Time
    Body      string
    Reactions []Reaction
    Deleted   bool // 标记是否已删除
}

// Issue GitHub Issue
type Issue struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed"
    Body      string
    Comments  []Comment
}

// PullRequest GitHub Pull Request
type PullRequest struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed", "merged"
    Body      string
    Comments  []Comment // 仅Review评论
}

// Discussion GitHub Discussion
type Discussion struct {
    Title     string
    URL       string
    User      User
    CreatedAt time.Time
    State     string // "open", "closed"
    Body      string
    Comments  []Comment // 包含主楼和所有回复，按时间排序
}
```

---

### 4.3 internal/config

```go
// Config 应用配置
type Config struct {
    // 输入
    URL string

    // 输出
    OutputFile string // 空字符串表示stdout

    // 功能开关
    EnableReactions bool
    EnableUserLinks bool

    // 认证
    Token string // 从环境变量GITHUB_TOKEN读取
}
```

---

### 4.4 internal/converter

```go
// Converter Markdown转换器
type Converter struct {
    enableReactions bool
    enableUserLinks bool
}

// Option 配置选项类型（函数式选项模式）
type Option func(*Converter)

// WithReactions 启用reactions
func WithReactions(enable bool) Option {
    return func(c *Converter) {
        c.enableReactions = enable
    }
}

// WithUserLinks 启用用户链接
func WithUserLinks(enable bool) Option {
    return func(c *Converter) {
        c.enableUserLinks = enable
    }
}
```

---

## 5. 接口设计

### 5.1 internal/parser

```go
// ParseURL 解析GitHub URL并返回Resource
//
// 支持的URL格式：
//   - Issue:      https://github.com/{owner}/{repo}/issues/{number}
//   - PR:         https://github.com/{owner}/{repo}/pull/{number}
//   - Discussion: https://github.com/{owner}/{repo}/discussions/{number}
//
// 返回错误：
//   - ErrInvalidURLFormat: URL格式无效
//   - ErrUnsupportedResourceType: 不支持的资源类型
func ParseURL(url string) (*Resource, error)
```

**错误定义**：

```go
var (
    ErrInvalidURLFormat        = errors.New("invalid GitHub URL format")
    ErrUnsupportedResourceType = errors.New("unsupported resource type")
)
```

---

### 5.2 internal/github

```go
// Client GitHub API客户端
type Client struct {
    client *github.Client
    ctx    context.Context
}

// NewClient 创建新的GitHub API客户端
// token: Personal Access Token（可为空，但受API限流）
func NewClient(token string) *Client

// FetchIssue 获取指定Issue及其评论
// owner: 仓库所有者
// repo: 仓库名称
// number: Issue编号
//
// 返回错误：
//   - ErrResourceNotFound: Issue不存在
//   - ErrAPIRateLimit: API限流
//   - ErrNetwork: 网络错误
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)

// FetchPullRequest 获取指定PR及其Review评论
// 不包含diff和提交历史
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)

// FetchDiscussion 获取指定Discussion及其所有评论（含回复）
// 评论已按时间正序排列，回复平铺展示
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
```

**错误定义**：

```go
var (
    ErrResourceNotFound = errors.New("resource not found")
    ErrAPIRateLimit     = errors.New("API rate limit exceeded")
    ErrNetwork          = errors.New("network error")
)

// APIError GitHub API返回的错误
type APIError struct {
    Message  string
    Status   int    // HTTP状态码
    Response string // 响应体
}

func (e *APIError) Error() string {
    return fmt.Sprintf("API error (status %d): %s", e.Status, e.Message)
}
```

---

### 5.3 internal/converter

```go
// Converter Markdown转换器
type Converter struct {
    enableReactions bool
    enableUserLinks bool
}

// NewConverter 创建新的Converter
// options: 可选配置（通过函数式选项模式）
func NewConverter(options ...Option) *Converter

// Option 配置选项类型
type Option func(*Converter)

// WithReactions 启用reactions
func WithReactions(enable bool) Option

// WithUserLinks 启用用户链接
func WithUserLinks(enable bool) Option

// ConvertIssue 将Issue转换为Markdown
func (c *Converter) ConvertIssue(issue *github.Issue) (string, error)

// ConvertPullRequest 将PR转换为Markdown
func (c *Converter) ConvertPullRequest(pr *github.PullRequest) (string, error)

// ConvertDiscussion 将Discussion转换为Markdown
func (c *Converter) ConvertDiscussion(discussion *github.Discussion) (string, error)
```

**内部辅助方法**（不对外暴露）：

```go
// formatYAMLFrontmatter 生成YAML Frontmatter
func (c *Converter) formatYAMLFrontmatter(title, url, author, createdAt, status string) string

// formatUser 格式化用户名（根据enableUserLinks）
func (c *Converter) formatUser(user github.User) string

// formatTimestamp 格式化时间戳（本地时区）
func (c *Converter) formatTimestamp(t time.Time) string

// formatReactions 格式化reactions（根据enableReactions）
func (c *Converter) formatReactions(reactions []github.Reaction) string

// convertEmojiShortcode 将GitHub emoji shortcode转换为Unicode emoji
func (c *Converter) convertEmojiShortcode(body string) string
```

---

### 5.4 internal/config

```go
// Config 应用配置
type Config struct {
    URL              string
    OutputFile       string
    EnableReactions  bool
    EnableUserLinks  bool
    Token            string
}

// LoadFromFlags 从命令行参数加载配置
// argv: 命令行参数（通常是os.Args[1:]）
// stdout, stderr: 输出流（用于--help和--version）
//
// 返回：
//   - config: 加载的配置
//   - exitCode: 如果需要立即退出（如--help），返回退出码；否则返回-1
func LoadFromFlags(argv []string, stdout, stderr io.Writer) (*Config, int)
```

**Flag 定义**（使用标准库`flag`包）：

```go
var (
    enableReactions   bool
    enableUserLinks   bool
    showVersion       bool
    showHelp          bool
)

func init() {
    flag.BoolVar(&enableReactions, "enable-reactions", false, "启用reactions统计")
    flag.BoolVar(&enableUserLinks, "enable-user-links", false, "启用用户链接")
    flag.BoolVar(&showVersion, "version", false, "输出版本号")
    flag.BoolVar(&showHelp, "help", false, "输出帮助信息")
    flag.BoolVar(&showHelp, "h", false, "输出帮助信息（ shorthand）")
}
```

**版本信息**：

```go
var (
    // Version 版本号（通过-ldflags注入）
    Version   = "dev"
    BuildDate = "unknown"
)
```

---

### 5.5 internal/cli

```go
// Run 运行CLI应用
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
//   2. 解析URL (parser.ParseURL)
//   3. 创建GitHub客户端 (github.NewClient)
//   4. 获取数据 (FetchIssue/FetchPullRequest/FetchDiscussion)
//   5. 转换为Markdown (converter.Convert*)
//   6. 输出到文件或stdout
func Run(argv []string, stdout, stderr io.Writer) int
```

**实现骨架**：

```go
func Run(argv []string, stdout, stderr io.Writer) int {
    // 1. 加载配置
    cfg, exitCode := config.LoadFromFlags(argv, stdout, stderr)
    if exitCode != -1 {
        return exitCode
    }

    // 2. 解析URL
    resource, err := parser.ParseURL(cfg.URL)
    if err != nil {
        fmt.Fprintf(stderr, "Error: %v\n", err)
        return 1
    }

    // 3. 创建GitHub客户端
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
        if err != nil {
            fmt.Fprintf(stderr, "Error converting: %v\n", err)
            return 1
        }

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

### 5.6 cmd/issue2md

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

## 6. 技术难点与解决方案

### 6.1 GitHub API 限流处理

**问题**：GitHub API 对未认证请求限制为 60 次/小时，认证请求为 5000 次/小时。

**解决方案**：

1. 优先从环境变量 `GITHUB_TOKEN` 读取 token
2. 在 API 响应中检查 `X-RateLimit-Remaining` 头
3. 遇到 403 错误时，透传 GitHub 的错误信息
4. 在文档中说明如何设置 token

**代码示例**：

```go
resp, err := c.client.Issues.Get(ctx, owner, repo, number)
if err != nil {
    if resp != nil && resp.StatusCode == 403 {
        return fmt.Errorf("API rate limit exceeded: %w", ErrAPIRateLimit)
    }
    return fmt.Errorf("failed to fetch issue: %w", err)
}
```

---

### 6.2 Discussion API 的 GraphQL 查询

**问题**：GitHub REST API v3 **不支持** Discussions 评论，需要使用 GraphQL API v4。

**解决方案**：

1. `google/go-github` 提供 GraphQL 支持
2. 使用 `github.Client.GraphQL` 方法执行查询
3. 编写 GraphQL 查询语句获取 Discussion 及其评论

**GraphQL 查询示例**：

```graphql
query GetDiscussion($owner: String!, $repo: String!, $number: Int!) {
  repository(owner: $owner, name: $repo) {
    discussion(number: $number) {
      title
      url
      author {
        login
      }
      createdAt
      state
      body
      comments(first: 100) {
        nodes {
          author {
            login
          }
          createdAt
          body
          reactions(first: 10) {
            nodes {
              content
              content
            }
          }
        }
      }
    }
  }
}
```

**Go 代码**：

```go
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error) {
    query := `
        query GetDiscussion($owner: String!, $repo: String!, $number: Int!) {
            repository(owner: $owner, name: $repo) {
                discussion(number: $number) {
                    title
                    url
                    author { login }
                    createdAt
                    state
                    body
                    comments(first: 100) {
                        nodes {
                            author { login }
                            createdAt
                            body
                            reactions(first: 10) {
                                nodes {
                                    content
                                }
                            }
                        }
                    }
                }
            }
        }
    `

    variables := map[string]interface{}{
        "owner": owner,
        "repo":  repo,
        "number": number,
    }

    var result struct {
        Repository struct {
            Discussion *Discussion
        }
    }

    err := c.client.GraphQL(c.ctx, query, variables, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch discussion: %w", err)
    }

    return result.Repository.Discussion, nil
}
```

---

### 6.3 Emoji Shortcode 转换

**问题**：GitHub 支持 `:thumbsup:` 这样的 emoji shortcode，需要转换为 Unicode emoji（👍）。

**解决方案**：

1. 维护一个 shortcode 到 emoji 的映射表
2. 使用正则表达式匹配并替换

**代码示例**：

```go
var emojiMap = map[string]string{
    ":thumbsup:":   "👍",
    ":+1:":         "👍",
    ":thumbsdown:": "👎",
    ":-1:":         "👎",
    ":heart:":      "❤️",
    ":smile:":      "😄",
    ":laugh:":      "😄",
    ":tada:":       "🎉",
    ":hooray:":     "🎉",
    ":confused:":   "😕",
    ":rocket:":     "🚀",
    ":eyes:":       "👀",
}

func convertEmojiShortcode(body string) string {
    for shortcode, emoji := range emojiMap {
        body = strings.ReplaceAll(body, shortcode, emoji)
    }
    return body
}
```

---

### 6.4 评论排序（平铺展示）

**问题**：Discussion 有嵌套回复结构，但要求平铺展示并按时间排序。

**解决方案**：

1. 递归获取所有评论（主楼 + 回复）
2. 展平为一维数组
3. 按 `CreatedAt` 字段排序

**代码示例**：

```go
func flattenComments(comments []Comment) []Comment {
    result := make([]Comment, 0, len(comments))

    for _, comment := range comments {
        result = append(result, comment)
        if len(comment.Replies) > 0 {
            result = append(result, flattenComments(comment.Replies)...)
        }
    }

    // 按时间排序
    sort.Slice(result, func(i, j int) bool {
        return result[i].CreatedAt.Before(result[j].CreatedAt)
    })

    return result
}
```

---

## 7. 测试策略

### 7.1 测试金字塔

```
        /\
       /  \      端到端测试 (internal/cli)
      /____\     - 完整流程测试
     /      \    - 使用fake stdout/stderr
    /________\
   /          \  集成测试 (internal/github)
  /____________\ - 真实GitHub API
 /              \
/________________\
 单元测试 (parser, converter, config)
- 表格驱动测试
- 快速、隔离
```

---

### 7.2 测试用例设计

#### internal/parser（表格驱动测试）

```go
func TestParseURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        want    *Resource
        wantErr error
    }{
        {
            name: "valid issue URL",
            url:  "https://github.com/owner/repo/issues/123",
            want: &Resource{
                Type:   ResourceTypeIssue,
                Owner:  "owner",
                Repo:   "repo",
                Number: 123,
            },
            wantErr: nil,
        },
        {
            name:    "invalid URL format",
            url:     "https://invalid-url.com/abc",
            want:    nil,
            wantErr: ErrInvalidURLFormat,
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseURL(tt.url)
            if err != tt.wantErr {
                t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("ParseURL() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

---

#### internal/github（集成测试）

```go
func TestClient_FetchIssue(t *testing.T) {
    // 跳过测试（如果环境变量GITHUB_TOKEN未设置）
    if os.Getenv("GITHUB_TOKEN") == "" {
        t.Skip("GITHUB_TOKEN not set, skipping integration test")
    }

    client := NewClient(os.Getenv("GITHUB_TOKEN"))

    // 使用真实的public issue（golang/go#1）
    issue, err := client.FetchIssue("golang", "go", 1)
    if err != nil {
        t.Fatalf("FetchIssue() error = %v", err)
    }

    // 验证关键字段
    if issue.Title == "" {
        t.Error("Issue title is empty")
    }
    if issue.State != "open" && issue.State != "closed" {
        t.Errorf("Issue state = %s, want 'open' or 'closed'", issue.State)
    }
}
```

---

#### internal/converter（表格驱动测试）

```go
func TestConverter_ConvertIssue(t *testing.T) {
    tests := []struct {
        name     string
        issue    *github.Issue
        options  []Option
        contains []string // 输出应该包含的字符串
    }{
        {
            name: "basic issue",
            issue: &github.Issue{
                Title:     "Test Issue",
                URL:       "https://github.com/owner/repo/issues/1",
                User:      github.User{Login: "testuser"},
                CreatedAt: time.Date(2025, 1, 4, 10, 0, 0, 0, time.UTC),
                State:     "open",
                Body:      "This is a test issue",
            },
            contains: []string{
                "title: \"Test Issue\"",
                "author: \"@testuser\"",
                "status: \"open\"",
                "# Test Issue",
            },
        },
        // ... 更多测试用例
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            conv := NewConverter(tt.options...)
            got, err := conv.ConvertIssue(tt.issue)
            if err != nil {
                t.Fatalf("ConvertIssue() error = %v", err)
            }

            for _, substr := range tt.contains {
                if !strings.Contains(got, substr) {
                    t.Errorf("ConvertIssue() output does not contain %q", substr)
                }
            }
        })
    }
}
```

---

#### internal/cli（端到端测试）

```go
func TestRun(t *testing.T) {
    tests := []struct {
        name     string
        argv     []string
        wantCode int
        wantOut  string
        wantErr  string
    }{
        {
            name:     "--help flag",
            argv:     []string{"--help"},
            wantCode: 0,
            wantOut:  "Usage:",
        },
        {
            name:     "--version flag",
            argv:     []string{"--version"},
            wantCode: 0,
            wantOut:  "version:",
        },
        {
            name:     "invalid URL",
            argv:     []string{"https://invalid-url.com"},
            wantCode: 1,
            wantErr:  "invalid GitHub URL format",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            stdout := &strings.Builder{}
            stderr := &strings.Builder{}

            gotCode := Run(tt.argv, stdout, stderr)

            if gotCode != tt.wantCode {
                t.Errorf("Run() code = %v, want %v", gotCode, tt.wantCode)
            }

            if tt.wantOut != "" && !strings.Contains(stdout.String(), tt.wantOut) {
                t.Errorf("Run() stdout = %v, want contain %q", stdout.String(), tt.wantOut)
            }

            if tt.wantErr != "" && !strings.Contains(stderr.String(), tt.wantErr) {
                t.Errorf("Run() stderr = %v, want contain %q", stderr.String(), tt.wantErr)
            }
        })
    }
}
```

---

### 7.3 测试覆盖率目标

| 包                   | 目标覆盖率 | 关键覆盖点                         |
| -------------------- | ---------- | ---------------------------------- |
| `internal/parser`    | 100%       | 所有 URL 格式、错误路径            |
| `internal/github`    | 80%+       | API 调用、错误处理（跳过网络错误） |
| `internal/converter` | 95%+       | 所有格式化逻辑、emoji 转换         |
| `internal/config`    | 90%+       | 所有 flag 组合                     |
| `internal/cli`       | 70%+       | 主要流程路径                       |

---

## 8. 部署与构建

### 8.1 构建脚本

```bash
#!/bin/bash
# Makefile

.PHONY: build test clean

# 版本信息（通过ldflags注入）
VERSION?=$(shell git describe --tags --always --dirty)
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# 构建二进制文件
build:
	go build -ldflags \
		"-X github.com/issue2md/internal/cli.Version=$(VERSION) \
		-X github.com/issue2md/internal/cli.BuildDate=$(BUILD_DATE)" \
		-o bin/issue2md \
		cmd/issue2md/main.go

# 运行测试
test:
	go test -v -race -cover ./...

# 交叉编译
build-all:
	GOOS=linux GOARCH=amd64 go build -o bin/issue2md-linux-amd64 cmd/issue2md/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/issue2md-darwin-amd64 cmd/issue2md/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/issue2md-windows-amd64.exe cmd/issue2md/main.go

# 清理
clean:
	rm -rf bin/
```

---

### 8.2 Go 版本管理

```go
// go.mod
module github.com/issue2md

go 1.24.9

require (
    github.com/google/go-github v68.0.0
    github.com/google/go-querystring v1.1.0 // indirect
)
```

---

## 9. 性能考量

### 9.1 性能指标

| 操作             | 预期性能  | 优化措施              |
| ---------------- | --------- | --------------------- |
| URL 解析         | < 1ms     | 正则表达式预编译      |
| API 调用（单次） | 100-500ms | 取决于网络延迟        |
| Markdown 生成    | < 10ms    | 使用`strings.Builder` |
| 总体端到端       | < 1s      | 无并发需求（MVP）     |

---

### 9.2 内存管理

- **Comment 缓存**：使用固定容量数组（预分配 100 条评论）
- **字符串拼接**：使用`strings.Builder`而非`+`操作符
- **API 响应**：及时释放`google/go-github`的响应体

---

## 10. 安全考量

### 10.1 Token 管理

- ✅ 仅从环境变量读取 `GITHUB_TOKEN`
- ✅ 不在日志中输出 token
- ✅ 不在错误信息中泄露 token
- ❌ 不提供 `--token` 参数（防止 Shell 历史泄露）

---

### 10.2 输入验证

- **URL 验证**：严格匹配 GitHub URL 格式
- **文件路径验证**：检查输出文件是否可写
- **参数边界检查**：Issue/PR 编号 > 0

---

## 11. 未来扩展（Post-MVP）

### 11.1 功能扩展

| 优先级 | 功能         | 技术方案                       |
| ------ | ------------ | ------------------------------ |
| P1     | 支持私有仓库 | 增强认证逻辑                   |
| P2     | 批量处理     | 循环处理多个 URL               |
| P3     | 本地缓存     | 使用`os.TempDir()` + JSON 存储 |
| P4     | 自定义模板   | 使用 Go `text/template`        |

---

### 11.2 Web 版本架构

```
┌─────────────┐
│  Web UI     │  (HTML/CSS/JS)
└──────┬──────┘
       │ HTTP
       ▼
┌─────────────┐
│  Web Server │  (net/http标准库)
└──────┬──────┘
       │
       ▼
┌─────────────┐
│ internal/*  │  (复用现有包)
└─────────────┘
```

**关键点**：

- Web 层独立为 `cmd/issue2mdweb/`
- 复用 `internal/github` 和 `internal/converter`
- 使用标准库 `net/http`，不引入 Gin/Echo

---

## 12. 实施路线图

### Phase 1: 基础设施（Week 1）

- [ ] 项目结构搭建
- [ ] `internal/parser` 实现 + 测试
- [ ] `internal/config` 实现 + 测试

### Phase 2: 数据获取（Week 2）

- [ ] `internal/github` 实现（Issue/PR）
- [ ] `internal/github` 实现（Discussion + GraphQL）
- [ ] 集成测试

### Phase 3: Markdown 生成（Week 3）

- [ ] `internal/converter` 实现
- [ ] Emoji 转换
- [ ] 格式化辅助方法
- [ ] 单元测试

### Phase 4: CLI 集成（Week 4）

- [ ] `internal/cli` 实现
- [ ] `cmd/issue2md/main.go`
- [ ] 端到端测试
- [ ] 文档完善

---

## 13. 总结

### 13.1 技术方案亮点

✅ **简单性**：仅使用 2 个外部依赖（`google/go-github` + 间接依赖）
✅ **可测试性**：清晰的接口设计，支持单元测试、集成测试、端到端测试
✅ **可扩展性**：分层架构，为 Web 版本预留扩展点
✅ **合规性**：100%符合项目宪法

---

### 13.2 风险评估

| 风险                      | 影响 | 缓解措施                         |
| ------------------------- | ---- | -------------------------------- |
| GitHub GraphQL API 复杂性 | 中   | 使用官方库，参考官方文档         |
| 测试用例维护成本          | 低   | 使用真实 public Issue，减少 mock |
| 性能瓶颈（API 限流）      | 中   | 文档说明，推荐设置 token         |

---

### 13.3 下一步行动

1. ✅ 技术方案审查（与团队确认）
2. ✅ 开始 Phase 1 实施（`internal/parser`）
3. ✅ 设置 CI/CD（GitHub Actions，自动运行测试）
4. ✅ 编写 README.md（使用说明）

---

**文档结束**
