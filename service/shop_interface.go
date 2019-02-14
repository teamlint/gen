package service

import (
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

// ShopService shop service interface
type ShopService interface {
	Create(item *model.Shop) error
	Get(id interface{}, unscoped ...bool) (*model.Shop, error)
	Update(item *model.Shop) error
	UpdateSel(item *model.Shop, sel []string) error
	Delete(id interface{}, unscoped ...bool) error
	Undelete(id interface{}) error
	GetList(base *query.Base, q *query.Shop) ([]*model.Shop, int, error)
}
