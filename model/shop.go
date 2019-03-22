package model

import (
	"time"
)

var (
	_ = time.Second
)

type Shop struct {
	ID           string     `gorm:&#34;column:id;primary_key&#34; json:&#34;id&#34;`
	Name         string     `gorm:&#34;column:name&#34; json:&#34;name&#34;`
	Description  string     `gorm:&#34;column:description&#34; json:&#34;description&#34;`
	Logo         string     `gorm:&#34;column:logo&#34; json:&#34;logo&#34;`
	License      string     `gorm:&#34;column:license&#34; json:&#34;license&#34;`
	AreaCode     *string    `gorm:&#34;column:area_code&#34; json:&#34;area_code&#34;`
	Province     string     `gorm:&#34;column:province&#34; json:&#34;province&#34;`
	City         string     `gorm:&#34;column:city&#34; json:&#34;city&#34;`
	District     string     `gorm:&#34;column:district&#34; json:&#34;district&#34;`
	Address      string     `gorm:&#34;column:address&#34; json:&#34;address&#34;`
	Flag         int        `gorm:&#34;column:flag&#34; json:&#34;flag&#34;`
	Status       int        `gorm:&#34;column:status&#34; json:&#34;status&#34;`
	Pinyin       *string    `gorm:&#34;column:pinyin&#34; json:&#34;pinyin&#34;`
	PinyinFl     *string    `gorm:&#34;column:pinyin_fl&#34; json:&#34;pinyin_fl&#34;`
	ContactName  string     `gorm:&#34;column:contact_name&#34; json:&#34;contact_name&#34;`
	ContactPhone string     `gorm:&#34;column:contact_phone&#34; json:&#34;contact_phone&#34;`
	AgentID      string     `gorm:&#34;column:agent_id&#34; json:&#34;agent_id&#34;`
	ManagerID    string     `gorm:&#34;column:manager_id&#34; json:&#34;manager_id&#34;`
	CreatedAt    time.Time  `gorm:&#34;column:created_at&#34; json:&#34;created_at&#34;`
	UpdatedAt    time.Time  `gorm:&#34;column:updated_at&#34; json:&#34;updated_at&#34;`
	DeletedAt    *time.Time `gorm:&#34;column:deleted_at&#34; json:&#34;deleted_at&#34;`
	LicensePhoto string     `gorm:&#34;column:license_photo&#34; json:&#34;license_photo&#34;`
}

func (s *Shop) TableName() string {
	return "shops"
}
