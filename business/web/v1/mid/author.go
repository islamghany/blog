package mid

import (
	"context"
	"errors"
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

func Authorization(allowedRole []user.Role) web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			user := auth.GetUser(ctx)
			isValid := false
			for _, role := range allowedRole {
				roleStr := role.Name()
				for _, userRole := range user.Roles {
					userRoleStr := userRole.Name()
					if roleStr == userRoleStr {
						isValid = true
						break
					}
				}
			}
			if !isValid && len(allowedRole) > 0 {
				return response.NewError(errors.New("unauthorized endpoint"), http.StatusUnauthorized)
			}
			return handler(ctx, w, r)
		}
	}
}
