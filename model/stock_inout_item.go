package model

import (
	"time"
)

var (
	_ = time.Second
)

type StockInoutItem struct {
	ID          string     `gorm:"column:id;primary_key" json:"id"`
	EntID       string     `gorm:"column:ent_id" json:"ent_id"`
	StockType   int        `gorm:"column:stock_type" json:"stock_type"`
	StockID     string     `gorm:"column:stock_id" json:"stock_id"`
	StockNumber string     `gorm:"column:stock_number" json:"stock_number"`
	ProductID   string     `gorm:"column:product_id" json:"product_id"`
	Price       int64      `gorm:"column:price" json:"price"`
	Quantity    int        `gorm:"column:quantity" json:"quantity"`
	Fee         int64      `gorm:"column:fee" json:"fee"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (s *StockInoutItem) TableName() string {
	return "stock_inout_items"
}
