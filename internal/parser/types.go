package parser

// ResourceType GitHub资源类型
type ResourceType int

const (
	Unknown ResourceType = iota
	Issue
	PullRequest
	Discussion
)

// String 实现Stringer接口
func (rt ResourceType) String() string {
	switch rt {
	case Issue:
		return "issue"
	case PullRequest:
		return "pull_request"
	case Discussion:
		return "discussion"
	default:
		return "unknown"
	}
}

// Resource 解析后的GitHub资源
type Resource struct {
	Type        ResourceType
	Owner       string
	Repo        string
	Number      int
	OriginalURL string
}
