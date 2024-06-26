package handlers

import (
	"github/islamghany/blog/apis/blog-api/handlers/articleapi"
	"github/islamghany/blog/apis/blog-api/handlers/authapi"
	"github/islamghany/blog/apis/blog-api/handlers/healthcheck"
	"github/islamghany/blog/apis/blog-api/handlers/usergrp"
	v1 "github/islamghany/blog/business/web/v1"
	"github/islamghany/blog/foundation/web"
)

type Routes struct{}

func (Routes) Add(app *web.App, cfg *v1.WebMuxConfig) {
	healthcheck.Routes(app, healthcheck.Config{
		DB:    cfg.DB,
		Log:   cfg.Log,
		Build: cfg.Build,
	})
	usergrp.Routes(app, usergrp.Config{
		DB:   cfg.DB,
		Log:  cfg.Log,
		Auth: cfg.Auth,
	})
	authapi.Routes(app, authapi.Config{
		Auth: cfg.Auth,
	})
	articleapi.Routes(app, articleapi.Config{
		Log:         cfg.Log,
		DB:          cfg.DB,
		Auth:        cfg.Auth,
		ArticleCore: cfg.ArticleCore,
	})

}
