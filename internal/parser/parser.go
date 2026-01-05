package parser

import (
	"errors"
)

// 错误定义
var (
	ErrInvalidURLFormat        = errors.New("invalid GitHub URL format")
	ErrUnsupportedResourceType = errors.New("unsupported resource type")
)

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
func ParseURL(url string) (*Resource, error) {
	// TODO: 待实现
	return nil, ErrInvalidURLFormat
}
