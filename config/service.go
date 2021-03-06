package config

// Service 服务生成配置
type Service struct {
	Package   string `toml:"package" json:"package"`     // 包名
	Enabled   bool   `toml:"enabled" json:"enabled"`     // 是否允许生成
	Interface bool   `toml:"interface" json:"interface"` // 是否生成接口
	Template  string `toml:"template" json:"template"`   // 模板名称
}
