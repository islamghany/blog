package usergrp

import (
	"context"
	"errors"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/web/v1/paging"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

type UserHandler struct {
	user *user.Core
}

func NewUserHandler(user *user.Core) *UserHandler {
	return &UserHandler{
		user: user,
	}
}

func (h *UserHandler) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var nu ApiNewUser
	if err := web.Decode(w, r, &nu); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	coreUser, err := nu.toCoreNewUser()
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	usr, err := h.user.Create(ctx, coreUser)
	if err != nil {
		if errors.Is(err, user.ErrDuplicateEmail) {
			return response.NewError(user.ErrDuplicateEmail, http.StatusConflict)
		}
		if errors.Is(err, user.ErrDuplicateUsername) {
			return response.NewError(user.ErrDuplicateUsername, http.StatusConflict)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}
	return web.Response(ctx, w, toApiUser(usr), http.StatusCreated)
}

func (h *UserHandler) QueryByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamUUID(r, "id")
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	usr, err := h.user.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return response.NewError(user.ErrNotFound, http.StatusNotFound)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}
	return web.Response(ctx, w, toApiUser(usr), http.StatusOK)
}

func (h *UserHandler) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id, err := web.ParamUUID(r, "id")
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	var uu ApiUpdateUser
	if err := web.Decode(w, r, &uu); err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	usr, err := h.user.QueryByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return response.NewError(user.ErrNotFound, http.StatusNotFound)
		}
		return response.NewError(err, http.StatusInternalServerError)
	}
	coreUser, err := uu.toCoreUpdateUser()

	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	usr, err = h.user.Update(ctx, usr, coreUser)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}
	return web.Response(ctx, w, toApiUser(usr), http.StatusOK)
}

func (h *UserHandler) Query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	page, err := paging.ParseRequest(r)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	filter, err := parseFilter(r)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}
	order, err := parseOrder(r)
	if err != nil {
		return response.NewError(err, http.StatusBadRequest)
	}

	users, t, err := h.user.Query(ctx, filter, order, page.Number, page.Size)
	if err != nil {
		return response.NewError(err, http.StatusInternalServerError)
	}
	items := make([]ApiUser, len(users))
	for i, usr := range users {
		items[i] = toApiUser(usr)
	}

	return web.Response(ctx, w, paging.NewResponse(items, t, page.Number, page.Size), http.StatusOK)
}
