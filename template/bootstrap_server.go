package template

var BootstrapServerTmpl = `package server

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/teamlint/iris/mvc"

	"github.com/teamlint/iris"
	"github.com/teamlint/iris/sessions"
	"github.com/teamlint/iris/view"
	"github.com/teamlint/iris/websocket"
)

const (
	// StaticAssets is the root directory for public assets like images, css, js.
	defaultStaticAssets = "./static"
	// Favicon
	favicon = "favicon.ico"
)

// Configurator server configurator function
type Configurator func(*Server)

// ViewConfigurator view engine configurator function
type ViewConfigurator func(engine *view.HTMLEngine)

// Server app server struct
type Server struct {
	*iris.Application
	*iris.Configuration
	MVC          *mvc.Application
	AppName      string
	AppTitle     string
	AppSpawnDate time.Time
	Debug        bool

	Sessions *sessions.Sessions
}

// New returns a new server.
func New(appName, appTitle string, cfs ...Configurator) *Server {
	app := iris.Default()
	s := &Server{
		AppName:      appName,
		AppTitle:     appTitle,
		AppSpawnDate: time.Now(),
		Application:  app,
		MVC:          mvc.New(app),
	}

	for _, cf := range cfs {
		cf(s)
	}

	return s
}

// Default default server instance
// debug,views don't configure
func Default(appName, appTitle string, cfs ...Configurator) *Server {
	s := New(appName, appTitle, cfs...)
	s.SetupDebug(false)
	s.SetupViews("./views", "shared/layout.html")
	s.SetupAssets(defaultStaticAssets)
	s.SetupErrors()
	return s
}

// SetupConfiguration setup server configuration
func (s *Server) SetupConfiguration(cfg *iris.Configuration) {
	s.Configuration = cfg
	s.Application.Configure(iris.WithConfiguration(*cfg))
}

// SetupAssets setup static assets resources
func (s *Server) SetupAssets(assetsDir string) *Server {
	faviconPath := filepath.Join(assetsDir, "img", favicon)
	if _, err := os.Stat(faviconPath); !os.IsNotExist(err) {
		s.Favicon(faviconPath)
	}
	s.StaticWeb(assetsDir[1:], assetsDir)
	return s
}

// SetupViews loads the templates.
func (s *Server) SetupViews(viewsDir string, defaultLayout string, vcfs ...ViewConfigurator) *Server {
	tmpl := iris.HTML(viewsDir, ".html")
	if defaultLayout != "" {
		tmpl.Layout(defaultLayout)
	}
	tmpl.Reload(s.Debug)
	for _, vcf := range vcfs {
		vcf(tmpl)
	}
	s.Application.RegisterView(tmpl)
	return s
}

// SetupDebug
func (s *Server) SetupDebug(dbg bool) *Server {
	s.Debug = dbg
	if dbg {
		s.Application.Logger().SetLevel("debug")
	} else {
		s.Application.Logger().SetLevel("error")
	}
	return s
}

// SetupSessions initializes the sessions, optionally.
func (s *Server) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey string) *Server {
	cookieName := strings.ToUpper(strings.Replace(s.AppName, " ", "_", -1)) + "_SESS"
	config := sessions.Config{
		Cookie:       cookieName,
		Expires:      expires,
		AllowReclaim: true,
	}
	if cookieHashKey != "" && cookieBlockKey != "" {
		config.Encoding = securecookie.New([]byte(cookieHashKey), []byte(cookieBlockKey))
	}
	s.Sessions = sessions.New(config)
	s.Register(s.Sessions.Start)
	return s
}

// SetupWebsockets prepares the websocket server.
func (s *Server) SetupWebsockets(endpoint string, onConnection websocket.ConnectionFunc) *Server {
	ws := websocket.New(websocket.Config{})
	ws.OnConnection(onConnection)

	s.Get(endpoint, ws.Handler())
	s.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})
	return s
}

// SetupErrors prepares the http error handlers
func (s *Server) SetupErrors() *Server {
	s.Application.OnAnyErrorCode(func(ctx iris.Context) {
		statusCode := ctx.GetStatusCode()
		errCode := ctx.Values().GetIntDefault("err.code", statusCode)
		errMsg := ctx.Values().GetStringDefault("err.message", http.StatusText(statusCode))
		err := iris.Map{
			"Title":      s.AppTitle,
			"StatusCode": statusCode,
			"Code":       errCode,
			"Message":    errMsg,
		}
		// json
		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}
		// views
		ctx.ViewLayout(iris.NoLayout)
		ctx.ViewData("Error", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
		ctx.StopExecution()
	})
	return s
}

// Register appends one or more values as dependencies.
func (s *Server) Register(deps ...interface{}) *Server {
	s.MVC.Register(deps...)
	return s
}

// Configure accepts configurations and runs them inside the server's context.
func (s *Server) Configure(cs ...Configurator) *Server {
	for _, c := range cs {
		c(s)
	}
	return s
}

// Run starts the http server with the specified "addr".
func (s *Server) Run(addr string, cfgs ...iris.Configurator) {
	if addr != "" {
		s.Application.Run(iris.Addr(addr), cfgs...)
		return
	}
	if s.Configuration != nil {
		if port, ok := s.Configuration.GetOther()["ServerPort"].(int); ok {
			s.Application.Run(iris.Addr(":"+strconv.Itoa(port)), cfgs...)
			return
		}
	}
	s.Application.Logger().Fatal("must be configure server address")
}`
