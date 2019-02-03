package template

var ModelTmpl = `package {{.PackageName}}

import (
    "time"
)

var (
    _ = time.Second
)


type {{.StructName}} struct {
    {{range .Fields}}{{.}}
    {{end}}
}

// TableName sets the insert table name for this struct type
func ({{.ShortStructName}} *{{.StructName}}) TableName() string {
	return "{{.TableName}}"
}
`
