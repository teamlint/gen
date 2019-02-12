package template

var ModelTmpl = `package {{.PackageName}}

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

// TableName sets the table name for this struct type
func ({{.ShortStructName}} *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
