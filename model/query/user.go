package query

import (
	"github.com/jinzhu/gorm"
)

type User struct {
}

func (q *User) QueryScopes() func(db *gorm.DB) *gorm.DB {
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
