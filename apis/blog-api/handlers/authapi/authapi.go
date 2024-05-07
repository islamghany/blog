package authapi

import (
	"context"
	"github/islamghany/blog/apis/blog-api/handlers/usergrp"
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/web"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	Auth *auth.Auth
}

func NewAuthHandler(auth *auth.Auth) *AuthHandler {
	return &AuthHandler{
		Auth: auth,
	}
}

func (h *AuthHandler) Authorize(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var cred Credentials
	if err := web.Decode(w, r, &cred); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	usr, err := h.Auth.CoreUsr.QueryByUsername(ctx, cred.Username)
	if err != nil {
		return response.NewError(ErrInvalidCredentials, http.StatusUnauthorized)
	}
	// check password
	err = bcrypt.CompareHashAndPassword(usr.PasswordHashed, []byte(cred.Password))
	if err != nil {
		return response.NewError(ErrInvalidCredentials, http.StatusUnauthorized)
	}
	newV := usr.Version + 1
	token, _, err := h.Auth.Sign(usr.ID, usr.Roles, newV, time.Hour*24)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}
	var uu user.UpdateUser
	uu.Version = &newV
	u, err := h.Auth.CoreUsr.Update(ctx, usr, uu)
	type UserResp struct {
		User  usergrp.ApiUser `json:"user"`
		Token string          `json:"token"`
	}

	return web.Response(ctx, w, UserResp{
		User:  usergrp.ToApiUser(u),
		Token: token,
	}, http.StatusOK)
}
