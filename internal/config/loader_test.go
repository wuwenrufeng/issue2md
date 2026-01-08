package config

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestLoadFromFlags_ValidURL 测试基本的URL参数解析
func TestLoadFromFlags_ValidURL(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	cfg, exitCode := LoadFromFlags([]string{"https://github.com/owner/repo/issues/123"}, stdout, stderr)

	if exitCode != -1 {
		t.Errorf("expected exitCode -1, got %d", exitCode)
	}

	if cfg.URL != "https://github.com/owner/repo/issues/123" {
		t.Errorf("expected URL to be 'https://github.com/owner/repo/issues/123', got '%s'", cfg.URL)
	}

	if cfg.OutputFile != "" {
		t.Errorf("expected OutputFile to be empty, got '%s'", cfg.OutputFile)
	}

	if cfg.EnableReactions {
		t.Error("expected EnableReactions to be false")
	}

	if cfg.EnableUserLinks {
		t.Error("expected EnableUserLinks to be false")
	}
}

// TestLoadFromFlags_EnableReactionsFlag 测试 --enable-reactions flag
func TestLoadFromFlags_EnableReactionsFlag(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		wantEnableReactions bool
	}{
		{
			name:                "enable reactions flag set",
			args:                []string{"-enable-reactions", "https://github.com/owner/repo/issues/123"},
			wantEnableReactions: true,
		},
		{
			name:                "enable reactions flag not set",
			args:                []string{"https://github.com/owner/repo/issues/123"},
			wantEnableReactions: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cfg, exitCode := LoadFromFlags(tt.args, stdout, stderr)

			if exitCode != -1 {
				t.Errorf("expected exitCode -1, got %d", exitCode)
			}

			if cfg.EnableReactions != tt.wantEnableReactions {
				t.Errorf("expected EnableReactions to be %v, got %v", tt.wantEnableReactions, cfg.EnableReactions)
			}
		})
	}
}

// TestLoadFromFlags_EnableUserLinksFlag 测试 --enable-user-links flag
func TestLoadFromFlags_EnableUserLinksFlag(t *testing.T) {
	tests := []struct {
		name                string
		args                []string
		wantEnableUserLinks bool
	}{
		{
			name:                "enable user links flag set",
			args:                []string{"-enable-user-links", "https://github.com/owner/repo/issues/123"},
			wantEnableUserLinks: true,
		},
		{
			name:                "enable user links flag not set",
			args:                []string{"https://github.com/owner/repo/issues/123"},
			wantEnableUserLinks: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cfg, exitCode := LoadFromFlags(tt.args, stdout, stderr)

			if exitCode != -1 {
				t.Errorf("expected exitCode -1, got %d", exitCode)
			}

			if cfg.EnableUserLinks != tt.wantEnableUserLinks {
				t.Errorf("expected EnableUserLinks to be %v, got %v", tt.wantEnableUserLinks, cfg.EnableUserLinks)
			}
		})
	}
}

// TestLoadFromFlags_BothFlags 测试同时启用两个flag
func TestLoadFromFlags_BothFlags(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{"-enable-reactions", "-enable-user-links", "https://github.com/owner/repo/issues/123"}
	cfg, exitCode := LoadFromFlags(args, stdout, stderr)

	if exitCode != -1 {
		t.Errorf("expected exitCode -1, got %d", exitCode)
	}

	if !cfg.EnableReactions {
		t.Error("expected EnableReactions to be true")
	}

	if !cfg.EnableUserLinks {
		t.Error("expected EnableUserLinks to be true")
	}
}

// TestLoadFromFlags_GitHubToken 测试 GITHUB_TOKEN 环境变量
func TestLoadFromFlags_GitHubToken(t *testing.T) {
	// 保存原始环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	defer os.Setenv("GITHUB_TOKEN", originalToken)

	tests := []struct {
		name      string
		envToken  string
		wantToken string
	}{
		{
			name:      "token set in environment",
			envToken:  "ghp_test_token_123",
			wantToken: "ghp_test_token_123",
		},
		{
			name:      "token not set in environment",
			envToken:  "",
			wantToken: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			if tt.envToken != "" {
				os.Setenv("GITHUB_TOKEN", tt.envToken)
			} else {
				os.Unsetenv("GITHUB_TOKEN")
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cfg, exitCode := LoadFromFlags([]string{"https://github.com/owner/repo/issues/123"}, stdout, stderr)

			if exitCode != -1 {
				t.Errorf("expected exitCode -1, got %d", exitCode)
			}

			if cfg.Token != tt.wantToken {
				t.Errorf("expected Token to be '%s', got '%s'", tt.wantToken, cfg.Token)
			}
		})
	}
}

// TestLoadFromFlags_OutputFile 测试输出文件位置参数
func TestLoadFromFlags_OutputFile(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		wantOutputFile string
	}{
		{
			name:           "with output file",
			args:           []string{"https://github.com/owner/repo/issues/123", "output.md"},
			wantOutputFile: "output.md",
		},
		{
			name:           "without output file",
			args:           []string{"https://github.com/owner/repo/issues/123"},
			wantOutputFile: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			cfg, exitCode := LoadFromFlags(tt.args, stdout, stderr)

			if exitCode != -1 {
				t.Errorf("expected exitCode -1, got %d", exitCode)
			}

			if cfg.OutputFile != tt.wantOutputFile {
				t.Errorf("expected OutputFile to be '%s', got '%s'", tt.wantOutputFile, cfg.OutputFile)
			}
		})
	}
}

// TestLoadFromFlags_HelpFlag 测试 --help flag
func TestLoadFromFlags_HelpFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{"--help"}
	cfg, exitCode := LoadFromFlags(args, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("expected exitCode 0 for --help, got %d", exitCode)
	}

	if cfg != nil {
		t.Error("expected config to be nil for --help")
	}

	output := stdout.String()
	if !strings.Contains(output, "Usage:") {
		t.Error("expected help output to contain 'Usage:'")
	}

	if !strings.Contains(output, "issue2md") {
		t.Error("expected help output to contain 'issue2md'")
	}
}

// TestLoadFromFlags_VersionFlag 测试 --version flag
func TestLoadFromFlags_VersionFlag(t *testing.T) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{"--version"}
	cfg, exitCode := LoadFromFlags(args, stdout, stderr)

	if exitCode != 0 {
		t.Errorf("expected exitCode 0 for --version, got %d", exitCode)
	}

	if cfg != nil {
		t.Error("expected config to be nil for --version")
	}

	output := stdout.String()
	if !strings.Contains(output, "version:") {
		t.Error("expected version output to contain 'version:'")
	}
}

// TestLoadFromFlags_AllOptions 测试所有选项的组合
func TestLoadFromFlags_AllOptions(t *testing.T) {
	// 设置环境变量
	originalToken := os.Getenv("GITHUB_TOKEN")
	os.Setenv("GITHUB_TOKEN", "ghp_test_token")
	defer os.Setenv("GITHUB_TOKEN", originalToken)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	args := []string{
		"-enable-reactions",
		"-enable-user-links",
		"https://github.com/owner/repo/issues/123",
		"output.md",
	}

	cfg, exitCode := LoadFromFlags(args, stdout, stderr)

	if exitCode != -1 {
		t.Errorf("expected exitCode -1, got %d", exitCode)
	}

	if cfg.URL != "https://github.com/owner/repo/issues/123" {
		t.Errorf("expected URL to be 'https://github.com/owner/repo/issues/123', got '%s'", cfg.URL)
	}

	if cfg.OutputFile != "output.md" {
		t.Errorf("expected OutputFile to be 'output.md', got '%s'", cfg.OutputFile)
	}

	if !cfg.EnableReactions {
		t.Error("expected EnableReactions to be true")
	}

	if !cfg.EnableUserLinks {
		t.Error("expected EnableUserLinks to be true")
	}

	if cfg.Token != "ghp_test_token" {
		t.Errorf("expected Token to be 'ghp_test_token', got '%s'", cfg.Token)
	}
}
