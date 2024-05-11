package articleapi

import (
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/article"
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

	app.Handle(http.MethodPost, version, "/article", articleHandler.Create)
	app.Handle(http.MethodGet, version, "/article/:id", articleHandler.QueryByID)
}
