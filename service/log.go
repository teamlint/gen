package service

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/teamlint/gen/model"
	"github.com/teamlint/gen/model/query"
)

type logService struct {
	DB *gorm.DB
}

func NewLogService(db *gorm.DB) LogService {
	return &logService{db}
}

func (s *logService) Create(item *model.Log) (err error) {
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

func (s *logService) Get(id interface{}, unscoped ...bool) (*model.Log, error) {
	var item model.Log

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

func (s *logService) Update(item *model.Log) (err error) {
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

func (s *logService) UpdateSel(item *model.Log, sel []string) (err error) {
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
	if e := tx.Unscoped().Select(sel).Save(item).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *logService) Delete(id interface{}, unscoped ...bool) (err error) {
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
	if e := tx.Scopes(query.Unscoped(permanently)).Where("id=?", id).Delete(&model.Log{}).Error; e != nil {
		tx.Rollback()
		return e
	}

	return tx.Commit().Error
}

func (s *logService) GetList(base *query.Base, q *query.Log) ([]*model.Log, int, error) {
	var items []*model.Log
	var total int

	db := s.DB.Model(&model.Log{}).
		Scopes(base.OrderScopes()).
		Scopes(base.OrderByScopes()).
		Scopes(q.QueryScopes())
	err := db.Count(&total).Scopes(base.PagedScopes()).Scan(&items).Error

	return items, total, err
}
