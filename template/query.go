package template

var QueryTmpl = `package {{.Config.Query.Package}}

import (
	"github.com/jinzhu/gorm"
)

type {{.StructName}} struct {
}

func (q *{{.StructName}}) QueryScopes() func(db *gorm.DB) *gorm.DB {
	if q == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		// todo
		return db
	}
}`
