package github

import "time"

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
