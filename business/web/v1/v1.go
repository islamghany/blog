package v1

import (
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/web/v1/mid"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"os"

	"github.com/jmoiron/sqlx"
)

// hande options

type Options struct {
	corsOptions []string
}

func WithCors(whitelist []string) func(*Options) {
	return func(o *Options) {
		o.corsOptions = whitelist
	}
}

// the config that we want to pass to the handlers
type WebMuxConfig struct {
	Log       *logger.Logger
	shutdown  chan os.Signal
	DB        *sqlx.DB
	Whitelist []string
	Auth      *auth.Auth
}

// an interface that we use to inject routes from the handlers
type RouterAdder interface {
	Add(app *web.App, cfg *WebMuxConfig)
}

// creating new WebMux which return the web.app which has the server mux inside it embeding
func WebMux(cfg *WebMuxConfig, routerAdder RouterAdder) *web.App {
	app := web.NewApp(cfg.shutdown, mid.Logger(cfg.Log), mid.Errors(cfg.Log), mid.Metrics(), mid.Panics())

	routerAdder.Add(app, cfg)
	return app
}
