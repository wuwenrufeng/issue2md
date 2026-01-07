package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestFetchIssue_Success 测试成功获取 Issue
func TestFetchIssue_Success(t *testing.T) {
	// 创建 Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 处理 Issue 请求或 comments 请求
		if r.URL.Path == "/repos/test-owner/test-repo/issues/1" {
			// 返回模拟的 Issue 数据
			mockIssue := map[string]interface{}{
				"title":    "Test Issue",
				"html_url": "https://github.com/test-owner/test-repo/issues/1",
				"user": map[string]interface{}{
					"login":    "testuser",
					"html_url": "https://github.com/testuser",
				},
				"created_at": "2024-01-01T00:00:00Z",
				"state":      "open",
				"body":       "This is a test issue",
				"comments":   0,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockIssue)
		} else if r.URL.Path == "/repos/test-owner/test-repo/issues/1/comments" {
			// 返回空评论列表
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]interface{}{})
		}
	}))
	defer mockServer.Close()

	// 创建 Client（使用 Mock Server 的地址）
	client := NewClient("", WithBaseURL(mockServer.URL))

	// 调用 FetchIssue
	issue, err := client.FetchIssue("test-owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("FetchIssue failed: %v", err)
	}

	// 验证返回的数据
	if issue.Title != "Test Issue" {
		t.Errorf("expected title 'Test Issue', got '%s'", issue.Title)
	}

	if issue.URL != "https://github.com/test-owner/test-repo/issues/1" {
		t.Errorf("expected URL 'https://github.com/test-owner/test-repo/issues/1', got '%s'", issue.URL)
	}

	if issue.User.Login != "testuser" {
		t.Errorf("expected user login 'testuser', got '%s'", issue.User.Login)
	}

	if issue.State != "open" {
		t.Errorf("expected state 'open', got '%s'", issue.State)
	}

	if issue.Body != "This is a test issue" {
		t.Errorf("expected body 'This is a test issue', got '%s'", issue.Body)
	}
}

// TestFetchIssue_WithComments 测试获取带评论的 Issue
func TestFetchIssue_WithComments(t *testing.T) {
	// 用于存储评论数据的变量
	var commentsRequest bool

	// 创建 Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test-owner/test-repo/issues/1" {
			// Issue 请求
			mockIssue := map[string]interface{}{
				"title":    "Test Issue with Comments",
				"html_url": "https://github.com/test-owner/test-repo/issues/1",
				"user": map[string]interface{}{
					"login":    "testuser",
					"html_url": "https://github.com/testuser",
				},
				"created_at": "2024-01-01T00:00:00Z",
				"state":      "open",
				"body":       "This is a test issue",
				"comments":   2,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockIssue)
		} else if r.URL.Path == "/repos/test-owner/test-repo/issues/1/comments" {
			// 评论请求
			commentsRequest = true
			mockComments := []map[string]interface{}{
				{
					"id": 1001,
					"user": map[string]interface{}{
						"login":    "commenter1",
						"html_url": "https://github.com/commenter1",
					},
					"created_at": "2024-01-02T00:00:00Z",
					"body":       "First comment",
					"reactions": map[string]interface{}{
						"total_count": 5,
						"plus_one":    3,
						"heart":       2,
					},
				},
				{
					"id": 1002,
					"user": map[string]interface{}{
						"login":    "commenter2",
						"html_url": "https://github.com/commenter2",
					},
					"created_at": "2024-01-03T00:00:00Z",
					"body":       "Second comment",
					"reactions": map[string]interface{}{
						"total_count": 1,
						"plus_one":    1,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockComments)
		}
	}))
	defer mockServer.Close()

	client := NewClient("", WithBaseURL(mockServer.URL))

	issue, err := client.FetchIssue("test-owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("FetchIssue failed: %v", err)
	}

	// 验证请求了评论接口
	if !commentsRequest {
		t.Error("expected to request comments endpoint")
	}

	// 验证评论数量
	if len(issue.Comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(issue.Comments))
	}

	// 验证第一条评论
	if issue.Comments[0].ID != 1001 {
		t.Errorf("expected first comment ID 1001, got %d", issue.Comments[0].ID)
	}

	if issue.Comments[0].User.Login != "commenter1" {
		t.Errorf("expected first comment user 'commenter1', got '%s'", issue.Comments[0].User.Login)
	}

	if issue.Comments[0].Body != "First comment" {
		t.Errorf("expected first comment body 'First comment', got '%s'", issue.Comments[0].Body)
	}

	// 验证 reactions
	if len(issue.Comments[0].Reactions) != 2 {
		t.Errorf("expected 2 reaction types, got %d", len(issue.Comments[0].Reactions))
	}
}

// TestFetchIssue_NotFound 测试 Issue 不存在的情况
func TestFetchIssue_NotFound(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		mockError := map[string]interface{}{
			"message": "Not Found",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockError)
	}))
	defer mockServer.Close()

	client := NewClient("", WithBaseURL(mockServer.URL))

	_, err := client.FetchIssue("test-owner", "test-repo", 999)
	if err == nil {
		t.Fatal("expected error for non-existent issue, got nil")
	}

	// 验证错误类型
	if err != ErrResourceNotFound {
		t.Errorf("expected ErrResourceNotFound, got %v", err)
	}
}

// TestFetchPullRequest_Success 测试成功获取 Pull Request
func TestFetchPullRequest_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/repos/test-owner/test-repo/pulls/1" {
			mockPR := map[string]interface{}{
				"title":    "Test PR",
				"html_url": "https://github.com/test-owner/test-repo/pull/1",
				"user": map[string]interface{}{
					"login":    "prauthor",
					"html_url": "https://github.com/prauthor",
				},
				"created_at": "2024-01-01T00:00:00Z",
				"state":      "open",
				"body":       "This is a test PR",
				"merged":     false,
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockPR)
		} else if r.URL.Path == "/repos/test-owner/test-repo/pulls/1/comments" {
			mockComments := []map[string]interface{}{
				{
					"id": 2001,
					"user": map[string]interface{}{
						"login":    "reviewer1",
						"html_url": "https://github.com/reviewer1",
					},
					"created_at": "2024-01-02T00:00:00Z",
					"body":       "Review comment",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockComments)
		}
	}))
	defer mockServer.Close()

	client := NewClient("", WithBaseURL(mockServer.URL))

	pr, err := client.FetchPullRequest("test-owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("FetchPullRequest failed: %v", err)
	}

	if pr.Title != "Test PR" {
		t.Errorf("expected title 'Test PR', got '%s'", pr.Title)
	}

	if pr.State != "open" {
		t.Errorf("expected state 'open', got '%s'", pr.State)
	}

	if len(pr.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(pr.Comments))
	}
}

// TestFetchPullRequest_Merged 测试获取已合并的 PR
func TestFetchPullRequest_Merged(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mockPR := map[string]interface{}{
			"title":    "Merged PR",
			"html_url": "https://github.com/test-owner/test-repo/pull/2",
			"user": map[string]interface{}{
				"login":    "prauthor",
				"html_url": "https://github.com/prauthor",
			},
			"created_at": "2024-01-01T00:00:00Z",
			"state":      "closed",
			"body":       "This PR was merged",
			"merged":     true,
			"merged_at":  "2024-01-02T00:00:00Z",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockPR)
	}))
	defer mockServer.Close()

	client := NewClient("", WithBaseURL(mockServer.URL))

	pr, err := client.FetchPullRequest("test-owner", "test-repo", 2)
	if err != nil {
		t.Fatalf("FetchPullRequest failed: %v", err)
	}

	// 验证 merged 状态被正确处理
	if pr.State != "merged" {
		t.Errorf("expected state 'merged', got '%s'", pr.State)
	}
}

// TestFetchDiscussion_Success 测试成功获取 Discussion
func TestFetchDiscussion_Success(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Discussion 使用 GraphQL API
		// 验证请求是 GraphQL
		if r.Method != "POST" {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// 返回模拟的 GraphQL 响应
		mockResponse := map[string]interface{}{
			"data": map[string]interface{}{
				"repository": map[string]interface{}{
					"discussion": map[string]interface{}{
						"title": "Test Discussion",
						"url":   "https://github.com/test-owner/test-repo/discussions/1",
						"author": map[string]interface{}{
							"login": "discussant",
							"url":   "https://github.com/discussant",
						},
						"createdAt": "2024-01-01T00:00:00Z",
						"closedAt":  nil,
						"body":      "This is a test discussion",
						"comments": map[string]interface{}{
							"nodes": []map[string]interface{}{
								{
									"id": "test-id-1",
									"author": map[string]interface{}{
										"login": "commenter1",
										"url":   "https://github.com/commenter1",
									},
									"createdAt": "2024-01-02T00:00:00Z",
									"body":      "Discussion comment",
									"reactions": map[string]interface{}{
										"nodes": []map[string]interface{}{
											{
												"content": "THUMBS_UP",
												"users": map[string]interface{}{
													"totalCount": 3,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client := NewClient("", WithBaseURL(mockServer.URL))

	discussion, err := client.FetchDiscussion("test-owner", "test-repo", 1)
	if err != nil {
		t.Fatalf("FetchDiscussion failed: %v", err)
	}

	if discussion.Title != "Test Discussion" {
		t.Errorf("expected title 'Test Discussion', got '%s'", discussion.Title)
	}

	if discussion.State != "open" {
		t.Errorf("expected state 'open', got '%s'", discussion.State)
	}

	if len(discussion.Comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(discussion.Comments))
	}
}

// TestNewClient 测试创建 Client
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "client with token",
			token:   "test-token",
			wantErr: false,
		},
		{
			name:    "client without token",
			token:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			if client == nil {
				t.Error("expected non-nil client")
			}
		})
	}
}

// 辅助函数：解析时间字符串
func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}
