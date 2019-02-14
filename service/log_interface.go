package service

import (
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

// LogService log service interface
type LogService interface {
	Create(item *model.Log) error
	Get(id interface{}, unscoped ...bool) (*model.Log, error)
	Update(item *model.Log) error
	UpdateSel(item *model.Log, sel []string) error
	Delete(id interface{}, unscoped ...bool) error
	GetList(base *query.Base, q *query.Log) ([]*model.Log, int, error)
}
