package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/wuwenrufeng/issue2md/internal/config"
	"github.com/wuwenrufeng/issue2md/internal/converter"
	"github.com/wuwenrufeng/issue2md/internal/github"
	"github.com/wuwenrufeng/issue2md/internal/parser"
)

// Run 是CLI的主入口函数
// argv: 命令行参数切片（不包括程序名）
// stdout: 标准输出 writer
// stderr: 标准错误 writer
// 返回: exitCode (0: 成功, 1: 错误)
func Run(argv []string, stdout, stderr io.Writer) int {
	// 1. 加载配置
	cfg, exitCode := config.LoadFromFlags(argv, stdout, stderr)
	if exitCode != -1 {
		return exitCode
	}
	// 此时 cfg != nil，程序继续执行

	// 2. 解析URL
	resource, err := parser.ParseURL(cfg.URL)
	if err != nil {
		fmt.Fprintf(stderr, "URL解析错误: %v\n", err)
		return 1
	}

	// 3. 创建GitHub客户端
	client := github.NewClient(cfg.Token)

	// 4. 根据资源类型获取数据
	var markdown string
	var fetchErr error

	switch resource.Type {
	case parser.Issue:
		issue, err := client.FetchIssue(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			fetchErr = err
			break
		}
		conv := converter.NewConverter(
			converter.WithReactions(cfg.EnableReactions),
			converter.WithUserLinks(cfg.EnableUserLinks),
		)
		markdown, err = conv.ConvertIssue(issue)
		if err != nil {
			fetchErr = err
		}
	case parser.PullRequest:
		pr, err := client.FetchPullRequest(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			fetchErr = err
			break
		}
		conv := converter.NewConverter(
			converter.WithReactions(cfg.EnableReactions),
			converter.WithUserLinks(cfg.EnableUserLinks),
		)
		markdown, err = conv.ConvertPullRequest(pr)
		if err != nil {
			fetchErr = err
		}
	case parser.Discussion:
		discussion, err := client.FetchDiscussion(resource.Owner, resource.Repo, resource.Number)
		if err != nil {
			fetchErr = err
			break
		}
		conv := converter.NewConverter(
			converter.WithReactions(cfg.EnableReactions),
			converter.WithUserLinks(cfg.EnableUserLinks),
		)
		markdown, err = conv.ConvertDiscussion(discussion)
		if err != nil {
			fetchErr = err
		}
	default:
		fmt.Fprintf(stderr, "不支持的资源类型: %v\n", resource.Type)
		return 1
	}

	// 5. 处理获取/转换错误
	if fetchErr != nil {
		fmt.Fprintf(stderr, "GitHub API错误: %v\n", fetchErr)
		return 1
	}

	// 6. 输出结果
	if cfg.OutputFile == "" {
		// 输出到stdout
		fmt.Fprint(stdout, markdown)
	} else {
		// 输出到文件
		err := os.WriteFile(cfg.OutputFile, []byte(markdown), 0644)
		if err != nil {
			fmt.Fprintf(stderr, "文件写入错误: %v\n", err)
			return 1
		}
	}

	return 0
}
