package query

import (
	"github.com/jinzhu/gorm"
	"github.com/teamlint/gen/model"
)

type User struct {
	Gender model.UserGender
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
