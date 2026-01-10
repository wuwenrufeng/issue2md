package cli

// Version 版本号（通过-ldflags注入或SetVersion设置）
var Version = "dev"

// BuildDate 构建日期（通过-ldflags注入或SetVersion设置）
var BuildDate = "unknown"

// SetVersion 设置版本信息（由 main 包调用）
func SetVersion(version, buildDate string) {
	if version != "" {
		Version = version
	}
	if buildDate != "" {
		BuildDate = buildDate
	}
}
