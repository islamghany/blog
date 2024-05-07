package mid

import (
	"context"
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

func Authen(auth *auth.Auth) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Check if the user is authenticated.
			err := auth.Authenticate(ctx, r)
			if err != nil {
				return response.NewError(err, http.StatusUnauthorized)
			}
			return handler(ctx, w, r)
		}
	}
}
