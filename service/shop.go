package service

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

type shopService struct {
	DB *gorm.DB
}

func NewShopService(db *gorm.DB) ShopService {
	return &shopService{db}
}

func (s *shopService) Create(item *model.Shop) (err error) {
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

func (s *shopService) Get(id interface{}, unscoped ...bool) (*model.Shop, error) {
	var item model.Shop

	var permanently bool
	if len(unscoped) > 0 && unscoped[0] {
		permanently = true
	}
	if err := s.DB.Scopes(query.Unscoped(permanently)).Where("id=?", id).Take(&item).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrRecordNotFound
		}
		return nil, err
	}

	return &item, nil
}

func (s *shopService) Update(item *model.Shop) (err error) {
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

func (s *shopService) UpdateSel(item *model.Shop, sel []string) (err error) {
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

func (s *shopService) Delete(id interface{}, unscoped ...bool) (err error) {
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
	if e := tx.Scopes(query.Unscoped(permanently)).Where("id=?", id).Delete(&model.Shop{}).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *shopService) Undelete(id interface{}) (err error) {
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
	if e := tx.Model(&model.Shop{}).Unscoped().Where("id=?", id).Update("deleted_at", nil).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *shopService) GetList(base *query.Base, q *query.Shop) ([]*model.Shop, int, error) {
	var items []*model.Shop
	var total int

	db := s.DB.Model(&model.Shop{}).
		Scopes(base.OrderScopes()).
		Scopes(base.OrderByScopes()).
		Scopes(q.QueryScopes())
	err := db.Count(&total).Scopes(base.PagedScopes()).Scan(&items).Error

	return items, total, err
}
