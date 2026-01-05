package config

// Config 应用配置
type Config struct {
	// 输入
	URL string

	// 输出
	OutputFile string // 空字符串表示stdout

	// 功能开关
	EnableReactions bool
	EnableUserLinks bool

	// 认证
	Token string // 从环境变量GITHUB_TOKEN读取
}
