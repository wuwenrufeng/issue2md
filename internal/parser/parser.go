package parser

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// 错误定义
var (
	ErrInvalidURLFormat        = errors.New("invalid URL format")
	ErrUnsupportedResourceType = errors.New("unsupported resource type")
)

// parseAndValidateURL 解析URL并验证基础格式
// 返回路径分段和解析后的URL对象
func parseAndValidateURL(rawURL string) ([]string, *url.URL, error) {
	// 处理空字符串
	if rawURL == "" {
		return nil, nil, ErrInvalidURLFormat
	}

	// 使用标准库解析URL（自动处理query和fragment）
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse URL %q failed: %w", rawURL, ErrInvalidURLFormat)
	}

	// 验证host必须是github.com
	if parsed.Host != "github.com" {
		return nil, nil, fmt.Errorf("host must be github.com, got %q: %w", parsed.Host, ErrInvalidURLFormat)
	}

	// 分割路径，去掉开头的空字符串
	path := strings.Trim(parsed.Path, "/")
	if path == "" {
		return nil, nil, fmt.Errorf("empty URL path: %w", ErrInvalidURLFormat)
	}

	parts := strings.Split(path, "/")

	// 验证路径段数量并区分错误类型
	if len(parts) < 2 {
		return nil, nil, fmt.Errorf("invalid path segments count %d: %w", len(parts), ErrInvalidURLFormat)
	}

	// 2段：可能是 /owner/repo 或 /issues/123 等情况
	if len(parts) == 2 {
		firstPart := strings.ToLower(parts[0])
		githubReserved := map[string]bool{
			"issues": true, "pull": true, "discussions": true,
			"settings": true, "actions": true, "wiki": true,
		}
		if githubReserved[firstPart] {
			return nil, nil, fmt.Errorf("invalid path format: %w", ErrInvalidURLFormat)
		}
		return nil, nil, ErrUnsupportedResourceType
	}

	// 3段：需要判断是unsupported还是invalid
	if len(parts) == 3 {
		thirdPart := strings.ToLower(parts[2])
		unsupportedTypes := map[string]bool{
			"wiki": true, "actions": true, "security": true, "pulse": true,
		}
		if unsupportedTypes[thirdPart] {
			return nil, nil, fmt.Errorf("resource type %q: %w", thirdPart, ErrUnsupportedResourceType)
		}
		return nil, nil, fmt.Errorf("incomplete URL path: %w", ErrInvalidURLFormat)
	}

	// 4段及以上：验证通过
	return parts, parsed, nil
}

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
func ParseURL(rawURL string) (*Resource, error) {
	parts, parsed, err := parseAndValidateURL(rawURL)
	if err != nil {
		return nil, err
	}

	// 提取路径信息
	owner := parts[0]
	repo := parts[1]
	resourceType := strings.ToLower(parts[2])
	numberStr := parts[3]

	// 验证owner和repo非空
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner or repo is empty: %w", ErrInvalidURLFormat)
	}

	// 解析number
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return nil, fmt.Errorf("cannot parse number %q: %w", numberStr, ErrInvalidURLFormat)
	}

	// 识别资源类型
	var resType ResourceType
	switch resourceType {
	case "issues":
		resType = Issue
	case "pull":
		resType = PullRequest
	case "discussions":
		resType = Discussion
	default:
		return nil, fmt.Errorf("unsupported resource type %q: %w", resourceType, ErrUnsupportedResourceType)
	}

	// 构建干净的URL（去掉query和fragment，保留原始大小写）
	cleanURL := fmt.Sprintf("https://%s/%s/%s/%s/%d",
		parsed.Host, owner, repo, parts[2], number)

	return &Resource{
		Type:        resType,
		Owner:       owner,
		Repo:        repo,
		Number:      number,
		OriginalURL: cleanURL,
	}, nil
}
