package usergrp

import (
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/core/user/userdb"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log *logger.Logger
	DB  *sqlx.DB
}

func Routes(app *web.App, cfg Config) {
	version := "v1"

	userCore := user.NewCore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB))
	userHandler := NewUserHandler(userCore)

	app.Handle(http.MethodGet, version, "/user", userHandler.Query)
	app.Handle(http.MethodPost, version, "/user", userHandler.Create)
	app.Handle(http.MethodGet, version, "/user/:id", userHandler.QueryByID)
	app.Handle(http.MethodPatch, version, "/user/:id", userHandler.Update)

}
