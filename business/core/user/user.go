package user

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/data/order"
	"github/islamghany/blog/foundation/logger"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// possible errors that  can happen when interacting with the user store.
var (
	ErrNotFound              = fmt.Errorf("user not found")
	ErrDuplicateEmail        = fmt.Errorf("email already in use")
	ErrDuplicateUsername     = fmt.Errorf("username already in use")
	ErrAuthenticationFailure = fmt.Errorf("authentication failed")
)

// Storer defines the database operations for working with users.
type Storer interface {
	Create(context context.Context, nu User) error
	QueryByID(context context.Context, id uuid.UUID) (User, error)
	Update(context context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, int, error)
	QueryByUsername(ctx context.Context, username string) (User, error)
	Count(ctx context.Context) (int, error)
}

// =============================================================================

// Core manages the set of APIs for user access.
type Core struct {
	log   *logger.Logger
	store Storer
}

// NewCore constructs a Core for user api access.
func NewCore(log *logger.Logger, store Storer) *Core {
	return &Core{
		log:   log,
		store: store,
	}
}

// Create inserts a new user into the database.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generating password hash: %w", err)
	}

	now := time.Now()
	usr := User{
		ID:             uuid.New(),
		Email:          nu.Email,
		Username:       nu.Username,
		Roles:          nu.Roles,
		FirstName:      nu.FirstName,
		LastName:       nu.LastName,
		PasswordHashed: passwordHashed,
		CreatedAt:      now,
		UpdatedAt:      now,
		Enabled:        true,
	}

	err = c.store.Create(ctx, usr)
	if err != nil {
		return User{}, fmt.Errorf("creating user: %w", err)
	}

	return usr, nil
}

// QueryByID returns the user with the specified ID.
func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (User, error) {
	usr, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("querying user by id: %w", err)
	}
	return usr, nil
}

// QueryByID returns the user with the specified ID.
func (c *Core) QueryByUsername(ctx context.Context, username string) (User, error) {
	usr, err := c.store.QueryByUsername(ctx, username)
	if err != nil {
		return User{}, fmt.Errorf("querying user by username: %w", err)
	}
	return usr, nil
}

// Update patches the user with the specified usr.
func (c *Core) Update(ctx context.Context, usr User, uu UpdateUser) (User, error) {
	if uu.FirstName != nil {
		usr.FirstName = *uu.FirstName
	}
	if uu.LastName != nil {
		usr.LastName = *uu.LastName
	}
	if uu.Password != nil {
		passwordHashed, err := bcrypt.GenerateFromPassword([]byte(*uu.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generating password hash: %w", err)
		}
		usr.PasswordHashed = passwordHashed
	}
	if uu.Roles != nil {
		usr.Roles = uu.Roles
	}
	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}
	if uu.Version != nil {
		usr.Version = *uu.Version
	}
	usr.UpdatedAt = time.Now()

	err := c.store.Update(ctx, usr)
	if err != nil {
		return User{}, fmt.Errorf("updating user: %w", err)
	}
	return usr, nil

}

// Query
func (c *Core) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, int, error) {
	users, total, err := c.store.Query(ctx, filter, orderBy, pageNumber, rowsPerPage)
	if err != nil {
		return nil, 0, fmt.Errorf("querying users: %w", err)
	}
	return users, total, nil
}

// Count
func (c *Core) Count(ctx context.Context) (int, error) {
	t, err := c.store.Count(ctx)
	if err != nil {
		return 0, err
	}

	return t, nil
}
