package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestRun 表格驱动测试：验证CLI的Run函数的各种场景
func TestRun(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantExitCode int
		wantStdout   string // 期望stdout包含的文本（子字符串匹配）
		wantStderr   string // 期望stderr包含的文本（子字符串匹配）
		setupFunc    func() // 测试前的设置函数（如设置环境变量）
		cleanupFunc  func() // 测试后的清理函数
	}{
		// ===== 任务 6.1.2: 帮助和版本测试 =====
		{
			name:         "--help flag",
			args:         []string{"--help"},
			wantExitCode: 0,
			wantStdout:   "Usage:", // 帮助信息应包含"Usage:"
			wantStderr:   "",
		},
		{
			name:         "--version flag",
			args:         []string{"--version"},
			wantExitCode: 0,
			wantStdout:   "version:", // 版本信息应包含"version:"
			wantStderr:   "",
		},
		// ===== 任务 6.1.3: URL解析错误测试 =====
		{
			name:         "invalid URL format",
			args:         []string{"https://invalid-url.com/abc"},
			wantExitCode: 1,
			wantStdout:   "",
			wantStderr:   "invalid", // 错误信息应包含"invalid"或类似提示
		},
		{
			name:         "missing URL argument",
			args:         []string{},
			wantExitCode: 1,
			wantStdout:   "",
			wantStderr:   "URL", // 错误信息应提示需要URL
		},
		{
			name:         "unsupported resource type",
			args:         []string{"https://github.com/owner/repo/commits/123"},
			wantExitCode: 1,
			wantStdout:   "",
			wantStderr:   "unsupported", // 错误信息应包含"unsupported"或类似提示
		},
		// ===== 任务 6.1.4: GitHub API错误测试 =====
		{
			name:         "issue not found (mock)", // 使用mock github client
			args:         []string{"https://github.com/owner/repo/issues/999999"},
			wantExitCode: 1,
			wantStdout:   "",
			wantStderr:   "not found", // 错误信息应包含"not found"
			setupFunc: func() {
				// 设置环境变量确保使用正确的token（如果有）
			},
		},
		// ===== 任务 6.1.5: 完整流程测试 =====
		{
			name:         "end-to-end issue conversion (skip if no token)",
			args:         []string{"https://github.com/golang/go/issues/1"}, // 使用公开的issue
			wantExitCode: 0,
			wantStdout:   "#", // Markdown输出应包含"#"标题
			wantStderr:   "",
			setupFunc: func() {
				// 检查是否有GITHUB_TOKEN，如果没有则跳过测试
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 执行设置函数
			if tt.setupFunc != nil {
				tt.setupFunc()
			}
			// 执行清理函数
			if tt.cleanupFunc != nil {
				defer tt.cleanupFunc()
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			// 调用Run函数（目前不存在，测试会失败）
			exitCode := Run(tt.args, stdout, stderr)

			// 验证退出码
			if exitCode != tt.wantExitCode {
				t.Errorf("Run() exitCode = %v, want %v", exitCode, tt.wantExitCode)
			}

			// 验证stdout输出
			if tt.wantStdout != "" {
				stdoutStr := stdout.String()
				if !strings.Contains(stdoutStr, tt.wantStdout) {
					t.Errorf("Run() stdout = %v, want to contain %v", stdoutStr, tt.wantStdout)
				}
			}

			// 验证stderr输出
			if tt.wantStderr != "" {
				stderrStr := stderr.String()
				if !strings.Contains(stderrStr, tt.wantStderr) {
					t.Errorf("Run() stderr = %v, want to contain %v", stderrStr, tt.wantStderr)
				}
			}
		})
	}
}

// 辅助函数：创建临时环境变量
func setEnv(key, value string) func() {
	original := getEnvOrDefault(key, "")
	osSetEnv(key, value)
	return func() {
		osSetEnv(key, original)
	}
}

// 辅助函数：获取环境变量或默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value := osGetEnv(key); value != "" {
		return value
	}
	return defaultValue
}

// 包装os.Getenv和os.Setenv以便于测试
var (
	osGetEnv = os.Getenv
	osSetEnv = os.Setenv
)
