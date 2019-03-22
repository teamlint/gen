package template

var BootstrapMiddlewareTmpl = `package middleware

import (
	"{{.Config.Bootstrap.Import}}/server"

	"github.com/teamlint/iris"
)

func Sessions(s *server.Server) iris.Handler {
	return func(ctx iris.Context) {
		ctx.Application().Logger().Info("sessions shift expiration")
		s.Sessions.ShiftExpiration(ctx)
		ctx.Next()
	}
}`
