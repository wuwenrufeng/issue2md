package config

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// LoadFromFlags 从命令行参数和环境变量加载配置
// argv: 命令行参数切片（不包括程序名）
// stdout: 标准输出 writer
// stderr: 标准错误 writer
// 返回: (*Config, exitCode)
//   - exitCode = -1: 继续执行
//   - exitCode = 0: 正常退出（如 --help, --version）
//   - exitCode = 1: 错误退出
func LoadFromFlags(argv []string, stdout, stderr io.Writer) (*Config, int) {
	// 创建独立的 FlagSet，避免污染全局 flag
	fs := flag.NewFlagSet("issue2md", flag.ContinueOnError)

	// 定义 flag 变量
	var enableReactions bool
	var enableUserLinks bool
	var showVersion bool
	var showHelp bool

	// 注册 flag
	fs.BoolVar(&enableReactions, "enable-reactions", false, "显示 reactions 统计")
	fs.BoolVar(&enableUserLinks, "enable-user-links", false, "用户名显示为可点击链接")
	fs.BoolVar(&showVersion, "version", false, "显示版本信息")
	fs.BoolVar(&showHelp, "help", false, "显示帮助信息")

	// 设置 flag 的输出方向（用于错误信息）
	fs.SetOutput(stderr)

	// 解析命令行参数
	if err := fs.Parse(argv); err != nil {
		// 参数解析失败
		fmt.Fprintf(stderr, "参数解析错误: %v\n", err)
		return nil, 1
	}

	// 处理 --help
	if showHelp {
		printHelp(stdout)
		return nil, 0
	}

	// 处理 --version
	if showVersion {
		printVersion(stdout)
		return nil, 0
	}

	// 获取位置参数
	args := fs.Args()

	// 检查是否提供了 URL 参数
	if len(args) == 0 {
		fmt.Fprintln(stderr, "错误: 缺少必需参数 URL")
		fmt.Fprintln(stderr, "使用 --help 查看使用说明")
		return nil, 1
	}

	// 第一个位置参数是 URL
	url := args[0]

	// 第二个位置参数（可选）是输出文件
	outputFile := ""
	if len(args) > 1 {
		outputFile = args[1]
	}

	// 从环境变量读取 GitHub Token
	token := os.Getenv("GITHUB_TOKEN")

	// 构建配置
	cfg := &Config{
		URL:             url,
		OutputFile:      outputFile,
		EnableReactions: enableReactions,
		EnableUserLinks: enableUserLinks,
		Token:           token,
	}

	return cfg, -1
}

// printHelp 输出帮助信息
func printHelp(w io.Writer) {
	fmt.Fprintln(w, "issue2md - 将 GitHub Issue/PR/Discussion 转换为 Markdown")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  issue2md [flags] <URL> [output_file]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Arguments:")
	fmt.Fprintln(w, "  URL          GitHub Issue/PR/Discussion 的完整 URL")
	fmt.Fprintln(w, "  output_file  输出文件路径（可选，不提供则输出到 stdout）")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  -enable-reactions   显示 reactions 统计（如 👍 3 ❤️ 1）")
	fmt.Fprintln(w, "  -enable-user-links  用户名显示为可点击链接")
	fmt.Fprintln(w, "  -version            显示版本信息")
	fmt.Fprintln(w, "  -help               显示此帮助信息")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Environment Variables:")
	fmt.Fprintln(w, "  GITHUB_TOKEN  GitHub Personal Access Token（可选）")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  issue2md https://github.com/owner/repo/issues/123")
	fmt.Fprintln(w, "  issue2md -enable-reactions https://github.com/owner/repo/issues/123 output.md")
	fmt.Fprintln(w, "  GITHUB_TOKEN=ghp_xxx issue2md https://github.com/owner/repo/issues/123")
}

// printVersion 输出版本信息
func printVersion(w io.Writer) {
	// 版本信息从 internal/cli 包读取
	// 这里我们需要动态导入，但由于循环依赖问题，
	// 我们通过一个简单的方式来获取版本
	fmt.Fprintf(w, "issue2md version: %s\n", getVersion())
	fmt.Fprintf(w, "build date: %s\n", getBuildDate())
}

// getVersion 获取版本号
func getVersion() string {
	// 从环境变量或编译时注入的值获取
	// 这里使用默认值 "dev"
	if v := os.Getenv("ISSUE2MD_VERSION"); v != "" {
		return v
	}
	return "dev"
}

// getBuildDate 获取构建日期
func getBuildDate() string {
	// 从环境变量或编译时注入的值获取
	// 这里使用默认值 "unknown"
	if d := os.Getenv("ISSUE2MD_BUILD_DATE"); d != "" {
		return d
	}
	return "unknown"
}
