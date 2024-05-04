package usergrp

import (
	"fmt"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/foundation/validate"
	"time"
)

// ================================================================================================
// User

type ApiUser struct {
	ID           string   `json:"id"`
	Email        string   `json:"email"`
	Username     string   `json:"username"`
	FirstName    string   `json:"first_name"`
	LastName     string   `json:"last_name"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Enabled      bool     `json:"enabled"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

func toApiUser(usr user.User) ApiUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.String()
	}
	return ApiUser{
		ID:        usr.ID.String(),
		Email:     usr.Email,
		Username:  usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Roles:     roles,
		Enabled:   usr.Enabled,
		CreatedAt: usr.CreatedAt.Format(time.RFC3339),
		UpdatedAt: usr.UpdatedAt.Format(time.RFC3339),
	}
}

// ================================================================================================
// New User

type ApiNewUser struct {
	Email           string   `json:"email" validate:"required,email"`
	Username        string   `json:"username" validate:"required,min=6,max=32,alphanum"`
	FirstName       string   `json:"first_name" validate:"required,min=1,max=32"`
	LastName        string   `json:"last_name" validate:"required,min=1,max=32"`
	Password        string   `json:"password" validate:"required,min=6,max=32"`
	ConfrimPassword string   `json:"confirm_password" validate:"required,min=6,max=32,eqfield=Password"`
	Roles           []string `json:"roles" validate:"required"`
}

func (nu *ApiNewUser) toCoreNewUser() (user.NewUser, error) {
	roles := make([]user.Role, len(nu.Roles))
	var err error
	for i, role := range nu.Roles {
		roles[i], err = user.ParseRole(role)
		if err != nil {
			return user.NewUser{}, fmt.Errorf("parsing role: %w", err)
		}
	}

	return user.NewUser{
		Email:             nu.Email,
		Username:          nu.Username,
		FirstName:         nu.FirstName,
		LastName:          nu.LastName,
		Password:          nu.Password,
		Roles:             roles,
		ConfirmedPassword: nu.ConfrimPassword,
	}, nil
}

func (nu *ApiNewUser) Validate() error {
	return validate.Check(nu)
}

// ================================================================================================
// Update User

type ApiUpdateUser struct {
	FirstName       *string  `json:"first_name" validate:"omitempty,min=1,max=32"`
	LastName        *string  `json:"last_name" validate:"omitempty,min=1,max=32"`
	Roles           []string `json:"roles" validate:"omitempty"`
	Enabled         *bool    `json:"enabled"`
	Password        *string  `json:"password" validate:"omitempty,min=6,max=32"`
	ConfirmPassword *string  `json:"confirm_password" validate:"omitempty,min=6,max=32,eqfield=Password"`
}

func (uu *ApiUpdateUser) toCoreUpdateUser() (user.UpdateUser, error) {
	var roles []user.Role
	if len(uu.Roles) > 0 {
		roles = make([]user.Role, len(uu.Roles))
		var err error
		for i, role := range uu.Roles {
			roles[i], err = user.ParseRole(role)
			if err != nil {
				return user.UpdateUser{}, fmt.Errorf("parsing role: %w", err)
			}
		}
	}
	return user.UpdateUser{
		FirstName:         uu.FirstName,
		LastName:          uu.LastName,
		Roles:             roles,
		Enabled:           uu.Enabled,
		Password:          uu.Password,
		ConfirmedPassword: uu.ConfirmPassword,
	}, nil

}

func (uu *ApiUpdateUser) Validate() error {
	return validate.Check(uu)
}

// ================================================================================================
// User Query

type ApiUserQuery struct {
	Page          int    `json:"page" validate:"omitempty,gte=1"`
	PageSize      int    `json:"page_size" validate:"omitempty,gte=1,lte=100"`
	SortBy        string `json:"sort_by" validate:"omitempty"`
	SortDirection string `json:"sort_direction" validate:"omitempty"`
}
