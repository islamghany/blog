package healthcheck

import (
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	DB  *sqlx.DB
	Log *logger.Logger
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	healthCheck := New(cfg.DB, cfg.Log)
	app.Handle(http.MethodGet, version, "/healthcheck", healthCheck.HealthCheckHandler)
}
