package model

import (
	"time"
)

var (
	_ = time.Second
)

type StockInoutItem struct {
	ID          string     `gorm:&#34;column:id;primary_key&#34; json:&#34;id&#34;`
	EntID       string     `gorm:&#34;column:ent_id&#34; json:&#34;ent_id&#34;`
	StockType   int        `gorm:&#34;column:stock_type&#34; json:&#34;stock_type&#34;`
	StockID     string     `gorm:&#34;column:stock_id&#34; json:&#34;stock_id&#34;`
	StockNumber string     `gorm:&#34;column:stock_number&#34; json:&#34;stock_number&#34;`
	ProductID   string     `gorm:&#34;column:product_id&#34; json:&#34;product_id&#34;`
	Price       int64      `gorm:&#34;column:price&#34; json:&#34;price&#34;`
	Quantity    int        `gorm:&#34;column:quantity&#34; json:&#34;quantity&#34;`
	Fee         int64      `gorm:&#34;column:fee&#34; json:&#34;fee&#34;`
	CreatedAt   time.Time  `gorm:&#34;column:created_at&#34; json:&#34;created_at&#34;`
	UpdatedAt   time.Time  `gorm:&#34;column:updated_at&#34; json:&#34;updated_at&#34;`
	DeletedAt   *time.Time `gorm:&#34;column:deleted_at&#34; json:&#34;deleted_at&#34;`
}

func (s *StockInoutItem) TableName() string {
	return "stock_inout_items"
}
