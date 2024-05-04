package db

import (
	"context"
	"database/sql"
	"fmt"
	"github/islamghany/blog/foundation/logger"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

var (
	ErrDBNotFound     = sql.ErrNoRows
	ErrDBDuplicate    = fmt.Errorf("duplicate key")
	ErrUndefinedTable = fmt.Errorf("undefined table")
)

type Config struct {
	Name         string
	Host         string
	Password     string
	DisabelTLS   bool
	User         string
	MaxOpenConns int
	MaxIdleConns int
}

func Open(cfg Config) (*sqlx.DB, error) {
	sslMode := "require"
	if cfg.DisabelTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	DSN := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}
	db, err := sqlx.Connect("pgx", DSN.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	return db, nil
}

func NamedExecContext(ctx context.Context, log *logger.Logger, db sqlx.ExtContext, query string, data any) error {
	q := QueryString(query, data)
	if _, ok := data.(struct{}); ok {
		log.Infoc(ctx, 5, "database.NamedExecContext", "query", q)
	} else {
		log.Infoc(ctx, 4, "database.NamedExecContext", "query", q)
	}

	if _, err := sqlx.NamedExecContext(ctx, db, query, data); err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			switch pgerr.Code {
			case uniqueViolation:
				// wrap pgerr with the duplicate key error
				return fmt.Errorf("namedexeccontext:%s : %w", pgerr.Error(), ErrDBDuplicate)
			case undefinedTable:
				return ErrUndefinedTable
			}
			return err
		}
	}
	return nil

}

// NamedQueryStruct is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type where field replacement is necessary.
func NamedQueryStruct(ctx context.Context, log *logger.Logger, db sqlx.ExtContext, query string, data any, dest any) error {
	return namedQueryStruct(ctx, log, db, query, data, dest, false)
}
func namedQueryStruct(ctx context.Context, log *logger.Logger, db sqlx.ExtContext, query string, data any, dest any, withIn bool) error {
	q := QueryString(query, data)
	log.Infoc(ctx, 5, "database.NamedQueryStruct", "query", q)
	var rows *sqlx.Rows
	var err error
	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			namedQ, args, err := sqlx.Named(query, data)
			if err != nil {
				return nil, err
			}
			query, args, err := sqlx.In(namedQ, args...)
			if err != nil {
				return nil, err
			}
			query = db.Rebind(query)
			return db.QueryxContext(ctx, query, args...)
		}()
	default:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, data)
	}

	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == uniqueViolation {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()
	if rows.Next() {
		if err := rows.StructScan(dest); err != nil {
			return err
		}
	}
	return nil
}

// NamedQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshalled into a slice where field replacement is
// necessary.
func NamedQuerySlice[T any](ctx context.Context, log *logger.Logger, db sqlx.ExtContext, query string, data any, dest *[]T) error {
	return namedQuerySlice(ctx, log, db, query, data, dest, false)
}

func namedQuerySlice[T any](ctx context.Context, log *logger.Logger, db sqlx.ExtContext, query string, args any, dest *[]T, withIn bool) error {
	q := QueryString(query, args)
	log.Infoc(ctx, 5, "database.NamedQuerySlice", "query", q)
	var rows *sqlx.Rows
	var err error
	switch withIn {
	case true:
		rows, err = func() (*sqlx.Rows, error) {
			query, args, err := sqlx.Named(query, args)
			if err != nil {
				return nil, err
			}
			query, args, err = sqlx.In(query, args...)
			if err != nil {
				return nil, err
			}
			return db.QueryxContext(ctx, db.Rebind(query), args...)

		}()
	case false:
		rows, err = sqlx.NamedQueryContext(ctx, db, query, args)

	}
	if err != nil {
		if pgError, ok := err.(*pgconn.PgError); ok && pgError.Code == uniqueViolation {
			return ErrUndefinedTable
		}
		return err
	}

	defer rows.Close()
	var slice []T
	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}
	*dest = slice
	return nil
}

// QueryString returns the query string with the values replaced, and pritty print
// used for debugging
func QueryString(query string, args any) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}
	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}
	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")
	return strings.Trim(query, " ")
}
