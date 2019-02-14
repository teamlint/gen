package query

import (
	"github.com/jinzhu/gorm"
)

type StockInoutItem struct {
}

func (q *StockInoutItem) QueryScopes() func(db *gorm.DB) *gorm.DB {
	if q == nil {
		return func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		// todo
		return db
	}
}
