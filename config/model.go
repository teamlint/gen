package config

// Model 模型生成配置
type Model struct {
	Package  string   `toml:"package" json:"package"`   // 包名
	Import   string   `toml:"import" json:"import"`     // 导入包路径
	Enabled  bool     `toml:"enabled" json:"enabled"`   // 是否允许生成
	Base     bool     `toml:"base" json:"base"`         // 是否生成基础模型
	Guregu   bool     `toml:"guregu" json:"guregu"`     // 是否使用guregu/null类型包
	Tags     []string `toml:"tags" json:"tags"`         // 模型标记
	Template string   `toml:"template" json:"template"` // 模板名称
}
