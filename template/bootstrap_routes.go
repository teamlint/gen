package template

var BootstrapRoutesTmpl = `package routes

import (
	"{{.Config.Bootstrap.Import}}/controller"
	"{{.Config.Bootstrap.Import}}/middleware"
	"{{.Config.Bootstrap.Import}}/server"
)

func Configure(s *server.Server) {
	s.Use(middleware.Sessions(s))
	home := s.MVC.Party("/")
	home.Handle(new(controller.HomeController))
}`
