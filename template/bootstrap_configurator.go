package template

var BootstrapConfiguratorTmpl = `package configurator

import (
	"{{.Config.Bootstrap.Import}}/server"

	"github.com/jinzhu/gorm"
)

func Services(db *gorm.DB) server.Configurator {
	return func(s *server.Server) {
		// todo
		// s.Register(service.NewUserService(db))
	}
}`
