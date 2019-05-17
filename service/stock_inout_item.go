package service

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

type stockInoutItemService struct {
	DB *gorm.DB
}

func NewStockInoutItemService(db *gorm.DB) StockInoutItemService {
	return &stockInoutItemService{db}
}

func (s *stockInoutItemService) Create(item *model.StockInoutItem) (err error) {
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

func (s *stockInoutItemService) Get(id interface{}, unscoped ...bool) (*model.StockInoutItem, error) {
	var item model.StockInoutItem

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

func (s *stockInoutItemService) Update(item *model.StockInoutItem) (err error) {
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
	if e := tx.Unscoped().Omit("created_at").Save(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *stockInoutItemService) UpdateSel(item *model.StockInoutItem, sel []string) (err error) {
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

func (s *stockInoutItemService) Delete(id interface{}, unscoped ...bool) (err error) {
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
	if e := tx.Scopes(query.Unscoped(permanently)).Where("id=?", id).Delete(&model.StockInoutItem{}).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *stockInoutItemService) Undelete(id interface{}) (err error) {
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
	if e := tx.Model(&model.StockInoutItem{}).Unscoped().Where("id=?", id).Update("deleted_at", nil).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *stockInoutItemService) GetList(base *query.Base, q *query.StockInoutItem) ([]*model.StockInoutItem, int, error) {
	var items []*model.StockInoutItem
	var total int

	db := s.DB.Model(&model.StockInoutItem{}).
		Scopes(base.OrderScopes()).
		Scopes(base.OrderByScopes()).
		Scopes(q.QueryScopes())
	err := db.Count(&total).Scopes(base.PagedScopes()).Scan(&items).Error

	return items, total, err
}
