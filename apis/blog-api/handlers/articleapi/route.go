package articleapi

import (
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/article"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/business/web/v1/mid"
	"net/http"

	// "github/islamghany/blog/business/core/user"
	// "github/islamghany/blog/business/core/user/userdb"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Log         *logger.Logger
	DB          *sqlx.DB
	Auth        *auth.Auth
	ArticleCore *article.Core
}

func Routes(app *web.App, cfg Config) {
	version := "v1"
	articleHandler := NewArticleHandler(cfg.Log, cfg.ArticleCore, cfg.Auth.CoreUsr)
	tran := mid.ExecuteInTransaction(cfg.Log, db.NewBeginner(cfg.DB))
	app.Handle(http.MethodPost, version, "/article", articleHandler.Create, mid.Authen(cfg.Auth))
	app.Handle(http.MethodPost, version, "/articlewithtran", articleHandler.CreateWithTran, mid.Authen(cfg.Auth), tran)
	app.Handle(http.MethodGet, version, "/article/:id", articleHandler.QueryByID)
	app.Handle(http.MethodPatch, version, "/article/:id", articleHandler.Update, mid.Authen(cfg.Auth))
	app.Handle(http.MethodGet, version, "/article", articleHandler.Query)
}
