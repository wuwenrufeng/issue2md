package converter

import (
	"strings"
	"testing"
	"time"

	"github.com/wuwenrufeng/issue2md/internal/github"
)

// 辅助函数：创建测试用的 Issue
func createTestIssue(title, body string, comments []github.Comment) *github.Issue {
	return &github.Issue{
		Title: title,
		URL:   "https://github.com/test/repo/issues/1",
		User: github.User{
			Login:   "testuser",
			HTMLURL: "https://github.com/testuser",
		},
		CreatedAt: time.Date(2025, 1, 4, 10, 30, 0, 0, time.Local),
		State:     "open",
		Body:      body,
		Comments:  comments,
	}
}

// 辅助函数：创建测试用的 Comment
func createTestComment(login, body string, createdAt time.Time, reactions []github.Reaction) github.Comment {
	return github.Comment{
		ID:        1,
		User:      github.User{Login: login, HTMLURL: "https://github.com/" + login},
		CreatedAt: createdAt,
		Body:      body,
		Reactions: reactions,
		Deleted:   false,
	}
}

// TestConvertIssue_Basic 测试基本的 Issue 转换功能
func TestConvertIssue_Basic(t *testing.T) {
	tests := []struct {
		name     string
		issue    *github.Issue
		options  []Option
		validate func(t *testing.T, output string)
	}{
		{
			name:  "basic issue conversion",
			issue: createTestIssue("Test Issue", "This is a test issue", nil),
			validate: func(t *testing.T, output string) {
				// 验证 YAML Frontmatter 存在
				if !strings.Contains(output, "---") {
					t.Errorf("output should contain YAML frontmatter delimiter")
				}
				// 验证标题
				if !strings.Contains(output, "title: \"Test Issue\"") {
					t.Errorf("output should contain title in frontmatter")
				}
				// 验证 URL
				if !strings.Contains(output, "url: \"https://github.com/test/repo/issues/1\"") {
					t.Errorf("output should contain URL in frontmatter")
				}
				// 验证作者
				if !strings.Contains(output, "author: \"@testuser\"") {
					t.Errorf("output should contain author in frontmatter")
				}
				// 验证创建时间
				if !strings.Contains(output, "created_at:") {
					t.Errorf("output should contain created_at in frontmatter")
				}
				// 验证状态
				if !strings.Contains(output, "status: \"open\"") {
					t.Errorf("output should contain status in frontmatter")
				}
				// 验证标题行
				if !strings.Contains(output, "# Test Issue") {
					t.Errorf("output should contain H1 title")
				}
				// 验证正文
				if !strings.Contains(output, "This is a test issue") {
					t.Errorf("output should contain issue body")
				}
			},
		},
		{
			name: "issue with comments",
			issue: createTestIssue("Issue with Comments", "Main body", []github.Comment{
				createTestComment("alice", "First comment", time.Date(2025, 1, 4, 11, 0, 0, 0, time.Local), nil),
				createTestComment("bob", "Second comment", time.Date(2025, 1, 4, 12, 0, 0, 0, time.Local), nil),
			}),
			validate: func(t *testing.T, output string) {
				// 验证评论部分存在
				if !strings.Contains(output, "## 评论") {
					t.Errorf("output should contain comments section")
				}
				// 验证两个评论都存在
				if !strings.Contains(output, "### @alice") {
					t.Errorf("output should contain first comment author")
				}
				if !strings.Contains(output, "First comment") {
					t.Errorf("output should contain first comment body")
				}
				if !strings.Contains(output, "### @bob") {
					t.Errorf("output should contain second comment author")
				}
				if !strings.Contains(output, "Second comment") {
					t.Errorf("output should contain second comment body")
				}
				// 验证评论按时间正序排列
				firstIdx := strings.Index(output, "First comment")
				secondIdx := strings.Index(output, "Second comment")
				if firstIdx >= secondIdx {
					t.Errorf("comments should be in chronological order")
				}
			},
		},
		{
			name: "with user links enabled",
			issue: createTestIssue("Test", "Body", []github.Comment{
				createTestComment("user1", "Comment", time.Now(), nil),
			}),
			options: []Option{WithUserLinks(true)},
			validate: func(t *testing.T, output string) {
				// 验证用户名格式为链接
				if !strings.Contains(output, "[@user1](https://github.com/user1)") {
					t.Errorf("output should contain user links when enabled")
				}
				// 不应该包含纯 @username 格式
				if strings.Contains(output, "**作者**: @user1") {
					t.Errorf("output should not contain plain @username when links enabled")
				}
			},
		},
		{
			name: "with reactions enabled",
			issue: createTestIssue("Test", "Body", []github.Comment{
				createTestComment("user1", "Comment", time.Now(), []github.Reaction{
					{Content: "+1", Count: 3},
					{Content: "heart", Count: 1},
					{Content: "laugh", Count: 2},
				}),
			}),
			options: []Option{WithReactions(true)},
			validate: func(t *testing.T, output string) {
				// 验证 reactions 存在
				if !strings.Contains(output, "👍 3") {
					t.Errorf("output should contain +1 reactions")
				}
				if !strings.Contains(output, "❤️ 1") {
					t.Errorf("output should contain heart reactions")
				}
				if !strings.Contains(output, "😄 2") {
					t.Errorf("output should contain laugh reactions")
				}
				// reactions 应该在评论末尾
				commentIdx := strings.Index(output, "Comment")
				reactionIdx := strings.Index(output, "👍 3")
				if reactionIdx <= commentIdx {
					t.Errorf("reactions should be after comment body")
				}
			},
		},
		{
			name:  "emoji shortcode conversion",
			issue: createTestIssue("Test Emoji", "This has :thumbsup: and :heart:", nil),
			validate: func(t *testing.T, output string) {
				// 验证 emoji shortcode 被转换
				if !strings.Contains(output, "👍") {
					t.Errorf("output should convert :thumbsup: to emoji")
				}
				if !strings.Contains(output, "❤️") {
					t.Errorf("output should convert :heart: to emoji")
				}
				// 不应该包含原始 shortcode
				if strings.Contains(output, ":thumbsup:") {
					t.Errorf("output should not contain original shortcode")
				}
				if strings.Contains(output, ":heart:") {
					t.Errorf("output should not contain original shortcode")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter(tt.options...)
			output, err := c.ConvertIssue(tt.issue)

			if err != nil {
				t.Fatalf("ConvertIssue() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestConvertIssue_DeletedComment 测试已删除评论的处理
func TestConvertIssue_DeletedComment(t *testing.T) {
	issue := createTestIssue("Test", "Body", []github.Comment{
		createTestComment("user1", "Normal comment", time.Now(), nil),
		{
			ID:        2,
			User:      github.User{Login: "deleted"},
			CreatedAt: time.Now(),
			Body:      "",
			Deleted:   true,
		},
		createTestComment("user2", "Another comment", time.Now(), nil),
	})

	c := NewConverter()
	output, err := c.ConvertIssue(issue)

	if err != nil {
		t.Fatalf("ConvertIssue() error = %v", err)
	}

	// 验证已删除评论显示为 ~~deleted~~
	if !strings.Contains(output, "~~deleted~~") {
		t.Errorf("deleted comment should show as ~~deleted~~")
	}

	// 验证其他评论仍然存在
	if !strings.Contains(output, "Normal comment") {
		t.Errorf("normal comments should still be present")
	}
	if !strings.Contains(output, "Another comment") {
		t.Errorf("other comments should still be present")
	}
}

// TestConvertPullRequest 测试 PR 转换功能
func TestConvertPullRequest(t *testing.T) {
	tests := []struct {
		name     string
		pr       *github.PullRequest
		validate func(t *testing.T, output string)
	}{
		{
			name: "basic PR conversion",
			pr: &github.PullRequest{
				Title: "Test PR",
				URL:   "https://github.com/test/repo/pull/1",
				User: github.User{
					Login:   "pruser",
					HTMLURL: "https://github.com/pruser",
				},
				CreatedAt: time.Date(2025, 1, 4, 10, 0, 0, 0, time.Local),
				State:     "open",
				Body:      "PR description",
				Comments:  nil,
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "title: \"Test PR\"") {
					t.Errorf("output should contain PR title")
				}
				if !strings.Contains(output, "status: \"open\"") {
					t.Errorf("output should contain PR state")
				}
				if !strings.Contains(output, "PR description") {
					t.Errorf("output should contain PR body")
				}
			},
		},
		{
			name: "merged PR",
			pr: &github.PullRequest{
				Title:     "Merged PR",
				URL:       "https://github.com/test/repo/pull/2",
				User:      github.User{Login: "dev", HTMLURL: "https://github.com/dev"},
				CreatedAt: time.Now(),
				State:     "merged",
				Body:      "This was merged",
				Comments:  nil,
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "status: \"merged\"") {
					t.Errorf("output should show merged status")
				}
				if !strings.Contains(output, "**状态**: Merged") {
					t.Errorf("output should display Merged status in body")
				}
			},
		},
		{
			name: "PR with review comments",
			pr: &github.PullRequest{
				Title:     "PR with Reviews",
				URL:       "https://github.com/test/repo/pull/3",
				User:      github.User{Login: "author", HTMLURL: "https://github.com/author"},
				CreatedAt: time.Now(),
				State:     "open",
				Body:      "Description",
				Comments: []github.Comment{
					createTestComment("reviewer1", "LGTM", time.Now(), nil),
					createTestComment("reviewer2", "Needs changes", time.Now(), nil),
				},
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "LGTM") {
					t.Errorf("output should contain review comments")
				}
				if !strings.Contains(output, "Needs changes") {
					t.Errorf("output should contain all review comments")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter()
			output, err := c.ConvertPullRequest(tt.pr)

			if err != nil {
				t.Fatalf("ConvertPullRequest() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestConvertDiscussion 测试 Discussion 转换功能
func TestConvertDiscussion(t *testing.T) {
	tests := []struct {
		name       string
		discussion *github.Discussion
		validate   func(t *testing.T, output string)
	}{
		{
			name: "basic discussion conversion",
			discussion: &github.Discussion{
				Title: "Test Discussion",
				URL:   "https://github.com/test/repo/discussions/1",
				User: github.User{
					Login:   "discussuser",
					HTMLURL: "https://github.com/discussuser",
				},
				CreatedAt: time.Date(2025, 1, 4, 9, 0, 0, 0, time.Local),
				State:     "open",
				Body:      "Discussion body",
				Comments:  nil,
			},
			validate: func(t *testing.T, output string) {
				if !strings.Contains(output, "title: \"Test Discussion\"") {
					t.Errorf("output should contain discussion title")
				}
				if !strings.Contains(output, "Discussion body") {
					t.Errorf("output should contain discussion body")
				}
			},
		},
		{
			name: "discussion with replies (flattened)",
			discussion: &github.Discussion{
				Title:     "Discussion with Replies",
				URL:       "https://github.com/test/repo/discussions/2",
				User:      github.User{Login: "op", HTMLURL: "https://github.com/op"},
				CreatedAt: time.Now(),
				State:     "open",
				Body:      "Main post",
				Comments: []github.Comment{
					createTestComment("reply1", "First reply", time.Date(2025, 1, 4, 10, 0, 0, 0, time.Local), nil),
					createTestComment("reply2", "Second reply", time.Date(2025, 1, 4, 11, 0, 0, 0, time.Local), nil),
					createTestComment("reply3", "Third reply", time.Date(2025, 1, 4, 12, 0, 0, 0, time.Local), nil),
				},
			},
			validate: func(t *testing.T, output string) {
				// 验证所有回复都存在
				if !strings.Contains(output, "First reply") {
					t.Errorf("output should contain first reply")
				}
				if !strings.Contains(output, "Second reply") {
					t.Errorf("output should contain second reply")
				}
				if !strings.Contains(output, "Third reply") {
					t.Errorf("output should contain third reply")
				}
				// 验证是平铺展示（没有缩进层级）
				// Discussion 的所有评论应该按时间顺序平铺显示
				firstIdx := strings.Index(output, "First reply")
				secondIdx := strings.Index(output, "Second reply")
				thirdIdx := strings.Index(output, "Third reply")

				if firstIdx >= secondIdx || secondIdx >= thirdIdx {
					t.Errorf("replies should be in chronological order")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConverter()
			output, err := c.ConvertDiscussion(tt.discussion)

			if err != nil {
				t.Fatalf("ConvertDiscussion() error = %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, output)
			}
		})
	}
}

// TestFormatTimestamp 测试时间格式化
func TestFormatTimestamp(t *testing.T) {
	testTime := time.Date(2025, 1, 4, 14, 30, 45, 0, time.Local)
	expected := "2025-01-04 14:30:45"

	c := NewConverter()
	output := c.formatTimestamp(testTime)
	if output != expected {
		t.Errorf("formatTimestamp() = %v, want %v", output, expected)
	}
}

// TestConvertEmojiShortcode 测试 Emoji 转换
func TestConvertEmojiShortcode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{":thumbsup:", "👍"},
		{":heart:", "❤️"},
		{":+1:", "👍"},
		{":-1:", "👎"},
		{":smile:", "😄"},
		{":laugh:", "😄"},
		{":tada:", "🎉"},
		{":hooray:", "🎉"},
		{":confused:", "😕"},
		{":rocket:", "🚀"},
		{":eyes:", "👀"},
		{"Mixed :thumbsup: and :heart: text", "Mixed 👍 and ❤️ text"},
		{"No emoji here", "No emoji here"},
		{"Multiple :thumbsup: :thumbsup: :heart:", "Multiple 👍 👍 ❤️"},
	}

	c := NewConverter()
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			output := c.convertEmojiShortcode(tt.input)
			if output != tt.expected {
				t.Errorf("convertEmojiShortcode(%q) = %q, want %q", tt.input, output, tt.expected)
			}
		})
	}
}
