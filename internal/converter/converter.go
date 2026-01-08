package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/wuwenrufeng/issue2md/internal/github"
)

// emojiMap GitHub shortcode 到 Unicode emoji 的映射表
var emojiMap = map[string]string{
	// Thumbs up/down
	":thumbsup:":   "👍",
	":+1:":         "👍",
	":thumbsdown:": "👎",
	":-1:":         "👎",
	// Heart
	":heart:": "❤️",
	// Laugh/smile
	":smile:": "😄",
	":laugh:": "😄",
	// Hooray/tada
	":tada:":   "🎉",
	":hooray:": "🎉",
	// Confused
	":confused:": "😕",
	// Rocket
	":rocket:": "🚀",
	// Eyes
	":eyes:": "👀",
}

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

// NewConverter 创建新的Converter
func NewConverter(options ...Option) *Converter {
	c := &Converter{
		enableReactions: false,
		enableUserLinks: false,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

// formatYAMLFrontmatter 格式化 YAML Frontmatter
func (c *Converter) formatYAMLFrontatter(title, url, author, createdAt, status string) string {
	return fmt.Sprintf("---\ntitle: %q\nurl: %q\nauthor: %q\ncreated_at: %q\nstatus: %q\n---\n\n",
		title, url, author, createdAt, status)
}

// formatUser 格式化用户名
func (c *Converter) formatUser(user github.User) string {
	if c.enableUserLinks {
		return fmt.Sprintf("[@%s](%s)", user.Login, user.HTMLURL)
	}
	return fmt.Sprintf("@%s", user.Login)
}

// formatTimestamp 格式化时间戳
func (c *Converter) formatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// formatReactions 格式化 reactions
func (c *Converter) formatReactions(reactions []github.Reaction) string {
	if !c.enableReactions || len(reactions) == 0 {
		return ""
	}

	var builder strings.Builder
	// 按照 spec.md 定义的顺序：👍❤️😄🎉
	order := map[string]int{
		"+1":     1,
		"heart":  2,
		"laugh":  3,
		"hooray": 4,
	}

	// 按类型分组统计
	counts := make(map[string]int)
	for _, r := range reactions {
		counts[r.Content] += r.Count
	}

	// 按顺序输出
	type reactionItem struct {
		emoji string
		count int
		order int
	}
	var items []reactionItem
	emojiLookup := map[string]string{
		"+1":     "👍",
		"heart":  "❤️",
		"laugh":  "😄",
		"hooray": "🎉",
	}

	for content, count := range counts {
		if count > 0 {
			if emoji, ok := emojiLookup[content]; ok {
				items = append(items, reactionItem{
					emoji: emoji,
					count: count,
					order: order[content],
				})
			}
		}
	}

	// 排序
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[i].order > items[j].order {
				items[i], items[j] = items[j], items[i]
			}
		}
	}

	// 构建输出
	for i, item := range items {
		if i > 0 {
			builder.WriteString(" ")
		}
		builder.WriteString(fmt.Sprintf("%s %d", item.emoji, item.count))
	}

	return builder.String()
}

// convertEmojiShortcode 转换 emoji shortcode 为 Unicode emoji
func (c *Converter) convertEmojiShortcode(body string) string {
	result := body
	for shortcode, emoji := range emojiMap {
		result = strings.ReplaceAll(result, shortcode, emoji)
	}
	return result
}

// ConvertIssue 转换 Issue 为 Markdown
func (c *Converter) ConvertIssue(issue *github.Issue) (string, error) {
	var builder strings.Builder

	// 1. YAML Frontmatter
	author := fmt.Sprintf("@%s", issue.User.Login)
	createdAt := c.formatTimestamp(issue.CreatedAt)
	builder.WriteString(c.formatYAMLFrontatter(
		issue.Title,
		issue.URL,
		author,
		createdAt,
		issue.State,
	))

	// 2. 标题
	builder.WriteString(fmt.Sprintf("# %s\n\n", issue.Title))

	// 3. 元数据
	builder.WriteString(fmt.Sprintf("**作者**: %s\n", c.formatUser(issue.User)))
	builder.WriteString(fmt.Sprintf("**创建时间**: %s\n", createdAt))
	statusDisplay := strings.Title(issue.State)
	builder.WriteString(fmt.Sprintf("**状态**: %s\n\n", statusDisplay))

	// 4. 正文
	if issue.Body != "" {
		body := c.convertEmojiShortcode(issue.Body)
		builder.WriteString(body)
		builder.WriteString("\n\n")
	}

	// 5. 评论
	if len(issue.Comments) > 0 {
		builder.WriteString("## 评论\n\n")
		for _, comment := range issue.Comments {
			// 评论标题
			commentTime := c.formatTimestamp(comment.CreatedAt)
			builder.WriteString(fmt.Sprintf("### %s - %s\n\n", c.formatUser(comment.User), commentTime))

			// 评论内容
			if comment.Deleted {
				builder.WriteString("~~deleted~~\n\n")
			} else if comment.Body != "" {
				commentBody := c.convertEmojiShortcode(comment.Body)
				builder.WriteString(commentBody)
				builder.WriteString("\n\n")
			}

			// Reactions
			reactions := c.formatReactions(comment.Reactions)
			if reactions != "" {
				builder.WriteString(reactions)
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String(), nil
}

// ConvertPullRequest 转换 PullRequest 为 Markdown
func (c *Converter) ConvertPullRequest(pr *github.PullRequest) (string, error) {
	var builder strings.Builder

	// 1. YAML Frontmatter
	author := fmt.Sprintf("@%s", pr.User.Login)
	createdAt := c.formatTimestamp(pr.CreatedAt)
	builder.WriteString(c.formatYAMLFrontatter(
		pr.Title,
		pr.URL,
		author,
		createdAt,
		pr.State,
	))

	// 2. 标题
	builder.WriteString(fmt.Sprintf("# %s\n\n", pr.Title))

	// 3. 元数据
	builder.WriteString(fmt.Sprintf("**作者**: %s\n", c.formatUser(pr.User)))
	builder.WriteString(fmt.Sprintf("**创建时间**: %s\n", createdAt))
	statusDisplay := strings.Title(pr.State)
	builder.WriteString(fmt.Sprintf("**状态**: %s\n\n", statusDisplay))

	// 4. 正文
	if pr.Body != "" {
		body := c.convertEmojiShortcode(pr.Body)
		builder.WriteString(body)
		builder.WriteString("\n\n")
	}

	// 5. 评论 (Review 评论)
	if len(pr.Comments) > 0 {
		builder.WriteString("## 评论\n\n")
		for _, comment := range pr.Comments {
			// 评论标题
			commentTime := c.formatTimestamp(comment.CreatedAt)
			builder.WriteString(fmt.Sprintf("### %s - %s\n\n", c.formatUser(comment.User), commentTime))

			// 评论内容
			if comment.Deleted {
				builder.WriteString("~~deleted~~\n\n")
			} else if comment.Body != "" {
				commentBody := c.convertEmojiShortcode(comment.Body)
				builder.WriteString(commentBody)
				builder.WriteString("\n\n")
			}

			// Reactions
			reactions := c.formatReactions(comment.Reactions)
			if reactions != "" {
				builder.WriteString(reactions)
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String(), nil
}

// ConvertDiscussion 转换 Discussion 为 Markdown
func (c *Converter) ConvertDiscussion(discussion *github.Discussion) (string, error) {
	var builder strings.Builder

	// 1. YAML Frontmatter
	author := fmt.Sprintf("@%s", discussion.User.Login)
	createdAt := c.formatTimestamp(discussion.CreatedAt)
	builder.WriteString(c.formatYAMLFrontatter(
		discussion.Title,
		discussion.URL,
		author,
		createdAt,
		discussion.State,
	))

	// 2. 标题
	builder.WriteString(fmt.Sprintf("# %s\n\n", discussion.Title))

	// 3. 元数据
	builder.WriteString(fmt.Sprintf("**作者**: %s\n", c.formatUser(discussion.User)))
	builder.WriteString(fmt.Sprintf("**创建时间**: %s\n", createdAt))
	statusDisplay := strings.Title(discussion.State)
	builder.WriteString(fmt.Sprintf("**状态**: %s\n\n", statusDisplay))

	// 4. 正文
	if discussion.Body != "" {
		body := c.convertEmojiShortcode(discussion.Body)
		builder.WriteString(body)
		builder.WriteString("\n\n")
	}

	// 5. 评论（包含主楼和所有回复，已按时间排序）
	if len(discussion.Comments) > 0 {
		builder.WriteString("## 评论\n\n")
		for _, comment := range discussion.Comments {
			// 评论标题
			commentTime := c.formatTimestamp(comment.CreatedAt)
			builder.WriteString(fmt.Sprintf("### %s - %s\n\n", c.formatUser(comment.User), commentTime))

			// 评论内容
			if comment.Deleted {
				builder.WriteString("~~deleted~~\n\n")
			} else if comment.Body != "" {
				commentBody := c.convertEmojiShortcode(comment.Body)
				builder.WriteString(commentBody)
				builder.WriteString("\n\n")
			}

			// Reactions
			reactions := c.formatReactions(comment.Reactions)
			if reactions != "" {
				builder.WriteString(reactions)
				builder.WriteString("\n\n")
			}
		}
	}

	return builder.String(), nil
}
