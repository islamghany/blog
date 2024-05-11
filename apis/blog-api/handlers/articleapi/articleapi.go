package articleapi

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/core/article"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

type ArticleHandler struct {
	log         *logger.Logger
	articleCore *article.Core
	userCore    *user.Core
}

func NewArticleHandler(log *logger.Logger, articleCore *article.Core, userCore *user.Core) *ArticleHandler {
	return &ArticleHandler{
		log:         log,
		articleCore: articleCore,
		userCore:    userCore,
	}
}

func (h *ArticleHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var na ApiNewArticle
	if err := web.Decode(w, r, &na); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	coreNewArticle := toNewArticleCore(na)

	id, err := h.articleCore.Create(ctx, coreNewArticle)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}

	w.Header().Set("Location", fmt.Sprintf("/v1/article/%d", id))

	return web.Response(ctx, w, nil, http.StatusCreated)
}

func (h *ArticleHandler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamID(r, "id")
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	art, err := h.articleCore.QueryByID(ctx, id)
	if err != nil {
		if err == article.ErrorNotFound {
			return response.NewError(err, http.StatusNotFound)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}

	return web.Response(ctx, w, toApiArticle(art), http.StatusOK)
}
