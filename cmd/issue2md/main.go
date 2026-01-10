package main

import (
	"fmt"
	"os"

	"github.com/wuwenrufeng/issue2md/internal/cli"
)

// 版本信息（构建时通过 -ldflags 注入）
var (
	// Version 版本号
	Version = "dev"
	// BuildDate 构建日期
	BuildDate = "unknown"
)

func main() {
	// 设置版本信息到 cli 包
	cli.SetVersion(Version, BuildDate)

	// 预处理：检查 --version 和 --help（避免 config 包依赖 cli 包）
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version":
			fmt.Printf("issue2md version: %s\n", Version)
			fmt.Printf("build date: %s\n", BuildDate)
			os.Exit(0)
		case "-version":
			fmt.Printf("issue2md version: %s\n", Version)
			fmt.Printf("build date: %s\n", BuildDate)
			os.Exit(0)
		}
	}

	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
