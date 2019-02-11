package {{.PackageName}}

import (
    "time"
	{{if .Config.Model.Guregu}}
	"github.com/guregu/null"
	{{end}}
)

var (
    _ = time.Second
)


type {{.StructName}} struct {
    {{range .Fields}}{{.}}
    {{end}}
}

// TableName 表名称
func ({{.ShortStructName}} *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}

