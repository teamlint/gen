package template

var ServiceTmpl = `package {{.Config.Service.Package}}

import (
	"fmt"
	{{if .Helper.HasField "updated_at" -}}
	"time"
	{{end}}
	"{{.Config.Model.Import}}"
	"{{.Config.Query.Import}}"
	"github.com/jinzhu/gorm"
)

type {{.StructName|toLowerCamelCase}}Service struct {
	DB *gorm.DB
}

func New{{.StructName}}Service(db *gorm.DB) {{.StructName}}Service {
	return &{{.StructName|toLowerCamelCase}}Service{db}
}

func (s *{{.StructName|toLowerCamelCase}}Service) Create(item *{{.PackageName}}.{{.StructName}}) (err error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if e := tx.Create(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *{{.StructName|toLowerCamelCase}}Service) Get(id interface{}, unscoped ...bool) (*{{.PackageName}}.{{.StructName}}, error) {
	var item {{.PackageName}}.{{.StructName}}

	var permanently bool
	if len(unscoped) > 0 && unscoped[0] {
		permanently = true
	}
	if err := s.DB.Scopes({{.Config.Query.Package}}.Unscoped(permanently)).Where("id=?", id).Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, {{.Config.Model.Package}}.ErrRecordNotFound
		}
		return nil, err
	}

	return &item, nil
}

func (s *{{.StructName|toLowerCamelCase}}Service) Update(item *{{.PackageName}}.{{.StructName}}) (err error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if e := tx.Unscoped().Save(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *{{.StructName|toLowerCamelCase}}Service) UpdateSel(item *{{.PackageName}}.{{.StructName}}, sel []string) (err error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	{{if .Helper.HasField "updated_at" -}}
	item.UpdatedAt = time.Now()
	sel = append(sel, "updated_at")
	{{end -}}
	if e := tx.Unscoped().Select(sel).Save(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *{{.StructName|toLowerCamelCase}}Service) Delete(id interface{}, unscoped ...bool) (err error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	var permanently bool
	if len(unscoped) > 0 && unscoped[0] {
		permanently = true
	}
	if e := tx.Scopes({{.Config.Query.Package}}.Unscoped(permanently )).Where("id=?", id).Delete(&{{.PackageName}}.{{.StructName}}{}).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

{{if .Helper.HasField "deleted_at" -}}
func (s *{{.StructName|toLowerCamelCase}}Service) Undelete(id interface{}) (err error) {
	tx := s.DB.Begin()
	if tx.Error != nil {
		err = tx.Error
		return
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	if e := tx.Model(&{{.PackageName}}.{{.StructName}}{}).Unscoped().Where("id=?", id).Update("deleted_at",nil).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}
{{end}}

func (s *{{.StructName|toLowerCamelCase}}Service) GetList(base *query.Base, q *query.{{.StructName}}) ([]*{{.PackageName}}.{{.StructName}}, int, error) {
	var items []*{{.PackageName}}.{{.StructName}}
	var total int

	db := s.DB.Model(&{{.PackageName}}.{{.StructName}}{}).
		Scopes(base.OrderScopes()).
		Scopes(base.OrderByScopes()).
		Scopes(q.QueryScopes())
	err := db.Count(&total).Scopes(base.PagedScopes()).Scan(&items).Error

	return items, total, err
}`
