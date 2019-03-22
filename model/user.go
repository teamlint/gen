package model

import (
	"time"
)

var (
	_ = time.Second
)

type User struct {
	ID          int        `gorm:&#34;column:id;primary_key&#34; json:&#34;id&#34;`
	PID         int        `gorm:&#34;column:pid&#34; json:&#34;pid&#34;`
	Username    string     `gorm:&#34;column:username&#34; json:&#34;username&#34;`
	Nickname    *string    `gorm:&#34;column:nickname&#34; json:&#34;nickname&#34;`
	Password    string     `gorm:&#34;column:password&#34; json:&#34;password&#34;`
	Sex         UserGender `gorm:&#34;column:gender&#34; json:&#34;gender&#34;`
	Age         int        `gorm:&#34;column:age&#34; json:&#34;age&#34;`
	Money       float32    `gorm:&#34;column:money&#34; json:&#34;money&#34;`
	Double      float64    `gorm:&#34;column:double&#34; json:&#34;double&#34;`
	Like        string     `gorm:&#34;column:like&#34; json:&#34;like&#34;`
	Description string     `gorm:&#34;column:description&#34; json:&#34;description&#34;`
	IsApproved  bool       `gorm:&#34;column:is_approved&#34; json:&#34;is_approved&#34;`
	JoinTime    time.Time  `gorm:&#34;column:join_time&#34; json:&#34;join_time&#34;`
	LeaveTime   *time.Time `gorm:&#34;column:leave_time&#34; json:&#34;leave_time&#34;`
	CreatedAt   time.Time  `gorm:&#34;column:created_at&#34; json:&#34;created_at&#34;`
	UpdatedAt   time.Time  `gorm:&#34;column:updated_at&#34; json:&#34;updated_at&#34;`
	DeletedAt   *time.Time `gorm:&#34;column:deleted_at&#34; json:&#34;deleted_at&#34;`
	Fee         *float64   `gorm:&#34;column:fee&#34; json:&#34;fee&#34;`
	Fee2        float64    `gorm:&#34;column:fee2&#34; json:&#34;fee2&#34;`
}

func (u *User) TableName() string {
	return "users"
}
