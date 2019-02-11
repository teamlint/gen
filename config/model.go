package config

// Model 模型生成配置
type Model struct {
	Package  string   // 包名
	Enabled  bool     // 是否允许生成
	Guregu   bool     // 是否使用guregu/null类型包
	Tags     []string // 模型标记
	Template string   // 模板名称
}
