package articleapi

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/article"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/web/v1/paging"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/validate"
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
	usr := auth.GetUser(ctx)
	coreNewArticle := toNewArticleCore(na, usr.ID)

	id, err := h.articleCore.Create(ctx, coreNewArticle)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}

	w.Header().Set("Location", fmt.Sprintf("%d", id))

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

func (h *ArticleHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamID(r, "id")
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	ua := ApiUpdateArticle{}
	if err := web.Decode(w, r, &ua); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	art, err := h.articleCore.QueryByID(ctx, id)
	if err != nil {
		if err == article.ErrorNotFound {
			return response.NewError(err, http.StatusNotFound)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}

	err = h.articleCore.Update(ctx, art, toUpdateArticleCore(ua))
	if err != nil {
		if err == article.ErrorNotFound {
			return response.NewError(err, http.StatusNotFound)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}

	return web.Response(ctx, w, nil, http.StatusNoContent)
}

func (h *ArticleHandler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := paging.ParseRequest(r)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	values := r.URL.Query()
	if v := values.Get("search"); v == "" {
		return validate.NewFieldError("search", "must be a string")
	}
	arts, _, err := h.articleCore.Query(ctx, values.Get("search"), page.Number, page.Size)
	apiArtWithAuthor := toApiArticleWithAuthorSlice(arts)

	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}
	return web.Response(ctx, w, apiArtWithAuthor, http.StatusOK)

}
