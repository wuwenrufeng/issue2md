---
name: generate-commit-message
description: creates conventional commit messages based on staged changes. Use this skill whenever a user asks to generate git commit message.
allowed-tools: Bash(git add:_), Bash(git status:_), Bash(git commit:\*)
---

# Generate Commit Message Skill.

This skill creates conventional commit messages based on staged changes. Analyze the provided git diff output and generate appropriate conventional commit messages following the specification.

## Capabilities

Use 中文 generate commit messages

## How to Use

1. Determine Primary Type based on the nature of changes
2. Identify Scope from modified directories or modules
3. Craft Description focusing on the most significant change
4. Determine if there are Breaking Changes
5. For complex changes, include a detailed body explaining what and why
6. Add appropriate footers for issue references or breaking changes

For significant changes, include a detailed body explaining the changes.

Return ONLY the commit message in the conventional format, nothing else.

## Output Format

- DO NOT include any memory bank status indicators like "[Memory Bank: Active]" or "[Memory Bank: Missing]"
- DO NOT include any task-specific formatting or artifacts from other rules
- Use 中文 generate commit messages
- ONLY Generate a clean conventional commit message as specified below

### Context

-## Conventional Commits Format
Generate commit messages following this exact structure:

```
<type>[optional scope]: <description>
[optional body]
[optional footer(s)]
```

#### Core Types (Required)

- **feat**: New feature or functionality (MINOR version bump)
- **fix**: Bug fix or error correction (PATCH version bump)

### Additional Types (Extended)

- **docs**: Documentation changes only
- **style**: Code style changes (whitespace, formatting, semicolons, etc.)
- **refactor**: Code refactoring without feature changes or bug fixes
- **perf**: Performance improvements
- **test**: Adding or fixing tests
- **build**: Build system or external dependency changes
- **ci**: CI/CD configuration changes
- **chore**: Maintenance tasks, tooling changes
- **revert**: Reverting previous commits

### Scope Guidelines

- Use parentheses: `feat(api):`, `fix(ui):`
- Common scopes: `api`, `ui`, `auth`, `db`, `config`, `deps`, `docs`
- For monorepos: package or module names
- Keep scope concise and lowercase

### Description Rules

- Use imperative mood ("add" not "added" or "adds")
- Start with lowercase letter
- No period at the end
- Maximum 50 characters
- Be concise but descriptive

### Body Guidelines (Optional)

- Start one blank line after description
- Explain the "what" and "why", not the "how"
- Wrap at 72 characters per line
- Use for complex changes requiring explanation

### Footer Guidelines (Optional)

- Start one blank line after body
- **Breaking Changes**: `BREAKING CHANGE: description`

## Example Usage

feat: 初始化项目核心模块结构

- 添加 CLI 版本信息模块 (internal/cli/version.go)
- 新增应用配置结构体 (internal/config/config.go)
- 实现 Markdown 转换器基础框架 (internal/converter/converter.go)
- 定义 GitHub API 数据结构 (internal/github/types.go)
- 创建 URL 解析器类型定义 (internal/parser/types.go)
- 降级 Go 版本至 1.24.9 以提升兼容性
