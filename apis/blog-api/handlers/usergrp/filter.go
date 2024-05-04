package usergrp

import (
	"fmt"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/core/user/userdb"
	"github/islamghany/blog/business/data/order"
	"github/islamghany/blog/foundation/validate"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func parseFilter(r *http.Request) (user.QueryFilter, error) {
	// Create the default filter.
	filter := user.QueryFilter{}

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		uid, err := uuid.Parse(userID)
		if err != nil {
			return filter, validate.NewFieldError("user_id", "must be a valid UUID")
		}
		filter.WithUserID(uid)
	}
	if username := r.URL.Query().Get("username"); username != "" {
		filter.WithUsername(username)
	}
	if email := r.URL.Query().Get("email"); email != "" {
		filter.WithEmail(email)
	}
	if firstName := r.URL.Query().Get("first_name"); firstName != "" {
		filter.WithFirstName(firstName)
	}
	if lastName := r.URL.Query().Get("last_name"); lastName != "" {
		filter.WithLastName(lastName)
	}
	if startDate := r.URL.Query().Get("start_date"); startDate != "" {
		t, err := time.Parse(time.RFC3339, startDate)
		if err != nil {
			return filter, validate.NewFieldError("start_date", err.Error())
		}
		filter.WithCreatedAt(t)
	}
	if endDate := r.URL.Query().Get("end_date"); endDate != "" {
		t, err := time.Parse(time.RFC3339, endDate)
		if err != nil {
			return filter, validate.NewFieldError("end_date", err.Error())
		}
		filter.WithUpdatedAt(t)
	}

	if err := filter.Validate(); err != nil {
		return filter, fmt.Errorf("validating filter: %w", err)
	}

	return filter, nil

}

/// parsing Order

func parseOrder(r *http.Request) (order.By, error) {
	orderBy, err := order.Parse(r, user.DefaultOrderBy)
	if err != nil {
		return order.By{}, err
	}
	if _, exits := userdb.OrderByFields[orderBy.Field]; !exits {
		return order.By{}, validate.NewFieldError("order", "invalid field")
	}
	return orderBy, nil
}
