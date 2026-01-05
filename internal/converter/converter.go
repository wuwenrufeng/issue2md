package converter

// Converter Markdown转换器
type Converter struct {
	enableReactions bool
	enableUserLinks bool
}

// Option 配置选项类型（函数式选项模式）
type Option func(*Converter)

// WithReactions 启用reactions
func WithReactions(enable bool) Option {
	return func(c *Converter) {
		c.enableReactions = enable
	}
}

// WithUserLinks 启用用户链接
func WithUserLinks(enable bool) Option {
	return func(c *Converter) {
		c.enableUserLinks = enable
	}
}

// NewConverter 创建新的Converter
func NewConverter(options ...Option) *Converter {
	c := &Converter{
		enableReactions: false,
		enableUserLinks: false,
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}
