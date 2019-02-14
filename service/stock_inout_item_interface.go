package service

import (
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

// StockInoutItemService stockinoutitem service interface
type StockInoutItemService interface {
	Create(item *model.StockInoutItem) error
	Get(id interface{}, unscoped ...bool) (*model.StockInoutItem, error)
	Update(item *model.StockInoutItem) error
	UpdateSel(item *model.StockInoutItem, sel []string) error
	Delete(id interface{}, unscoped ...bool) error
	Undelete(id interface{}) error
	GetList(base *query.Base, q *query.StockInoutItem) ([]*model.StockInoutItem, int, error)
}
