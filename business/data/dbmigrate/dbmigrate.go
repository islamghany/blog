package dbmigrate

import (
	"context"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
)

var (
	//go:embed seeds/seeds.sql
	seeds string
	//go:embed migrations/*.sql
	migrations embed.FS
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
	// m, err := migrate.New(migrations, dsn)
	// if err != nil {
	// 	return err
	// }

	// if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	// 	return err
	// }

	// return nil
	d, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, dsn)
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
