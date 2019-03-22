package template

var BootstrapMainTmpl = `package main

import (
	"{{.Config.Bootstrap.Import}}/app"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/teamlint/golog"
	"github.com/teamlint/iris"
)

func main() {
	// config
	cfg := iris.YAML("./config.yml")
	// database
	constr, ok := cfg.GetOther()["DBConnectionString"]
	if !ok {
		golog.Error("database connection string not set")
		return
	}
	db, err := gorm.Open("mysql", constr)
	if err != nil {
		golog.Errorf("database connection error: %v\r\n", err)
	}
	defer db.Close()
	db.LogMode(cfg.GetOther()["Debug"].(bool))
	app.New(&cfg,db).Run("")
}`
