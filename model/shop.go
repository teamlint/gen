package model

import (
	"time"
)

var (
	_ = time.Second
)

type Shop struct {
	ID           string     `gorm:"column:id;primary_key" json:"id"`
	Name         string     `gorm:"column:name" json:"name"`
	Description  string     `gorm:"column:description" json:"description"`
	Logo         string     `gorm:"column:logo" json:"logo"`
	License      string     `gorm:"column:license" json:"license"`
	AreaCode     *string    `gorm:"column:area_code" json:"area_code"`
	Province     string     `gorm:"column:province" json:"province"`
	City         string     `gorm:"column:city" json:"city"`
	District     string     `gorm:"column:district" json:"district"`
	Address      string     `gorm:"column:address" json:"address"`
	Flag         int        `gorm:"column:flag" json:"flag"`
	Status       int        `gorm:"column:status" json:"status"`
	Pinyin       *string    `gorm:"column:pinyin" json:"pinyin"`
	PinyinFl     *string    `gorm:"column:pinyin_fl" json:"pinyin_fl"`
	ContactName  string     `gorm:"column:contact_name" json:"contact_name"`
	ContactPhone string     `gorm:"column:contact_phone" json:"contact_phone"`
	AgentID      string     `gorm:"column:agent_id" json:"agent_id"`
	ManagerID    string     `gorm:"column:manager_id" json:"manager_id"`
	UpdatedAt    time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	LicensePhoto string     `gorm:"column:license_photo" json:"license_photo"`
}

func (s *Shop) TableName() string {
	return "shops"
}
