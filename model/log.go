package model

import (
	"time"
)

var (
	_ = time.Second
)

type Log struct {
	ID        int        `gorm:&#34;column:id;primary_key&#34; json:&#34;id&#34;`
	Type      string     `gorm:&#34;column:type&#34; json:&#34;type&#34;`
	Username  *string    `gorm:&#34;column:username&#34; json:&#34;username&#34;`
	Content   *string    `gorm:&#34;column:content&#34; json:&#34;content&#34;`
	CreatedAt *time.Time `gorm:&#34;column:created_at&#34; json:&#34;created_at&#34;`
}

func (l *Log) TableName() string {
	return "logs"
}
