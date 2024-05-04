package dbmigrate

import (
	"context"
	_ "embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed seeds/seeds.sql
	seeds string
)

func Seeds(ctx context.Context, db *pgx.Conn) error {
	return pgx.BeginFunc(ctx, db, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, seeds)
		return err
	})
}

func SeedsWithSQLX(ctx context.Context, db *sqlx.DB) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	_, err = tx.ExecContext(ctx, seeds)
	return err
}

func MigrateUp(ctx context.Context, dsn string) error {
	m, err := migrate.New("file://business/data/dbmigrate/migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil

}
func MigrateDown(ctx context.Context, dsn string) error {
	m, err := migrate.New("embed://migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
