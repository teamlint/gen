package template

var BootstrapControllerTmpl = `package controller

import (
	"github.com/teamlint/iris"
	"github.com/teamlint/iris/mvc"
	"github.com/teamlint/iris/sessions"
)

type HomeController struct {
	Context iris.Context
	Session *sessions.Session
}

func (c *HomeController) Get() mvc.Result {
	c.Session.Set("foo", "foo value")

	data := iris.Map{
		"Title": "Home",
		"Message": "This is home page.",
	}
	
	return mvc.View{
		Name:   "home/index.html",
		Layout: "shared/layout.html",
		Data:   data,
	}
}

func (c *HomeController) GetAbout() mvc.Result {
	sv := c.Session.GetString("foo")

	data := iris.Map{
		"Title": "About",
		"Message": "This is about page.",
		"SessFoo": sv,
	}
	
	return mvc.View{
		Name:   "home/about.html",
		Data:   data,
	}
}`
