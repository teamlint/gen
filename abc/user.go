package abc

import (
	"time"
)

var (
	_ = time.Second
)

type User struct {
	ID          int        `gorm:"column:id;primary_key" json:"id"`
	PID         int        `gorm:"column:pid" json:"pid"`
	Username    string     `gorm:"column:username" json:"username"`
	Nickname    *string    `gorm:"column:nickname" json:"nickname"`
	Password    string     `gorm:"column:password" json:"password"`
	Gender      *int       `gorm:"column:gender" json:"gender"`
	Age         int        `gorm:"column:age" json:"age"`
	Money       float32    `gorm:"column:money" json:"money"`
	Double      float64    `gorm:"column:double" json:"double"`
	Like        string     `gorm:"column:like" json:"like"`
	Description string     `gorm:"column:description" json:"description"`
	IsApproved  int        `gorm:"column:is_approved" json:"is_approved"`
	JoinTime    time.Time  `gorm:"column:join_time" json:"join_time"`
	LeaveTime   *time.Time `gorm:"column:leave_time" json:"leave_time"`
	CreatedAt   time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
	Fee         *float64   `gorm:"column:fee" json:"fee"`
	Fee2        float64    `gorm:"column:fee2" json:"fee2"`
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "users"
}
