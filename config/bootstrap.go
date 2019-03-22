package config

// Bootstrap 程序框架生成配置
type Bootstrap struct {
	Package    string `toml:"package" json:"package"`       // 包名
	Import     string `toml:"import" json:"import"`         // 导入包路径
	Enabled    bool   `toml:"enabled" json:"enabled"`       // 是否允许生成
	Server     bool   `toml:"server" json:"server"`         // 服务主程序
	Router     bool   `toml:"router" json:"router"`         // 路由
	Controller bool   `toml:"controller" json:"controller"` // 控制器
	Middleware bool   `toml:"middleware" json:"middleware"` // 中间件
	Global     bool   `toml:"global" json:"global"`         // 全局处理程序
	Template   string `toml:"template" json:"template"`     // 模板名称
}
