package authapi

import (
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

type Config struct {
	Auth *auth.Auth
}

func Routes(app *web.App, cfg Config) {
	version := "v1"
	authHandler := NewAuthHandler(cfg.Auth)
	app.Handle(http.MethodPost, version, "/auth/login", authHandler.Authorize)
}
