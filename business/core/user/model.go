package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Email          string
	Username       string
	Roles          []Role
	FirstName      string
	LastName       string
	PasswordHashed []byte
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Enabled        bool
	Version        int
}

type NewUser struct {
	Email             string
	Username          string
	FirstName         string
	LastName          string
	Password          string
	Roles             []Role
	ConfirmedPassword string
}

type UpdateUser struct {
	FirstName         *string
	LastName          *string
	Password          *string
	Roles             []Role
	ConfirmedPassword *string
	Enabled           *bool
	Version           *int
}
