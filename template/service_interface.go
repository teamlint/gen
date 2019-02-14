package template

var ServiceInterfaceTmpl = `package {{.Config.Service.Package}}

import (
	"{{.Config.Model.Import}}"
	"{{.Config.Query.Import}}"
)

// {{.StructName}}Service {{.StructName|toLower}} service interface
type {{.StructName}}Service interface {
	Create(item *{{.PackageName}}.{{.StructName}}) error
	Get(id interface{}, unscoped ...bool) (*{{.PackageName}}.{{.StructName}}, error)
	Update(item *{{.PackageName}}.{{.StructName}}) error
	UpdateSel(item *{{.PackageName}}.{{.StructName}}, sel []string) error
	Delete(id interface{}, unscoped ...bool) error
	{{if .Helper.HasField "deleted_at" -}}
	Undelete(id interface{}) error
	{{end -}}
	GetList(base *query.Base, q *query.{{.StructName}}) ([]*{{.PackageName}}.{{.StructName}}, int, error)
}`
