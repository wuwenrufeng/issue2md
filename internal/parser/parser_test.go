package parser

import (
	"errors"
	"testing"
)

// TestParseURL 表格驱动测试：验证ParseURL函数的各种场景
func TestParseURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    *Resource
		wantErr error
	}{
		// ===== 有效的 Issue URL =====
		{
			name: "valid issue URL - basic",
			url:  "https://github.com/owner/repo/issues/123",
			want: &Resource{
				Type:        Issue,
				Owner:       "owner",
				Repo:        "repo",
				Number:      123,
				OriginalURL: "https://github.com/owner/repo/issues/123",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - golang/go",
			url:  "https://github.com/golang/go/issues/1",
			want: &Resource{
				Type:        Issue,
				Owner:       "golang",
				Repo:        "go",
				Number:      1,
				OriginalURL: "https://github.com/golang/go/issues/1",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - with hyphen in owner",
			url:  "https://github.com/my-org/my-repo/issues/42",
			want: &Resource{
				Type:        Issue,
				Owner:       "my-org",
				Repo:        "my-repo",
				Number:      42,
				OriginalURL: "https://github.com/my-org/my-repo/issues/42",
			},
			wantErr: nil,
		},

		// ===== 有效的 Pull Request URL =====
		{
			name: "valid PR URL - basic",
			url:  "https://github.com/owner/repo/pull/456",
			want: &Resource{
				Type:        PullRequest,
				Owner:       "owner",
				Repo:        "repo",
				Number:      456,
				OriginalURL: "https://github.com/owner/repo/pull/456",
			},
			wantErr: nil,
		},
		{
			name: "valid PR URL - large number",
			url:  "https://github.com/golang/go/pull/12345",
			want: &Resource{
				Type:        PullRequest,
				Owner:       "golang",
				Repo:        "go",
				Number:      12345,
				OriginalURL: "https://github.com/golang/go/pull/12345",
			},
			wantErr: nil,
		},

		// ===== 有效的 Discussion URL =====
		{
			name: "valid discussion URL - basic",
			url:  "https://github.com/owner/repo/discussions/789",
			want: &Resource{
				Type:        Discussion,
				Owner:       "owner",
				Repo:        "repo",
				Number:      789,
				OriginalURL: "https://github.com/owner/repo/discussions/789",
			},
			wantErr: nil,
		},
		{
			name: "valid discussion URL - github/discussions",
			url:  "https://github.com/github/discussions/discussions/1",
			want: &Resource{
				Type:        Discussion,
				Owner:       "github",
				Repo:        "discussions",
				Number:      1,
				OriginalURL: "https://github.com/github/discussions/discussions/1",
			},
			wantErr: nil,
		},

		// ===== 无效的 URL（格式错误） =====
		{
			name:    "invalid URL - not github.com",
			url:     "https://gitlab.com/owner/repo/issues/123",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - missing owner",
			url:     "https://github.com/issues/123",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - missing repo",
			url:     "https://github.com/owner/issues/123",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - missing resource type",
			url:     "https://github.com/owner/repo/123",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - missing number",
			url:     "https://github.com/owner/repo/issues",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - number is not numeric",
			url:     "https://github.com/owner/repo/issues/abc",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - empty string",
			url:     "",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},
		{
			name:    "invalid URL - random website",
			url:     "https://example.com",
			want:    nil,
			wantErr: ErrInvalidURLFormat,
		},

		// ===== 不支持的 URL 类型 =====
		{
			name:    "unsupported URL - repository home",
			url:     "https://github.com/owner/repo",
			want:    nil,
			wantErr: ErrUnsupportedResourceType,
		},
		{
			name:    "unsupported URL - actions",
			url:     "https://github.com/owner/repo/actions",
			want:    nil,
			wantErr: ErrUnsupportedResourceType,
		},
		{
			name:    "unsupported URL - wiki",
			url:     "https://github.com/owner/repo/wiki",
			want:    nil,
			wantErr: ErrUnsupportedResourceType,
		},
		{
			name:    "unsupported URL - security",
			url:     "https://github.com/owner/repo/security",
			want:    nil,
			wantErr: ErrUnsupportedResourceType,
		},

		// ===== 边界条件测试 =====
		{
			name: "valid issue URL - with query parameter (should ignore)",
			url:  "https://github.com/owner/repo/issues/123?q=1",
			want: &Resource{
				Type:        Issue,
				Owner:       "owner",
				Repo:        "repo",
				Number:      123,
				OriginalURL: "https://github.com/owner/repo/issues/123",
			},
			wantErr: nil,
		},
		{
			name: "valid PR URL - with fragment (should ignore)",
			url:  "https://github.com/owner/repo/pull/456#issuecomment-123456",
			want: &Resource{
				Type:        PullRequest,
				Owner:       "owner",
				Repo:        "repo",
				Number:      456,
				OriginalURL: "https://github.com/owner/repo/pull/456",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - with both query and fragment",
			url:  "https://github.com/owner/repo/issues/789?q=1#top",
			want: &Resource{
				Type:        Issue,
				Owner:       "owner",
				Repo:        "repo",
				Number:      789,
				OriginalURL: "https://github.com/owner/repo/issues/789",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - case insensitive (should work)",
			url:  "https://github.com/Owner/Repo/ISSUES/999",
			want: &Resource{
				Type:        Issue,
				Owner:       "Owner",
				Repo:        "Repo",
				Number:      999,
				OriginalURL: "https://github.com/Owner/Repo/ISSUES/999",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - number is 1",
			url:  "https://github.com/owner/repo/issues/1",
			want: &Resource{
				Type:        Issue,
				Owner:       "owner",
				Repo:        "repo",
				Number:      1,
				OriginalURL: "https://github.com/owner/repo/issues/1",
			},
			wantErr: nil,
		},
		{
			name: "valid issue URL - number is 0 (edge case)",
			url:  "https://github.com/owner/repo/issues/0",
			want: &Resource{
				Type:        Issue,
				Owner:       "owner",
				Repo:        "repo",
				Number:      0,
				OriginalURL: "https://github.com/owner/repo/issues/0",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURL(tt.url)

			// 验证错误（使用 errors.Is 以支持错误包装）
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// 验证返回值
			if tt.want == nil {
				if got != nil {
					t.Errorf("ParseURL() = %v, want nil", got)
				}
				return
			}

			// 深度比较 Resource 结构体
			if got == nil {
				t.Errorf("ParseURL() = nil, want %v", tt.want)
				return
			}

			if got.Type != tt.want.Type {
				t.Errorf("ParseURL() Type = %v, want %v", got.Type, tt.want.Type)
			}
			if got.Owner != tt.want.Owner {
				t.Errorf("ParseURL() Owner = %v, want %v", got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("ParseURL() Repo = %v, want %v", got.Repo, tt.want.Repo)
			}
			if got.Number != tt.want.Number {
				t.Errorf("ParseURL() Number = %v, want %v", got.Number, tt.want.Number)
			}
			if got.OriginalURL != tt.want.OriginalURL {
				t.Errorf("ParseURL() OriginalURL = %v, want %v", got.OriginalURL, tt.want.OriginalURL)
			}
		})
	}
}
