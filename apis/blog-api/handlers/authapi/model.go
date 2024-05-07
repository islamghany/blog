package authapi

import (
	"errors"
	"github/islamghany/blog/foundation/validate"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Credentials struct {
	Username string `json:"username" validate:"required,min=6,max=32,alphanum"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

func (c *Credentials) Validate() error {
	return validate.Check(c)
}
