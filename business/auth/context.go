package auth

import (
	"context"
	"github/islamghany/blog/business/core/user"
)

type authKey int

const key = authKey(1)

// SetUser set the user in the context
func SetUser(ctx context.Context, user *user.User) context.Context {
	return context.WithValue(ctx, key, user)
}

// User get the user from the context
func getUser(ctx context.Context) *user.User {
	usr, ok := ctx.Value(key).(*user.User)
	if !ok {
		return nil
	}
	return usr
}
