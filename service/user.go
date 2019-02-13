package service

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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
type userService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{db}
}

func (s *userService) Create(item *model.User) (err error) {
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

func (s *userService) Get(id interface{}, unscoped ...bool) (*model.User, error) {
	var item model.User

	var permanently bool
	if len(unscoped) > 0 && unscoped[0] {
		permanently = true
	}
	if err := s.DB.Scopes(query.Unscoped(permanently)).Where("id=?", id).Take(&item).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *userService) Update(item *model.User) (err error) {
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

func (s *userService) UpdateSel(item *model.User, sel []string) (err error) {
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
	item.UpdatedAt = time.Now()
	sel = append(sel, "updated_at")
	if e := tx.Unscoped().Select(sel).Save(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *userService) Delete(id interface{}, unscoped ...bool) (err error) {
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
	if e := tx.Scopes(query.Unscoped(permanently)).Where("id=?", id).Delete(&model.User{}).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *userService) Undelete(id interface{}) (err error) {
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
	if e := tx.Model(&model.User{}).Unscoped().Where("id=?", id).Update("deleted_at", nil).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *userService) GetList(base *query.Base, q *query.User) ([]*model.User, int, error) {
	var items []*model.User
	var total int

	db := s.DB.Model(&model.User{}).
		Scopes(base.OrderScopes()).
		Scopes(base.OrderByScopes()).
		Scopes(q.QueryScopes())
	err := db.Count(&total).Scopes(base.PagedScopes()).Scan(&items).Error

	return items, total, err
}
