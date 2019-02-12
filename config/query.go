package config

// Query 查询生成配置
type Query struct {
	Package  string `toml:"package" json:"package"`   // 包名
	Import   string `toml:"import" json:"import"`     // 导入包路径
	Enabled  bool   `toml:"enabled" json:"enabled"`   // 是否允许生成
	Base     bool   `toml:"base" json:"base"`         // 是否生成基础查询结构
	Template string `toml:"template" json:"template"` // 模板名称
}
