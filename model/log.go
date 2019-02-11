package model

import (
	"time"
)

var (
	_ = time.Second
)

type Log struct {
	ID        int        `gorm:"column:id;primary_key" json:"id"`
	Type      string     `gorm:"column:type" json:"type"`
	Username  *string    `gorm:"column:username" json:"username"`
	Content   *string    `gorm:"column:content" json:"content"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 表名称
func (l *Log) TableName() string {
	return "logs"
}
