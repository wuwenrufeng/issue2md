package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 错误定义
var (
	ErrResourceNotFound = fmt.Errorf("resource not found")
	ErrAPIRateLimit     = fmt.Errorf("API rate limit exceeded")
	ErrNetwork          = fmt.Errorf("network error")
)

// Option 配置 Client 的选项
type Option func(*Client)

// WithBaseURL 设置自定义的 BaseURL（用于测试）
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// Client GitHub API 客户端
type Client struct {
	baseURL string       // API Base URL
	token   string       // GitHub Token（可选）
	client  *http.Client // HTTP 客户端
}

// NewClient 创建新的 GitHub Client
func NewClient(token string, opts ...Option) *Client {
	client := &Client{
		baseURL: "https://api.github.com", // 默认 GitHub API URL
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// 应用选项
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// FetchIssue 获取 GitHub Issue
func (c *Client) FetchIssue(owner, repo string, number int) (*Issue, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/issues/%d", c.baseURL, owner, repo, number)

	// 获取 Issue 主数据
	var issueData struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login   string `json:"login"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		State     string    `json:"state"`
		Body      string    `json:"body"`
	}

	err := c.get(url, &issueData)
	if err != nil {
		return nil, err
	}

	// 构建Issue
	issue := &Issue{
		Title: issueData.Title,
		URL:   issueData.HTMLURL,
		User: User{
			Login:   issueData.User.Login,
			HTMLURL: issueData.User.HTMLURL,
		},
		CreatedAt: issueData.CreatedAt,
		State:     issueData.State,
		Body:      issueData.Body,
	}

	// 获取评论
	commentsURL := fmt.Sprintf("%s/repos/%s/%s/issues/%d/comments", c.baseURL, owner, repo, number)
	var commentsData []struct {
		ID   int64 `json:"id"`
		User struct {
			Login   string `json:"login"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		Body      string    `json:"body"`
		Reactions struct {
			TotalCount int `json:"total_count"`
			PlusOne    int `json:"plus_one"`
			Heart      int `json:"heart"`
		} `json:"reactions"`
	}

	err = c.get(commentsURL, &commentsData)
	if err == nil {
		issue.Comments = make([]Comment, len(commentsData))
		for i, cData := range commentsData {
			issue.Comments[i] = Comment{
				ID:        cData.ID,
				User:      User{Login: cData.User.Login, HTMLURL: cData.User.HTMLURL},
				CreatedAt: cData.CreatedAt,
				Body:      cData.Body,
				Reactions: buildReactions(cData.Reactions),
				Deleted:   false,
			}
		}
	}

	return issue, nil
}

// FetchPullRequest 获取 GitHub Pull Request
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d", c.baseURL, owner, repo, number)

	var prData struct {
		Title   string `json:"title"`
		HTMLURL string `json:"html_url"`
		User    struct {
			Login   string `json:"login"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		State     string    `json:"state"`
		Body      string    `json:"body"`
		Merged    bool      `json:"merged"`
	}

	err := c.get(url, &prData)
	if err != nil {
		return nil, err
	}

	// 确定 PR 状态
	state := prData.State
	if prData.Merged {
		state = "merged"
	}

	pr := &PullRequest{
		Title: prData.Title,
		URL:   prData.HTMLURL,
		User: User{
			Login:   prData.User.Login,
			HTMLURL: prData.User.HTMLURL,
		},
		CreatedAt: prData.CreatedAt,
		State:     state,
		Body:      prData.Body,
	}

	// 获取评论（Review 评论）
	commentsURL := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/comments", c.baseURL, owner, repo, number)
	var commentsData []struct {
		ID   int64 `json:"id"`
		User struct {
			Login   string `json:"login"`
			HTMLURL string `json:"html_url"`
		} `json:"user"`
		CreatedAt time.Time `json:"created_at"`
		Body      string    `json:"body"`
	}

	err = c.get(commentsURL, &commentsData)
	if err == nil {
		pr.Comments = make([]Comment, len(commentsData))
		for i, cData := range commentsData {
			pr.Comments[i] = Comment{
				ID:        cData.ID,
				User:      User{Login: cData.User.Login, HTMLURL: cData.User.HTMLURL},
				CreatedAt: cData.CreatedAt,
				Body:      cData.Body,
				Reactions: []Reaction{},
				Deleted:   false,
			}
		}
	}

	return pr, nil
}

// FetchDiscussion 获取 GitHub Discussion（使用 GraphQL）
func (c *Client) FetchDiscussion(owner, repo string, number int) (*Discussion, error) {
	// GraphQL 查询
	query := fmt.Sprintf(`{
		repository(owner: "%s", name: "%s") {
			discussion(number: %d) {
				title
				url
				author {
					login
					url
				}
				createdAt
				closedAt
				body
				comments(first: 100) {
					nodes {
						id
						author {
							login
							url
						}
						createdAt
						body
						reactions(first: 100) {
							nodes {
								content
								users {
									totalCount
								}
							}
						}
					}
				}
			}
		}
	}`, owner, repo, number)

	url := c.baseURL + "/graphql"

	var response struct {
		Data struct {
			Repository struct {
				Discussion *struct {
					Title  string `json:"title"`
					URL    string `json:"url"`
					Author struct {
						Login string `json:"login"`
						URL   string `json:"url"`
					} `json:"author"`
					CreatedAt time.Time  `json:"createdAt"`
					ClosedAt  *time.Time `json:"closedAt"`
					Body      string     `json:"body"`
					Comments  struct {
						Nodes []struct {
							ID     string `json:"id"`
							Author struct {
								Login string `json:"login"`
								URL   string `json:"url"`
							} `json:"author"`
							CreatedAt time.Time `json:"createdAt"`
							Body      string    `json:"body"`
							Reactions struct {
								Nodes []struct {
									Content string `json:"content"`
									Users   struct {
										TotalCount int `json:"totalCount"`
									} `json:"users"`
								} `json:"nodes"`
							} `json:"reactions"`
						} `json:"nodes"`
					} `json:"comments"`
				} `json:"discussion"`
			} `json:"repository"`
		} `json:"data"`
	}

	err := c.postGraphQL(url, query, &response)
	if err != nil {
		return nil, err
	}

	if response.Data.Repository.Discussion == nil {
		return nil, ErrResourceNotFound
	}

	d := response.Data.Repository.Discussion

	// 确定状态
	state := "open"
	if d.ClosedAt != nil {
		state = "closed"
	}

	// 构建评论
	comments := make([]Comment, len(d.Comments.Nodes))
	for i, node := range d.Comments.Nodes {
		reactions := make([]Reaction, len(node.Reactions.Nodes))
		for j, rNode := range node.Reactions.Nodes {
			reactions[j] = Reaction{
				Content: convertReactionContent(rNode.Content),
				Count:   rNode.Users.TotalCount,
			}
		}

		comments[i] = Comment{
			ID:        0, // GraphQL ID 是字符串，暂时设为 0
			User:      User{Login: node.Author.Login, HTMLURL: node.Author.URL},
			CreatedAt: node.CreatedAt,
			Body:      node.Body,
			Reactions: reactions,
			Deleted:   false,
		}
	}

	return &Discussion{
		Title:     d.Title,
		URL:       d.URL,
		User:      User{Login: d.Author.Login, HTMLURL: d.Author.URL},
		CreatedAt: d.CreatedAt,
		State:     state,
		Body:      d.Body,
		Comments:  comments,
	}, nil
}

// get 发送 GET 请求
func (c *Client) get(url string, v interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// 添加认证头
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNetwork, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrResourceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	return nil
}

// postGraphQL 发送 GraphQL POST 请求
func (c *Client) postGraphQL(url string, query string, v interface{}) error {
	payload := map[string]string{
		"query": query,
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrNetwork, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrResourceNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	return nil
}

// buildReactions 构建 reactions 列表
func buildReactions(r struct {
	TotalCount int `json:"total_count"`
	PlusOne    int `json:"plus_one"`
	Heart      int `json:"heart"`
}) []Reaction {
	reactions := []Reaction{}

	if r.PlusOne > 0 {
		reactions = append(reactions, Reaction{Content: "+1", Count: r.PlusOne})
	}
	if r.Heart > 0 {
		reactions = append(reactions, Reaction{Content: "heart", Count: r.Heart})
	}

	return reactions
}

// convertReactionContent 转换 GraphQL reaction content
func convertReactionContent(content string) string {
	// GraphQL 返回的是大写（如 THUMBS_UP），转换为小写
	switch content {
	case "THUMBS_UP":
		return "+1"
	case "HEART":
		return "heart"
	default:
		return content
	}
}
