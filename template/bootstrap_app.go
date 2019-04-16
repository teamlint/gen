package template

var BootstrapAppTmpl = `package app

import (
	"time"

	"github.com/jinzhu/gorm"
	"{{.Config.Bootstrap.Import}}/configurator"
	"{{.Config.Bootstrap.Import}}/routes"
	"{{.Config.Bootstrap.Import}}/server"
	"github.com/teamlint/iris"
)
func New(cfg *iris.Configuration, db *gorm.DB) *server.Server {
	s := server.New("Bootstrap", "Bootstarp")
	s.SetupDebug(cfg.GetOther()["Debug"].(bool))
	s.SetupViews("./views", "shared/layout.html")
	s.SetupAssets("./static")
	s.SetupConfiguration(cfg)
	s.SetupErrors()
	// sessions
	s.SetupSessions(20*time.Minute, cfg.GetOther()["HashKey"].(string), cfg.GetOther()["BlockKey"].(string))
	// services
	s.Configure(configurator.Services(db))
	// middleware
	// s.MVC.UseFunc(...)
	// routes
	s.Configure(routes.Configure)
	return s
}`
