# issue2md 核心功能需求规范

## 📋 文档信息
- **版本**: 1.0
- **创建日期**: 2025-01-04
- **状态**: Draft
- **优先级**: P0 (MVP)

---

## 🎯 用户故事

### 当前迭代：CLI 工具 (MVP)

**作为** 个人开发者
**我想要** 通过命令行工具输入 GitHub Issue/PR/Discussion 的 URL
**以便于** 将讨论内容快速转换为格式化的 Markdown 文件，用于归档和知识管理

**验收标准**：
- ✅ 支持公有仓库的 Issue/PR/Discussion
- ✅ 输出包含完整讨论内容的 Markdown
- ✅ 输出到 stdout 或指定文件
- ✅ 通过环境变量安全传入认证信息

---

### 未来迭代：Web 界面 (Post-MVP)

**作为** 个人开发者
**我想要** 通过 Web 界面粘贴 GitHub URL
**以便于** 在浏览器中直接获取 Markdown 内容并下载

**注意**：此功能不包含在当前 MVP 范围内。

---

## 📐 功能性需求

### 1. 输入处理

#### 1.1 URL 格式支持
工具必须能够识别并解析以下三种 GitHub URL 格式：

| 类型 | URL 模式 | 示例 |
|------|----------|------|
| Issue | `https://github.com/{owner}/{repo}/issues/{number}` | `https://github.com/owner/repo/issues/123` |
| PR | `https://github.com/{owner}/{repo}/pull/{number}` | `https://github.com/owner/repo/pull/456` |
| Discussion | `https://github.com/{owner}/{repo}/discussions/{number}` | `https://github.com/owner/repo/discussions/789` |

**约束**：
- ❌ 不支持简化格式（如 `owner/repo#123`）
- ❌ 不处理 URL 查询参数（如 `?sort=oldest`）
- ✅ URL 格式错误时，程序退出并输出清晰的错误信息到 stderr

#### 1.2 命令行参数

```bash
issue2md [flags] [output_file]
```

**参数定义**：
- `output_file` (位置参数，可选)
  - 不提供：输出到 stdout
  - 提供文件路径：写入指定文件
- `-enable-reactions` (flag，可选)
  - 启用后，在评论末尾显示 reactions 统计（如 `👍 3 ❤️ 1`）
- `-enable-user-links` (flag，可选)
  - 启用后，用户名显示为可点击链接（如 `[@username](https://github.com/username)`）
- `--version`
  - 输出版本号
- `--help`
  - 输出使用帮助

#### 1.3 认证机制

**仅支持**通过环境变量 `GITHUB_TOKEN` 传入 Personal Access Token：
```bash
export GITHUB_TOKEN=ghp_xxxxxxxxxxxx
issue2md https://github.com/owner/repo/issues/123
```

**安全约束**：
- ❌ 不提供 `--token` 参数（防止 Shell 历史泄露密钥）
- ✅ 公有仓库可不提供 Token，但受 API 限流限制
- ✅ 推荐 Token 用于提高 API 限流配额

---

### 2. GitHub API 集成

#### 2.1 API 版本
- **使用**: GitHub REST API v3
- **Base URL**: `https://api.github.com`

#### 2.2 数据获取策略

| 资源类型 | 核心数据 | 排除内容 |
|----------|----------|----------|
| **Issue** | 标题、作者、创建时间、状态、正文、所有评论 | 附件下载 |
| **PR** | 标题、作者、创建时间、状态、正文、所有 Review 评论 | Diff、提交历史 |
| **Discussion** | 标题、作者、创建时间、状态、正文、所有评论（含回复） | - |

#### 2.3 评论排序规则
- **统一规则**: 按创建时间正序（从旧到新）
- **Discussion 特殊处理**:
  - 主楼评论和回复平铺展示，不保留缩进层级
  - 目标：归档"发生了什么对话"

#### 2.4 错误处理
| 错误场景 | 处理方式 |
|----------|----------|
| URL 格式无效 | 退出码 1，stderr 输出错误信息 |
| Issue/PR/Discussion 不存在 | 退出码 1，stderr 输出错误信息 |
| 网络错误 | 退出码 1，stderr 输出错误信息 |
| API 限流 (Rate Limit) | 透传 GitHub API 错误信息给用户 |
| 内容已被删除 | 在 Markdown 中用删除线标注（~~deleted~~） |

---

### 3. Markdown 输出格式

#### 3.1 文件结构

```markdown
---
title: "{标题}"
url: "{原始URL}"
author: "{@username}"
created_at: "YYYY-MM-DD HH:MM:SS"
status: "{open|closed|merged}"
---

# {标题}

**作者**: [@username](https://github.com/username)
**创建时间**: 2025-01-04 12:30:45
**状态**: Open

{主楼正文内容}

## 评论

### [@username1](https://github.com/username1) - 2025-01-04 13:00:00

{评论内容}

### [@username2](https://github.com/username2) - 2025-01-04 14:30:00

{评论内容}
👍 3 ❤️ 1
```

#### 3.2 YAML Frontmatter 字段

| 字段 | 类型 | 说明 | 示例 |
|------|------|------|------|
| `title` | string | Issue/PR/Discussion 标题 | `"Fix authentication bug"` |
| `url` | string | 原始 GitHub URL | `"https://github.com/owner/repo/issues/123"` |
| `author` | string | 作者用户名（带 @ 前缀） | `"@johndoe"` |
| `created_at` | string | 创建时间（本地化格式） | `"2025-01-04 12:30:45"` |
| `status` | string | 当前状态 | `"open"` / `"closed"` / `"merged"` |

#### 3.3 时间格式
- **格式**: `YYYY-MM-DD HH:MM:SS`（24小时制）
- **时区**: 本地时区

#### 3.4 用户名格式
| Flag 关闭 | Flag 开启 (`-enable-user-links`) |
|-----------|----------------------------------|
| `@username` | `[@username](https://github.com/username)` |

#### 3.5 Reactions 格式（`-enable-reactions` 启用时）
```markdown
👍 3 ❤️ 1 😄 2 🎉 1
```
- 仅统计数量 > 0 的 reactions
- 显示顺序：👍❤️😄🎉（对应 GitHub 的 +1, heart, laugh, hooray）

#### 3.6 内容安全处理
- **HTML 标签**: 转义为 Markdown 安全格式
- **GitHub Emoji shortcode**: 转换为 Unicode emoji
  - `:thumbsup:` → 👍
  - `:heart:` → ❤️

---

## 🏗️ 非功能性需求

### 4.1 架构设计原则

#### 解耦设计
- **核心业务逻辑**与**CLI 界面**分离
- **GitHub API 客户端**独立为可测试模块
- **Markdown 生成器**独立为可复用组件（为未来 Web 版预留）

#### 依赖管理
- ✅ 仅使用 Go 标准库
- ✅ 不引入第三方 HTTP 客户端或 CLI 框架
- ✅ 使用 `net/http`, `encoding/json`, `time`, `fmt` 等标准包

### 4.2 错误处理规范
- **所有错误必须被显式处理**（遵循 Go 惯用法）
- **错误传递**使用 `fmt.Errorf("context: %w", err)` 包装
- **用户可见错误**输出到 stderr，包含清晰的上下文信息
- **退出码**：
  - `0`: 成功
  - `1`: 任何错误（URL 无效、网络错误、API 错误等）

### 4.3 测试策略
- **表格驱动测试**（Table-Driven Tests）
- **集成测试**优先，使用真实的 GitHub API（测试仓库）
- **拒绝 Mocks**，确保测试真实可靠

---

## ✅ 验收标准

### 测试用例

#### TC-01: Issue 基础转换
**输入**:
```bash
issue2md https://github.com/owner/repo/issues/123
```

**预期输出**:
- ✅ 输出包含 YAML Frontmatter
- ✅ 输出包含标题、作者、创建时间、状态
- ✅ 输出包含主楼正文
- ✅ 输出包含所有评论（按时间正序）
- ✅ 用户名格式为 `@username`（未启用 `-enable-user-links`）

---

#### TC-02: 启用用户链接和 Reactions
**输入**:
```bash
issue2md -enable-user-links -enable-reactions https://github.com/owner/repo/issues/123
```

**预期输出**:
- ✅ 用户名格式为 `[@username](https://github.com/username)`
- ✅ 评论末尾显示 reactions 统计（如 `👍 3`）

---

#### TC-03: PR 转换（不包含 Diff）
**输入**:
```bash
issue2md https://github.com/owner/repo/pull/456
```

**预期输出**:
- ✅ 输出包含 PR 描述
- ✅ 输出包含所有 Review 评论
- ❌ 不包含代码变更或提交历史

---

#### TC-04: Discussion 转换（平铺回复）
**输入**:
```bash
issue2md https://github.com/owner/repo/discussions/789
```

**预期输出**:
- ✅ 输出包含主楼评论和所有回复
- ✅ 所有内容按时间正序平铺展示（无缩进）

---

#### TC-05: 输出到文件
**输入**:
```bash
issue2md https://github.com/owner/repo/issues/123 output.md
```

**预期输出**:
- ✅ 文件 `output.md` 被创建
- ✅ 文件内容包含完整的 Markdown
- ✅ stdout 无输出（或仅输出进度信息）

---

#### TC-06: 无效 URL 错误处理
**输入**:
```bash
issue2md https://invalid-url.com/abc
```

**预期输出**:
- ❌ 退出码为 1
- ✅ stderr 输出清晰的错误信息（如 "invalid GitHub URL format"）

---

#### TC-07: 不存在的 Issue
**输入**:
```bash
issue2md https://github.com/owner/repo/issues/999999
```

**预期输出**:
- ❌ 退出码为 1
- ✅ stderr 输出错误信息（如 "issue not found"）

---

#### TC-08: 已删除内容的处理
**场景**: Issue 的某条评论被删除

**预期输出**:
- ✅ 在评论位置显示 `~~deleted~~`
- ✅ 不中断程序执行

---

#### TC-09: Emoji 转换
**场景**: Issue 内容包含 `:thumbsup:` 和 `:heart:`

**预期输出**:
- ✅ `:thumbsup:` 转换为 👍
- ✅ `:heart:` 转换为 ❤️

---

#### TC-10: 帮助和版本信息
**输入**:
```bash
issue2md --help
issue2md --version
```

**预期输出**:
- ✅ `--help` 输出使用说明
- ✅ `--version` 输出版本号

---

## 📄 输出格式示例

### 示例 1: 基础 Issue 转换（无 Flag）

```markdown
---
title: "Fix authentication bug in login flow"
url: "https://github.com/example/myproject/issues/42"
author: "@johndoe"
created_at: "2025-01-04 10:30:00"
status: "open"
---

# Fix authentication bug in login flow

**作者**: @johndoe
**创建时间**: 2025-01-04 10:30:00
**状态**: Open

## 描述

用户在使用 OAuth2 登录时，access_token 过期后没有正确刷新。

### 复现步骤
1. 登录应用
2. 等待 1 小时（token 过期）
3. 尝试访问需要认证的 API
4. 返回 401 错误

### 期望行为
应该自动刷新 token 并重试请求。

### 实际行为
直接返回 401 错误，用户体验差。

## 评论

### @alice - 2025-01-04 11:15:00

我也遇到了这个问题，建议在 HTTP 客户端中间件中添加自动重试逻辑。

### @bob - 2025-01-04 12:00:00

同意。我已经在 PR #45 中实现了这个功能，请 review 一下。

### @johndoe - 2025-01-04 14:30:00

@bob 看了你的 PR，实现得很好！已经合并了 👍

### @carol - 2025-01-04 15:00:00

~~deleted~~
```

---

### 示例 2: 启用 Reactions 和用户链接

```markdown
---
title: "Add support for custom templates"
url: "https://github.com/example/myproject/issues/100"
author: "@sarahr"
created_at: "2025-01-03 09:00:00"
status: "closed"
---

# Add support for custom templates

**作者**: [@sarahr](https://github.com/sarahr)
**创建时间**: 2025-01-03 09:00:00
**状态**: Closed

## 描述

希望支持用户自定义 Markdown 模板，以便适配不同的文档系统。

## 评论

### [@mikedev](https://github.com/mikedev) - 2025-01-03 10:30:00

好主意！建议使用 Go template 语法。

👍 3 ❤️ 1

### [@sarahr](https://github.com/sarahr) - 2025-01-03 11:00:00

同意，我会设计一个简单的模板系统。

🎉 2

### [@tester](https://github.com/tester) - 2025-01-03 12:00:00

可以参考 Hugo 的模板设计。

😄 1
```

---

### 示例 3: PR 转换（仅描述和 Review 评论）

```markdown
---
title: "Implement auto-retry logic for OAuth2 token refresh"
url: "https://github.com/example/myproject/pull/45"
author: "@bob"
created_at: "2025-01-04 11:00:00"
status: "merged"
---

# Implement auto-retry logic for OAuth2 token refresh

**作者**: @bob
**创建时间**: 2025-01-04 11:00:00
**状态**: Merged

## 描述

在 HTTP 客户端中间件中添加自动重试逻辑，当收到 401 错误时自动刷新 access_token 并重试请求。

### 变更内容
- 添加 `TokenRefresher` 中间件
- 处理 401 错误并刷新 token
- 最多重试 3 次

## 评论

### @reviewer1 - 2025-01-04 13:00:00

整体实现很好，但建议把重试次数提取为配置参数。

👍 1

### @bob - 2025-01-04 13:30:00

好建议，已更新为 `MaxRetries` 配置项。

### @reviewer1 - 2025-01-04 14:00:00

LGTM, approved! 🎉

👍 1 🎉 1
```

---

## 📚 附录

### A. GitHub API 端点参考

| 资源 | API 端点 |
|------|----------|
| Issue | `GET /repos/{owner}/{repo}/issues/{issue_number}` |
| Issue Comments | `GET /repos/{owner}/{repo}/issues/{issue_number}/comments` |
| PR | `GET /repos/{owner}/{repo}/pulls/{pull_number}` |
| PR Review Comments | `GET /repos/{owner}/{repo}/pulls/{pull_number}/comments` |
| Discussion | `GET /repos/{owner}/{repo}/discussions/{discussion_number}` |
| Discussion Comments | `GET /repos/{owner}/{repo}/discussions/{discussion_number}/comments` |

### B. GitHub Reactions 类型映射

| GitHub API 类型 | Unicode Emoji | Shortcode |
|-----------------|---------------|-----------|
| `+1` | 👍 | `:thumbsup:` / `:+1:` |
| `-1` | 👎 | `:thumbsdown:` / `:-1:` |
| `laugh` | 😄 | `:smile:` / `:laugh:` |
| `hooray` | 🎉 | `:tada:` / `:hooray:` |
| `confused` | 😕 | `:confused:` |
| `heart` | ❤️ | `:heart:` |
| `rocket` | 🚀 | `:rocket:` |
| `eyes` | 👀 | `:eyes:` |

---

**文档结束**
