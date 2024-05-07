package auth

import (
	"github/islamghany/blog/business/core/user"
	"time"

	"github.com/google/uuid"
)

type Maker interface {
	Sign(id uuid.UUID, roles []user.Role, version int, duration time.Duration) (string, *Payload, error)
	Verify(token string) (*Payload, error)
}

