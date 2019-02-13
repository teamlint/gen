package query

import (
	"github.com/jinzhu/gorm"
)

type Base struct {
	PageNumber int    `form:"pageNumber" json:"page_number"`
	PageSize   int    `form:"pageSize" json:"page_size"`
	SortName   string `form:"sortName" json:"sort_name"`
	SortOrder  string `form:"sortOrder" json:"sort_order"`
	SearchText string `form:"searchText" json:"search_text"`
	OrderBy    string `form:"orderBy" json:"order_by"`
}

// PagedScopes 分页查询拼接
func (b *Base) PagedScopes() func(db *gorm.DB) *gorm.DB {
	if ok, fn := b.IsNil(); ok {
		return fn
	}
	return func(db *gorm.DB) *gorm.DB {
		if b.PageNumber < 1 {
			b.PageNumber = 1
		}
		if b.PageSize <= 0 {
			return db
		}
		offset := b.PageNumber*b.PageSize - b.PageSize
		return db.Offset(offset).Limit(b.PageSize)
	}
}

// OrderScopes 查询排序拼接
func (b *Base) OrderScopes(defaultOrder ...string) func(db *gorm.DB) *gorm.DB {
	if ok, fn := b.IsNil(); ok {
		return fn
	}
	return func(db *gorm.DB) *gorm.DB {
		if b.SortName == "" {
			if len(defaultOrder) > 0 {
				b.SortName = defaultOrder[0]
			} else {
				return db
			}
		}
		if b.SortOrder == "" {
			if len(defaultOrder) > 1 {
				b.SortOrder = defaultOrder[1]
			} else {
				b.SortOrder = "Desc"
			}
		}
		return db.Order(b.SortName + " " + b.SortOrder)
	}
}

// OrdeyByScopes 多排序
func (b *Base) OrderByScopes() func(db *gorm.DB) *gorm.DB {
	if ok, fn := b.IsNil(); ok {
		return fn
	}
	return func(db *gorm.DB) *gorm.DB {
		if b.OrderBy != "" {
			return db.Order(b.OrderBy)
		}
		return db
	}
}

func (b *Base) IsNil() (bool, func(db *gorm.DB) *gorm.DB) {
	if b == nil {
		return true, func(db *gorm.DB) *gorm.DB {
			return db
		}
	}
	return false, nil
}

func Unscoped(unscoped bool) func(db *gorm.DB) *gorm.DB {
	if unscoped {
		return func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}
	}
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}
