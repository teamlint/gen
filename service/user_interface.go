package service

import (
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

// UserService user service interface
type UserService interface {
	Create(item *model.User) error
	Get(id interface{}, unscoped ...bool) (*model.User, error)
	Update(item *model.User) error
	UpdateSel(item *model.User, sel []string) error
	Delete(id interface{}, unscoped ...bool) error
	Undelete(id interface{}) error
	GetList(base *query.Base, q *query.User) ([]*model.User, int, error)
}
