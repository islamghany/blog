package auth

import (
	"context"
	"errors"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/foundation/logger"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Different types of error returned by the VerifyToken function
var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type Payload struct {
	ID        uuid.UUID
	Roles     []user.Role
	IssuedAt  time.Time
	ExpiredAt time.Time
	Version   int
}

func NewPayload(id uuid.UUID, roles []user.Role, version int, duration time.Duration) *Payload {
	return &Payload{
		ID:        id,
		Roles:     roles,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Version:   version,
	}
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}

func (p *Payload) Outdated(currentVersion int) error {
	if p.Version != currentVersion {
		return ErrInvalidToken
	}
	return nil
}

func (p *Payload) HasRole(role user.Role) bool {
	for _, r := range p.Roles {
		if r == role {
			return true
		}
	}
	return false
}

/// Auth

type Auth struct {
	log       *logger.Logger
	secretKey string
	CoreUsr   *user.Core
}

func NewAuth(log *logger.Logger, secretKey string, coreUsr *user.Core) *Auth {
	return &Auth{
		log:       log,
		secretKey: secretKey,
		CoreUsr:   coreUsr,
	}
}

func (a *Auth) Authenticate(ctx context.Context, r *http.Request) (context.Context, error) {
	BearerToken := r.Header.Get("Authorization")
	if BearerToken == "" {
		return nil, ErrInvalidToken
	}
	t := strings.Split(BearerToken, "Bearer ")
	if len(t) != 2 {
		return nil, ErrInvalidToken
	}
	token := t[1]
	payload, err := a.Verify(token)
	if err != nil {
		return nil, err
	}
	if err := payload.Valid(); err != nil {
		return nil, err
	}
	user, err := a.CoreUsr.QueryByID(ctx, payload.ID)
	if err != nil {
		return nil, err
	}
	if err := payload.Outdated(user.Version); err != nil {
		return nil, err
	}
	newCtx := SetUser(ctx, &user)
	return newCtx, nil
}
