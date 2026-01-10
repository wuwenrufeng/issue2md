# issue2md 开发任务列表

---

## 📋 文档信息
- **版本**: 1.0
- **创建日期**: 2025-01-04
- **最后更新**: 2025-01-07
- **状态**: Active
- **基于文档**: spec.md, plan.md, constitution.md

---

## 📊 当前进度概览

| Phase | 名称 | 状态 | 完成度 |
|-------|------|------|--------|
| Phase 1 | Foundation（数据结构定义） | ✅ 完成 | 12/12 (100%) |
| Phase 2 | URL Parser（URL解析器） | ✅ 完成 | 10/10 (100%) |
| Phase 3 | Config Loader（配置加载器） | ✅ 完成 | 10/10 (100%) |
| Phase 4 | GitHub Fetcher（API客户端） | ✅ 完成 | 13/13 (100%) |
| Phase 5 | Markdown Converter（转换器） | ✅ 完成 | 15/15 (100%) |
| Phase 6 | CLI Assembly（命令行集成） | ✅ 完成 | 9/9 (100%) |
| Phase 7 | Main Entry Point（入口点） | ✅ 完成 | 1/1 (100%) |
| Phase 8 | Build & Documentation（构建和文档） | ✅ 完成 | 7/7 (100%) |
| Phase 9 | Code Review & Polish（代码审查和优化） | ❌ 未开始 | 0/4 (0%) |
| **总计** | | | **77/79 (97.5%)** |

### 下一步建议
🎯 **推荐优先级**：
1. ~~**Phase 4**（GitHub Fetcher）- 实现 API 客户端~~ ✅ 已完成
2. ~~**Phase 5**（Markdown Converter）- 完成转换逻辑~~ ✅ 已完成
3. ~~**Phase 3**（Config Loader）- 实现配置加载~~ ✅ 已完成
4. ~~**Phase 6**（CLI Assembly）- 集成所有模块~~ ✅ 已完成
5. ~~**Phase 7**（Main Entry Point）- 创建入口点~~ ✅ 已完成
6. ~~**Phase 8**（Build & Documentation）- 完善构建脚本和文档~~ ✅ 已完成
7. **Phase 9**（Code Review & Polish）- 代码质量检查和优化


## 📌 任务说明

### 符号标记
- `✅` : **已完成**
- `[P]` : **可并行**执行的任务（无依赖关系）
- **TDD强制**：测试任务必须在实现任务之前完成
- **依赖关系**：任务后面的括号表示依赖的前置任务

### 执行原则
1. **严格遵守TDD**：每个功能都必须先写测试（Red）→ 实现功能（Green）→ 重构（Refactor）
2. **小步提交**：每个任务完成后立即提交代码
3. **原子化**：每个任务只修改一个主要文件
4. **测试覆盖**：单元测试覆盖率目标：parser 100%, converter 95%+, github 80%+, config 90%+, cli 70%+

---

## Phase 1: Foundation（数据结构定义）

> **目标**：定义核心数据结构，搭建项目基础框架

### 1.1 项目初始化

- **任务 1.1.1** ✅ `[P]` 初始化Go模块
  - 创建 `go.mod` 文件
  - 设置Go版本要求 `go 1.24.9`
  - 添加模块路径 `module github.com/issue2md`

- **任务 1.1.2** ✅ `[P]` 创建基础目录结构
  - 创建 `cmd/issue2md/` 目录
  - 创建 `internal/{parser,github,converter,config,cli}/` 目录
  - 验证目录结构符合plan.md定义

- **任务 1.1.3** ✅ `[P]` 添加外部依赖
  - 执行 `go get github.com/google/go-github/v68`
  - 执行 `go mod tidy`
  - 验证 `go.sum` 文件生成

---

### 1.2 internal/parser - 数据结构定义

- **任务 1.2.1** ✅ `[P]` 创建 `internal/parser/types.go` - 定义ResourceType
  ```go
  - 定义 ResourceType 类型（int）
  - 定义4个常量：Unknown, Issue, PullRequest, Discussion
  - 实现 Stringer 接口（String() 方法）
  ```
  **文件**: `internal/parser/types.go`

- **任务 1.2.2** ✅ `[P]` 创建 `internal/parser/types.go` - 定义Resource结构体
  ```go
  - 定义 Resource 结构体
  - 字段：Type (ResourceType), Owner, Repo, Number (int), OriginalURL
  ```
  **文件**: `internal/parser/types.go`
  **依赖**: 任务 1.2.1

---

### 1.3 internal/github - 数据结构定义

- **任务 1.3.1** ✅ `[P]` 创建 `internal/github/types.go` - 定义基础类型
  ```go
  - 定义 User 结构体：Login, HTMLURL
  - 定义 Reaction 结构体：Content, Count
  ```
  **文件**: `internal/github/types.go`

- **任务 1.3.2** ✅ `[P]` 创建 `internal/github/types.go` - 定义Comment结构体
  ```go
  - 定义 Comment 结构体
  - 字段：ID (int64), User (User), CreatedAt (time.Time), Body, Reactions ([]Reaction), Deleted (bool)
  ```
  **文件**: `internal/github/types.go`
  **依赖**: 任务 1.3.1

- **任务 1.3.3** ✅ `[P]` 创建 `internal/github/types.go` - 定义Issue结构体
  ```go
  - 定义 Issue 结构体
  - 字段：Title, URL, User, CreatedAt, State, Body, Comments ([]Comment)
  ```
  **文件**: `internal/github/types.go`
  **依赖**: 任务 1.3.2

- **任务 1.3.4** ✅ `[P]` 创建 `internal/github/types.go` - 定义PullRequest结构体
  ```go
  - 定义 PullRequest 结构体
  - 字段：Title, URL, User, CreatedAt, State, Body, Comments ([]Comment)
  ```
  **文件**: `internal/github/types.go`
  **依赖**: 任务 1.3.2

- **任务 1.3.5** ✅ `[P]` 创建 `internal/github/types.go` - 定义Discussion结构体
  ```go
  - 定义 Discussion 结构体
  - 字段：Title, URL, User, CreatedAt, State, Body, Comments ([]Comment)
  ```
  **文件**: `internal/github/types.go`
  **依赖**: 任务 1.3.2

---

### 1.4 internal/config - 数据结构定义

- **任务 1.4.1** ✅ `[P]` 创建 `internal/config/config.go` - 定义Config结构体
  ```go
  - 定义 Config 结构体
  - 字段：URL, OutputFile, EnableReactions, EnableUserLinks, Token
  ```
  **文件**: `internal/config/config.go`

---

### 1.5 internal/converter - 数据结构定义

- **任务 1.5.1** ✅ `[P]` 创建 `internal/converter/converter.go` - 定义Converter结构体
  ```go
  - 定义 Converter 结构体
  - 字段：enableReactions, enableUserLinks (bool)
  - 定义 Option 类型：func(*Converter)
  ```
  **文件**: `internal/converter/converter.go`

- **任务 1.5.2** ✅ `[P]` 创建 `internal/converter/converter.go` - 定义选项函数
  ```go
  - 实现 NewConverter(options ...Option) *Converter
  - 实现 WithReactions(enable bool) Option
  - 实现 WithUserLinks(enable bool) Option
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 1.5.1

---

### 1.6 internal/cli - 数据结构定义

- **任务 1.6.1** ✅ `[P]` 创建 `internal/cli/version.go`
  ```go
  - 定义 Version 变量 = "dev"
  - 定义 BuildDate 变量 = "unknown"
  ```
  **文件**: `internal/cli/version.go`

---

## Phase 2: URL Parser（URL解析器，TDD）

> **目标**：实现GitHub URL解析功能，支持Issue/PR/Discussion

### 2.1 测试先行（Red Phase）

- **任务 2.1.1** ✅ 编写 `internal/parser/parser_test.go` - 表格驱动测试框架
  ```go
  - 创建测试文件
  - 定义测试用例结构（包含name, url, want, wantErr字段）
  - 创建空的测试函数 TestParseURL(t *testing.T)
  ```
  **文件**: `internal/parser/parser_test.go`

- **任务 2.1.2** ✅ 编写 `internal/parser/parser_test.go` - 有效URL测试用例
  ```go
  - 添加测试用例：valid issue URL
  - 添加测试用例：valid PR URL
  - 添加测试用例：valid Discussion URL
  - 每个用例验证返回的Resource字段正确
  ```
  **文件**: `internal/parser/parser_test.go`
  **依赖**: 任务 2.1.1

- **任务 2.1.3** ✅ 编写 `internal/parser/parser_test.go` - 无效URL测试用例
  ```go
  - 添加测试用例：invalid URL format
  - 添加测试用例：unsupported resource type
  - 添加测试用例：missing required fields
  - 验证返回正确的错误类型
  ```
  **文件**: `internal/parser/parser_test.go`
  **依赖**: 任务 2.1.1

- **任务 2.1.4** ✅ 编写 `internal/parser/parser_test.go` - 边界条件测试
  ```go
  - 添加测试用例：URL with query parameters（应忽略）
  - 添加测试用例：URL with fragment（应忽略）
  - 添加测试用例：case sensitivity
  ```
  **文件**: `internal/parser/parser_test.go`
  **依赖**: 任务 2.1.1

- **任务 2.1.5** ✅ 运行测试验证失败（Red）
  ```bash
  - 执行 go test ./internal/parser -v
  - 确认所有测试失败（因为功能未实现）
  ```
  **依赖**: 任务 2.1.2, 2.1.3, 2.1.4

---

### 2.2 实现功能（Green Phase）

- **任务 2.2.1** ✅ 创建 `internal/parser/parser.go` - 定义错误变量
  ```go
  - 定义 ErrInvalidURLFormat
  - 定义 ErrUnsupportedResourceType
  ```
  **文件**: `internal/parser/parser.go`

- **任务 2.2.2** ✅ 创建 `internal/parser/parser.go` - 实现ParseURL函数骨架
  ```go
  - 函数签名：func ParseURL(url string) (*Resource, error)
  - 实现基本的URL解析逻辑
  - 使用 net/url 包解析URL
  ```
  **文件**: `internal/parser/parser.go`
  **依赖**: 任务 2.2.1

- **任务 2.2.3** ✅ 创建 `internal/parser/parser.go` - 实现URL路径匹配逻辑
  ```go
  - 使用正则表达式匹配URL路径
  - 支持三种模式：/issues/{num}, /pull/{num}, /discussions/{num}
  - 提取 owner, repo, number
  ```
  **文件**: `internal/parser/parser.go`
  **依赖**: 任务 2.2.2

- **任务 2.2.4** ✅ 创建 `internal/parser/parser.go` - 实现资源类型识别
  ```go
  - 根据URL路径识别 ResourceType
  - 返回对应的 Resource 结构体
  ```
  **文件**: `internal/parser/parser.go`
  **依赖**: 任务 2.2.3

- **任务 2.2.5** ✅ 运行测试验证通过（Green）
  ```bash
  - 执行 go test ./internal/parser -v
  - 确认所有测试通过
  ```
  **依赖**: 任务 2.2.4

---

### 2.3 重构优化（Refactor Phase）

- **任务 2.3.1** ✅ `[P]` 重构 `internal/parser/parser.go` - 提取辅助函数优化代码结构
  ```go
  - 将正则表达式提取为包级常量
  - 使用 regexp.MustCompile 在 init 中编译
  ```
  **文件**: `internal/parser/parser.go`
  **依赖**: 任务 2.2.5

- **任务 2.3.2** ✅ `[P]` 重构 `internal/parser/parser.go` - 优化错误信息
  ```go
  - 使用 fmt.Errorf 包装错误
  - 提供更详细的错误上下文
  ```
  **文件**: `internal/parser/parser.go`
  **依赖**: 任务 2.3.1

- **任务 2.3.3** ✅ `[P]` 运行测试确保重构未破坏功能
  ```bash
  - 执行 go test ./internal/parser -v
  - 确认所有测试仍然通过
  ```
  **依赖**: 任务 2.3.1, 2.3.2

---

## Phase 3: Config Loader（配置加载器，TDD）

> **目标**：实现命令行参数解析和环境变量读取

### 3.1 测试先行（Red Phase）

- **任务 3.1.1** ✅ `[P]` 编写 `internal/config/loader_test.go` - 测试框架
  ```go
  - 创建测试文件
  - 定义测试辅助函数（创建fake stdout/stderr）
  - 创建空的测试函数 TestLoadFromFlags(t *testing.T)
  ```
  **文件**: `internal/config/loader_test.go`

- **任务 3.1.2** ✅ `[P]` 编写 `internal/config/loader_test.go` - 基本参数测试
  ```go
  - 添加测试用例：valid URL argument
  - 验证返回的Config.URL正确
  ```
  **文件**: `internal/config/loader_test.go`
  **依赖**: 任务 3.1.1

- **任务 3.1.3** ✅ `[P]` 编写 `internal/config/loader_test.go` - Flag测试
  ```go
  - 添加测试用例：--enable-reactions flag
  - 添加测试用例：--enable-user-links flag
  - 验证Config的布尔字段正确
  ```
  **文件**: `internal/config/loader_test.go`
  **依赖**: 任务 3.1.1

- **任务 3.1.4** ✅ `[P]` 编写 `internal/config/loader_test.go` - 环境变量测试
  ```go
  - 添加测试用例：GITHUB_TOKEN环境变量
  - 验证Config.Token正确读取
  ```
  **文件**: `internal/config/loader_test.go`
  **依赖**: 任务 3.1.1

- **任务 3.1.5** ✅ `[P]` 编写 `internal/config/loader_test.go` - 输出文件测试
  ```go
  - 添加测试用例：output file位置参数
  - 验证Config.OutputFile正确
  ```
  **文件**: `internal/config/loader_test.go`
  **依赖**: 任务 3.1.1

- **任务 3.1.6** ✅ `[P]` 编写 `internal/config/loader_test.go` - 帮助和版本测试
  ```go
  - 添加测试用例：--help flag（验证exitCode=0）
  - 添加测试用例：--version flag（验证exitCode=0）
  - 验证stdout输出包含预期内容
  ```
  **文件**: `internal/config/loader_test.go`
  **依赖**: 任务 3.1.1

- **任务 3.1.7** ✅ 运行测试验证失败（Red）
  ```bash
  - 执行 go test ./internal/config -v
  - 确认所有测试失败
  ```
  **依赖**: 任务 3.1.2, 3.1.3, 3.1.4, 3.1.5, 3.1.6

---

### 3.2 实现功能（Green Phase）

- **任务 3.2.1** ✅ 创建 `internal/config/loader.go` - 定义flag变量
  ```go
  - 定义包级flag变量：enableReactions, enableUserLinks, showVersion, showHelp
  - 在 init() 中使用 flag.BoolVar 注册
  ```
  **文件**: `internal/config/loader.go`

- **任务 3.2.2** ✅ 创建 `internal/config/loader.go` - 实现LoadFromFlags函数骨架
  ```go
  - 函数签名：func LoadFromFlags(argv []string, stdout, stderr io.Writer) (*Config, int)
  - 创建flag.FlagSet实例（避免污染全局flag）
  - 调用 flagSet.Parse(argv)
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.1

- **任务 3.2.3** ✅ 创建 `internal/config/loader.go` - 实现帮助和版本处理
  ```go
  - 检查 --help 和 --version 标志
  - 输出帮助信息到stdout
  - 返回exitCode=0（表示程序应正常退出）
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.2

- **任务 3.2.4** ✅ 创建 `internal/config/loader.go` - 实现位置参数解析
  ```go
  - 使用 flagSet.Args() 获取位置参数
  - 第一个参数是URL（必需）
  - 第二个参数是output_file（可选）
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.2

- **任务 3.2.5** ✅ 创建 `internal/config/loader.go` - 实现环境变量读取
  ```go
  - 使用 os.Getenv("GITHUB_TOKEN") 读取token
  - 赋值给Config.Token字段
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.4

- **任务 3.2.6** ✅ 创建 `internal/config/loader.go` - 构建并返回Config
  ```go
  - 组装Config结构体
  - 返回(*Config, -1)（-1表示不退出）
  - 处理错误情况（如缺少URL参数）
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.5

- **任务 3.2.7** ✅ 运行测试验证通过（Green）
  ```bash
  - 执行 go test ./internal/config -v
  - 确认所有测试通过
  - 测试覆盖率：88.5%
  ```
  **依赖**: 任务 3.2.6

---

### 3.3 重构优化（Refactor Phase）

- **任务 3.3.1** ✅ `[P]` 重构 `internal/config/loader.go` - 提取帮助文本
  ```go
  - 将帮助文本提取为常量
  - 使用 fmt.Fprintf 输出
  - 代码已经很清晰，帮助文本独立为 printHelp 函数
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.2.7

- **任务 3.3.2** ✅ `[P]` 重构 `internal/config/loader.go` - 提取版本信息格式化
  ```go
  - 创建版本信息格式化函数
  - 使用 Version 和 BuildDate 变量
  - 版本信息已独立为 printVersion, getVersion, getBuildDate 函数
  ```
  **文件**: `internal/config/loader.go`
  **依赖**: 任务 3.3.1

- **任务 3.3.3** ✅ `[P]` 运行测试确保重构未破坏功能
  ```bash
  - 执行 go test ./internal/config -v
  - 确认所有测试仍然通过
  - 测试覆盖率：88.5%
  ```
  **依赖**: 任务 3.3.1, 3.3.2

---

## Phase 4: GitHub Fetcher（API客户端，TDD）

> **目标**：实现GitHub API数据获取功能（Issue/PR/Discussion）

### 4.1 测试先行（Red Phase）

- **任务 4.1.1** ✅ `[P]` 编写 `internal/github/client_test.go` - 单元测试框架（使用 httptest）
  ```go
  - 创建测试文件
  - 使用 net/http/httptest 创建 Mock Server
  - 编写测试用例覆盖所有功能
  ```
  **文件**: `internal/github/client_test.go`

- **任务 4.1.2** ✅ `[P]` 编写 `internal/github/client_test.go` - NewClient测试
  ```go
  - 添加测试用例：create client with token
  - 添加测试用例：create client without token
  - 验证Client字段正确初始化
  ```
  **文件**: `internal/github/client_test.go`
  **依赖**: 任务 4.1.1

- **任务 4.1.3** ✅ `[P]` 编写 `internal/github/client_test.go` - FetchIssue单元测试
  ```go
  - 添加测试用例：fetch issue success
  - 添加测试用例：fetch issue with comments
  - 验证Issue字段正确（Title, URL, User等）
  - 验证Comments字段
  ```
  **文件**: `internal/github/client_test.go`
  **依赖**: 任务 4.1.1

- **任务 4.1.4** ✅ `[P]` 编写 `internal/github/client_test.go` - FetchPullRequest单元测试
  ```go
  - 添加测试用例：fetch PR success
  - 添加测试用例：fetch merged PR
  - 验证PR字段正确
  - 验证Comments为Review评论
  ```
  **文件**: `internal/github/client_test.go`
  **依赖**: 任务 4.1.1

- **任务 4.1.5** ✅ `[P]` 编写 `internal/github/client_test.go` - FetchDiscussion单元测试
  ```go
  - 添加测试用例：fetch discussion success（GraphQL）
  - 验证Discussion字段正确（Title, URL, User, State等）
  - 验证Comments包含所有回复
  ```
  **文件**: `internal/github/client_test.go`
  **依赖**: 任务 4.1.1

- **任务 4.1.6** ✅ `[P]` 编写 `internal/github/client_test.go` - 错误处理测试
  ```go
  - 添加测试用例：fetch non-existent issue（404）
  - 验证返回ErrResourceNotFound
  ```
  **文件**: `internal/github/client_test.go`
  **依赖**: 任务 4.1.1

- **任务 4.1.7** ✅ 运行测试验证失败（Red）
  ```bash
  - 执行 go test ./internal/github -v
  - 确认所有测试失败（编译错误）
  ```
  **依赖**: 任务 4.1.2, 4.1.3, 4.1.4, 4.1.5, 4.1.6

---

### 4.2 实现功能（Green Phase）

- **任务 4.2.1** ✅ 创建 `internal/github/client.go` - 定义错误变量
  ```go
  - 定义 ErrResourceNotFound
  - 定义 ErrAPIRateLimit
  - 定义 ErrNetwork
  ```
  **文件**: `internal/github/client.go`

- **任务 4.2.2** ✅ 创建 `internal/github/client.go` - 实现NewClient
  ```go
  - 函数签名：func NewClient(token string, opts ...Option) *Client
  - 支持选项模式（WithBaseURL）
  - 使用 net/http 创建 HTTP 客户端
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.1

- **任务 4.2.3** ✅ 创建 `internal/github/client.go` - 实现HTTP请求辅助函数
  ```go
  - 实现 get() 方法用于 REST API
  - 实现 postGraphQL() 方法用于 GraphQL
  - 实现 buildReactions() 辅助函数
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.2

- **任务 4.2.4** ✅ 创建 `internal/github/client.go` - 实现FetchIssue方法
  ```go
  - 方法签名：func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error)
  - 调用 GitHub REST API 获取 Issue
  - 调用 GitHub REST API 获取评论
  - 转换为内部 Issue 结构体
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.3

- **任务 4.2.5** ✅ 创建 `internal/github/client.go` - 实现FetchPullRequest方法
  ```go
  - 方法签名：func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error)
  - 调用 GitHub REST API 获取 PR
  - 调用 GitHub REST API 获取 Review 评论
  - 处理 merged 状态
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.4

- **任务 4.2.6** ✅ 创建 `internal/github/client.go` - 实现FetchDiscussion方法（GraphQL）
  ```go
  - 方法签名：func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error)
  - 使用 GraphQL API（因为REST不支持Discussion评论）
  - 编写GraphQL查询语句
  - 解析响应并转换为Discussion结构体
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.5

- **任务 4.2.7** ✅ 运行测试验证通过（Green）
  ```bash
  - 执行 go test ./internal/github -v
  - 确认所有测试通过
  - 测试覆盖率：82.6%
  ```
  **依赖**: 任务 4.2.6

---

### 4.3 重构优化（Refactor Phase）

- **任务 4.3.1** ✅ `[P]` 重构 `internal/github/client.go` - 代码已经符合最佳实践
  ```go
  - 使用 Go 标准库（net/http）
  - 错误处理清晰明确
  - 代码结构简单直接
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.2.7

- **任务 4.3.2** ✅ `[P]` 重构 `internal/github/client.go` - 错误处理已优化
  ```go
  - 检查HTTP响应状态码
  - 区分404错误
  - 返回对应的错误类型
  ```
  **文件**: `internal/github/client.go`
  **依赖**: 任务 4.3.1

- **任务 4.3.3** ✅ `[P]` 运行测试确保重构未破坏功能
  ```bash
  - 执行 go test ./internal/github -v
  - 确认所有测试仍然通过
  ```
  **依赖**: 任务 4.3.1, 4.3.2

---

## Phase 5: Markdown Converter（转换器，TDD）

> **目标**：实现GitHub数据到Markdown的转换功能

### 5.1 测试先行（Red Phase）

- **任务 5.1.1** ✅ `[P]` 编写 `internal/converter/converter_test.go` - 测试框架
  ```go
  - 创建测试文件
  - 定义测试辅助函数（创建mock github.Issue）
  - 创建空的测试函数 TestConvertIssue(t *testing.T)
  ```
  **文件**: `internal/converter/converter_test.go`

- **任务 5.1.2** ✅ `[P]` 编写 `internal/converter/converter_test.go` - ConvertIssue基础测试
  ```go
  - 添加测试用例：basic issue conversion
  - 验证输出包含YAML Frontmatter
  - 验证输出包含标题、作者、时间、状态
  - 验证输出包含正文
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.3** ✅ `[P]` 编写 `internal/converter/converter_test.go` - ConvertIssue评论测试
  ```go
  - 添加测试用例：issue with comments
  - 验证输出包含所有评论
  - 验证评论按时间正序排列
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.4** ✅ `[P]` 编写 `internal/converter/converter_test.go` - Reactions和UserLinks测试
  ```go
  - 添加测试用例：with reactions enabled
  - 验证输出包含reactions统计
  - 添加测试用例：with user links enabled
  - 验证用户名格式为[@username](url)
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.5** ✅ `[P]` 编写 `internal/converter/converter_test.go` - Emoji转换测试
  ```go
  - 添加测试用例：emoji shortcode conversion
  - 验证:thumbsup: → 👍
  - 验证:heart: → ❤️
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.6** ✅ `[P]` 编写 `internal/converter/converter_test.go` - ConvertPullRequest测试
  ```go
  - 添加测试用例：basic PR conversion
  - 验证输出包含PR标题、状态
  - 验证状态可能是"merged"
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.7** ✅ `[P]` 编写 `internal/converter/converter_test.go` - ConvertDiscussion测试
  ```go
  - 添加测试用例：basic discussion conversion
  - 验证输出包含所有评论（含回复）
  ```
  **文件**: `internal/converter/converter_test.go`
  **依赖**: 任务 5.1.1

- **任务 5.1.8** ✅ 运行测试验证失败（Red）
  ```bash
  - 执行 go test ./internal/converter -v
  - 确认所有测试失败
  ```
  **依赖**: 任务 5.1.2, 5.1.3, 5.1.4, 5.1.5, 5.1.6, 5.1.7

---

### 5.2 实现功能（Green Phase）

- **任务 5.2.1** ✅ 在 `internal/converter/converter.go` - 实现emoji映射表
  ```go
  - 定义 emojiMap 变量（map[string]string）
  - 包含所有GitHub shortcode到emoji的映射
  ```
  **文件**: `internal/converter/converter.go`

- **任务 5.2.2** ✅ 在 `internal/converter/converter.go` - 实现辅助函数
  ```go
  - 实现 formatYAMLFrontmatter(...) string
  - 实现 formatUser(user github.User) string
  - 实现 formatTimestamp(t time.Time) string
  - 实现 formatReactions(reactions []github.Reaction) string
  - 实现 convertEmojiShortcode(body string) string
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.1

- **任务 5.2.3** ✅ 在 `internal/converter/converter.go` - 实现ConvertIssue方法（骨架）
  ```go
  - 方法签名：func (c *Converter) ConvertIssue(issue *github.Issue) (string, error)
  - 创建strings.Builder实例
  - 调用formatYAMLFrontmatter
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.2

- **任务 5.2.4** ✅ 在 `internal/converter/converter.go` - 实现ConvertIssue方法（正文）
  ```go
  - 生成标题行（# {Title}）
  - 生成元数据行（作者、时间、状态）
  - 添加正文内容
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.3

- **任务 5.2.5** ✅ 在 `internal/converter/converter.go` - 实现ConvertIssue方法（评论）
  ```go
  - 遍历issue.Comments
  - 为每个评论生成标题（### {user} - {time}）
  - 添加评论内容
  - 添加reactions（如果启用）
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.4

- **任务 5.2.6** ✅ 在 `internal/converter/converter.go` - 实现ConvertPullRequest方法
  ```go
  - 类似ConvertIssue的逻辑
  - 处理merged状态
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.5

- **任务 5.2.7** ✅ 在 `internal/converter/converter.go` - 实现ConvertDiscussion方法
  ```go
  - 类似ConvertIssue的逻辑
  - 评论已按时间排序（github包保证）
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.6

- **任务 5.2.8** ✅ 运行测试验证通过（Green）
  ```bash
  - 执行 go test ./internal/converter -v
  - 确认所有测试通过
  ```
  **依赖**: 任务 5.2.7

---

### 5.3 重构优化（Refactor Phase）

- **任务 5.3.1`[P]` 重构 `internal/converter/converter.go` - 优化字符串构建
  ```go
  - 使用strings.Builder替代字符串拼接
  - 预分配容量（estimate size）
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.2.8

- **任务 5.3.2`[P]` 重构 `internal/converter/converter.go` - 提取时间格式化逻辑
  ```go
  - 使用time.Format()而非手动拼接
  - 格式："2006-01-02 15:04:05"
  ```
  **文件**: `internal/converter/converter.go`
  **依赖**: 任务 5.3.1

- **任务 5.3.3`[P]` 运行测试确保重构未破坏功能
  ```bash
  - 执行 go test ./internal/converter -v
  - 确认所有测试仍然通过
  ```
  **依赖**: 任务 5.3.1, 5.3.2

---

## Phase 6: CLI Assembly（命令行集成，TDD）

> **目标**：实现CLI业务逻辑编排，串联所有模块

### 6.1 测试先行（Red Phase）

- **任务 6.1.1`[P]` ✅ 编写 `internal/cli/cli_test.go` - 端到端测试框架
  ```go
  - 创建测试文件
  - 定义测试辅助函数（创建fake stdout/stderr）
  - 创建空的测试函数 TestRun(t *testing.T)
  ```
  **文件**: `internal/cli/cli_test.go`

- **任务 6.1.2`[P]` ✅ 编写 `internal/cli/cli_test.go` - 帮助和版本测试
  ```go
  - 添加测试用例：--help flag
  - 验证exitCode=0
  - 验证stdout包含"Usage:"
  - 添加测试用例：--version flag
  - 验证stdout包含"version:"
  ```
  **文件**: `internal/cli/cli_test.go`
  **依赖**: 任务 6.1.1

- **任务 6.1.3`[P]` ✅ 编写 `internal/cli/cli_test.go` - URL解析错误测试
  ```go
  - 添加测试用例：invalid URL
  - 验证exitCode=1
  - 验证stderr包含错误信息
  ```
  **文件**: `internal/cli/cli_test.go`
  **依赖**: 任务 6.1.1

- **任务 6.1.4`[P]` ✅ 编写 `internal/cli/cli_test.go` - GitHub API错误测试
  ```go
  - 添加测试用例：issue not found（使用mock github client）
  - 验证exitCode=1
  - 验证stderr包含"not found"
  ```
  **文件**: `internal/cli/cli_test.go`
  **依赖**: 任务 6.1.1

- **任务 6.1.5`[P]` ✅ 编写 `internal/cli/cli_test.go` - 完整流程测试（跳过，需要真实token）
  ```go
  - 添加测试用例：end-to-end issue conversion
  - 使用真实的public issue
  - 验证stdout输出包含Markdown
  - 添加skip条件（无GITHUB_TOKEN时跳过）
  ```
  **文件**: `internal/cli/cli_test.go`
  **依赖**: 任务 6.1.1

- **任务 6.1.6 ✅ 运行测试验证失败（Red）
  ```bash
  - 执行 go test ./internal/cli -v
  - 确认所有测试失败
  ```
  **依赖**: 任务 6.1.2, 6.1.3, 6.1.4, 6.1.5

---

### 6.2 实现功能（Green Phase）

- **任务 6.2.1 ✅ 创建 `internal/cli/cli.go` - 实现Run函数骨架
  ```go
  - 函数签名：func Run(argv []string, stdout, stderr io.Writer) int
  - 调用config.LoadFromFlags
  - 处理立即退出的情况（--help, --version）
  ```
  **文件**: `internal/cli/cli.go`

- **任务 6.2.2 ✅ 在 `internal/cli/cli.go` - 实现URL解析逻辑
  ```go
  - 调用parser.ParseURL
  - 捕获错误并输出到stderr
  - 返回exitCode=1
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.2.1

- **任务 6.2.3 ✅ 在 `internal/cli/cli.go` - 实现GitHub客户端创建
  ```go
  - 调用github.NewClient(cfg.Token)
  - 根据resource.Type分发到不同的Fetch方法
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.2.2

- **任务 6.2.4 ✅ 在 `internal/cli/cli.go` - 实现数据获取和转换逻辑
  ```go
  - 调用FetchIssue/FetchPullRequest/FetchDiscussion
  - 创建converter.Converter实例
  - 调用Convert方法
  - 处理错误
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.2.3

- **任务 6.2.5 ✅ 在 `internal/cli/cli.go` - 实现输出逻辑
  ```go
  - 检查cfg.OutputFile
  - 如果为空，输出到stdout
  - 如果不为空，写入文件
  - 处理文件写入错误
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.2.4

- **任务 6.2.6 ✅ 运行测试验证通过（Green）
  ```bash
  - 执行 go test ./internal/cli -v
  - 确认所有测试通过
  ```
  **依赖**: 任务 6.2.5

---

### 6.3 重构优化（Refactor Phase）

- **任务 6.3.1`[P]` 重构 `internal/cli/cli.go` - 提取资源类型分发逻辑
  ```go
  - 将根据type分发到不同Fetch方法的逻辑提取为独立函数
  - 函数名：fetchResource(client *github.Client, resource *parser.Resource)
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.2.6

- **任务 6.3.2`[P]` 重构 `internal/cli/cli.go` - 提取输出逻辑
  ```go
  - 将输出逻辑提取为独立函数
  - 函数名：writeOutput(markdown string, outputFile string, stdout io.Writer) error
  ```
  **文件**: `internal/cli/cli.go`
  **依赖**: 任务 6.3.1

- **任务 6.3.3`[P]` 运行测试确保重构未破坏功能
  ```bash
  - 执行 go test ./internal/cli -v
  - 确认所有测试仍然通过
  ```
  **依赖**: 任务 6.3.1, 6.3.2

---

## Phase 7: Main Entry Point（入口点）

> **目标**：实现cmd/issue2md/main.go入口文件

### 7.1 实现功能

- **任务 7.1.1** ✅ 创建 `cmd/issue2md/main.go`
  ```go
  - package main
  - import "github.com/issue2md/internal/cli"
  - func main() { os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr)) }
  ```
  **文件**: `cmd/issue2md/main.go`

---

## Phase 8: Build & Documentation（构建和文档）

> **目标**：完善构建脚本和文档

### 8.1 构建脚本

- **任务 8.1.1** ✅ `[P]` 创建 `Makefile`
  ```makefile
  - 定义build目标
  - 定义test目标
  - 定义clean目标
  - 定义version注入（-ldflags）
  ```
  **文件**: `Makefile`

- **任务 8.1.2** ✅ `[P]` 创建 `README.md`
  ```markdown
  - 项目简介
  - 安装说明
  - 使用示例
  - 环境变量说明（GITHUB_TOKEN）
  - 构建说明
  ```
  **文件**: `README.md`

---

### 8.2 验收测试

- **任务 8.2.1** ✅ 手动验收测试 - Issue转换
  ```bash
  - 构建二进制文件
  - 运行：./issue2md https://github.com/golang/go/issues/1
  - 验证输出包含正确的Markdown
  ```

- **任务 8.2.2** ✅ 手动验收测试 - PR转换
  ```bash
  - 运行：./issue2md https://github.com/golang/go/pull/1
  - 验证输出包含PR信息
  ```

- **任务 8.2.3** ✅ 手动验收测试 - Discussion转换
  ```bash
  - 运行：./issue2md https://github.com/github/discussions/1
  - 验证输出包含Discussion信息
  ```

- **任务 8.2.4** ✅ 手动验收测试 - Reactions和UserLinks
  ```bash
  - 运行：./issue2md -enable-reactions -enable-user-links <url>
  - 验证输出包含reactions和用户链接
  ```

- **任务 8.2.5** ✅ 手动验收测试 - 输出到文件
  ```bash
  - 运行：./issue2md <url> output.md
  - 验证output.md文件创建成功
  - 验证文件内容正确
  ```

- **任务 8.2.6** ✅ 手动验收测试 - 错误处理
  ```bash
  - 运行：./issue2md https://invalid-url.com
  - 验证返回exitCode=1
  - 验证stderr输出错误信息
  ```

---

## Phase 9: Code Review & Polish（代码审查和优化）

> **目标**：代码质量检查和优化

### 9.1 代码质量

- **任务 9.1.1`[P]` 运行 go vet 检查
  ```bash
  - 执行 go vet ./...
  - 修复所有警告
  ```

- **任务 9.1.2`[P]` 运行 gofmt 检查
  ```bash
  - 执行 gofmt -l .
  - 格式化所有不符合规范的文件
  ```

- **任务 9.1.3`[P]` 运行完整测试套件
  ```bash
  - 执行 go test -v -race -cover ./...
  - 确保所有测试通过
  - 确保无竞态条件
  - 检查测试覆盖率
  ```

- **任务 9.1.4`[P]` 执行交叉编译
  ```bash
  - 构建Linux版本
  - 构建macOS版本
  - 构建Windows版本
  ```

---

## 任务统计

| Phase | 任务数 | 可并行任务 | 核心任务 |
|-------|--------|-----------|---------|
| Phase 1: Foundation | 12 | 12 | 数据结构定义 |
| Phase 2: URL Parser | 10 | 3 | URL解析功能 |
| Phase 3: Config Loader | 10 | 7 | 配置加载功能 |
| Phase 4: GitHub Fetcher | 13 | 6 | API客户端功能 |
| Phase 5: Converter | 13 | 7 | Markdown转换功能 |
| Phase 6: CLI Assembly | 9 | 5 | CLI业务逻辑 |
| Phase 7: Main | 1 | 0 | 入口点 |
| Phase 8: Build & Doc | 7 | 2 | 构建和文档 |
| Phase 9: Polish | 4 | 4 | 代码质量 |
| **总计** | **79** | **46** | - |

---

## 依赖关系图

```
Phase 1 (Foundation)
    ├─> Phase 2 (Parser)
    │       └─> Phase 3 (Config Loader) [P]
    │               └─> Phase 6 (CLI Assembly)
    ├─> Phase 4 (GitHub Fetcher)
    │       └─> Phase 5 (Converter)
    │               └─> Phase 6 (CLI Assembly)
    └─> Phase 7 (Main Entry)
            └─> Phase 8 (Build & Doc)
                    └─> Phase 9 (Polish)
```

---

## 执行建议

### 优先级顺序
1. **高优先级**：Phase 1 → Phase 2 → Phase 4 → Phase 5 → Phase 6 → Phase 7
2. **中优先级**：Phase 3（可与Phase 4并行）
3. **低优先级**：Phase 8 → Phase 9

### 并行策略
- **同一Phase内**：标记`[P]`的任务可以并行执行
- **不同Phase间**：必须等待前置依赖完成

### 时间估算
- **Phase 1**: 0.5天（数据结构定义）
- **Phase 2**: 1天（URL解析器）
- **Phase 3**: 1天（配置加载）
- **Phase 4**: 2天（GitHub API客户端）
- **Phase 5**: 2天（Markdown转换）
- **Phase 6**: 1天（CLI集成）
- **Phase 7**: 0.5天（入口点）
- **Phase 8**: 0.5天（构建和文档）
- **Phase 9**: 0.5天（代码质量）
- **总计**: 约8.5天

---

**文档结束**
