package userdb

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github/islamghany/blog/business/core/user"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/business/data/order"
	"github/islamghany/blog/foundation/logger"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	log *logger.Logger
	DB  *sqlx.DB
}

func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		DB:  db,
	}
}

func (s *Store) Create(ctx context.Context, nu user.User) error {
	q := `
		INSERT INTO users
			(id, email, username, roles, first_name, last_name, password_hashed, created_at, updated_at, enabled)
		VALUES
		    (:id, :email, :username, :roles, :first_name, :last_name, :password_hashed, :created_at, :updated_at, :enabled)
	`
	dbusr := toDBUser(nu)
	if err := db.NamedExecContext(ctx, s.log, s.DB, q, dbusr); err != nil {
		if errors.Is(err, db.ErrDBDuplicate) {
			if strings.Contains(err.Error(), "users_username_key") {
				return fmt.Errorf("namedexeccontext: %w", user.ErrDuplicateUsername)
			}
			return fmt.Errorf("namedexeccontext: %w", user.ErrDuplicateEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

func (s *Store) QueryByID(ctx context.Context, id uuid.UUID) (user.User, error) {
	data := struct {
		ID uuid.UUID `db:"id"`
	}{
		ID: id,
	}
	q := `
		SELECT 
			id, email, username, roles, first_name, last_name, created_at, password_hashed, updated_at, enabled, version
		 FROM users
		WHERE id = :id
	`
	var dbusr dbuser
	if err := db.NamedQueryStruct(ctx, s.log, s.DB, q, data, &dbusr); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	if dbusr.ID == uuid.Nil {
		return user.User{}, user.ErrNotFound
	}

	return toCoreUser(dbusr), nil
}
func (s *Store) QueryByUsername(ctx context.Context, username string) (user.User, error) {
	data := struct {
		Username string `db:"username"`
	}{
		Username: username,
	}
	q := `
		SELECT 
			id, email, username, roles, first_name, last_name, created_at, password_hashed, updated_at, enabled, version
		 FROM users
		WHERE username = :username
	`
	var dbusr dbuser
	if err := db.NamedQueryStruct(ctx, s.log, s.DB, q, data, &dbusr); err != nil {
		if errors.Is(err, db.ErrDBNotFound) {
			return user.User{}, user.ErrNotFound
		}
		return user.User{}, fmt.Errorf("namedquerystruct: %w", err)
	}
	if dbusr.ID == uuid.Nil {
		return user.User{}, user.ErrNotFound
	}

	return toCoreUser(dbusr), nil
}

func (s *Store) Update(ctx context.Context, usr user.User) error {
	q := `
		UPDATE 
			users
		SET
			roles = :roles,
			first_name = :first_name,
			last_name = :last_name,
			password_hashed = :password_hashed,
			updated_at = :updated_at,
			enabled = :enabled,
			version = :version
		WHERE
			id = :id
	`
	dbusr := toDBUser(usr)
	if err := db.NamedExecContext(ctx, s.log, s.DB, q, dbusr); err != nil {
		if errors.Is(err, db.ErrDBDuplicate) {
			if strings.Contains(err.Error(), "users_username_key") {
				return fmt.Errorf("namedexeccontext: %w", user.ErrDuplicateUsername)
			}
			return fmt.Errorf("namedexeccontext: %w", user.ErrDuplicateEmail)
		}
		return fmt.Errorf("namedexeccontext: %w", err)
	}
	return nil
}

func (s *Store) Query(ctx context.Context, filter user.QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]user.User, int, error) {
	data := map[string]any{
		"offset": (pageNumber - 1) * rowsPerPage,
		"limit":  rowsPerPage,
	}
	q := `
	SELECT
	   count (*) OVER() AS total,
		*
	FROM
		users
	 
	`
	buf := bytes.NewBufferString(q)
	s.applyFilters(filter, data, buf)
	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, 0, err
	}
	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset LIMIT :limit")
	var dbusrs []dbuserWithTotal
	if err := db.NamedQuerySlice(ctx, s.log, s.DB, buf.String(), data, &dbusrs); err != nil {
		return nil, 0, fmt.Errorf("namedqueryslice: %w", err)
	}
	u, t := toCoreUserWithTotalSlice(dbusrs)
	return u, t, nil
}
