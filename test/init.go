package test

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/teamlint/gen/service"
)

var (
	db                *gorm.DB
	err               error
	userService       service.UserService
	customUserService service.CustomUserService
)

func init() {
	var constr = "root:123456@tcp(127.0.0.1:3306)/test?loc=Local&parseTime=true"
	db, err = gorm.Open("mysql", constr)
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	// defer db.Close()
	userService = service.NewUserService(db)
	customUserService = service.NewCustomUserService(db)
}
