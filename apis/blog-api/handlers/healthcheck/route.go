package healthcheck

import (
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	DB    *sqlx.DB
	Log   *logger.Logger
	Build string
}

func Routes(app *web.App, cfg Config) {
	const version = "v1"

	healthCheck := New(cfg.DB, cfg.Log, cfg.Build)
	app.Handle(http.MethodGet, version, "/healthcheck", healthCheck.HealthCheckHandler)
	app.Handle(http.MethodGet, version, "/liveness", healthCheck.LivenessHandler)
	app.Handle(http.MethodGet, version, "/readiness", healthCheck.ReadinessHandler)
}
