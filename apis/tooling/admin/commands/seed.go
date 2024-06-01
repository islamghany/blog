package commands

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/data/dbmigrate"

	"github.com/jmoiron/sqlx"
)

func Seed(ctx context.Context, conn *sqlx.DB) error {

	err := dbmigrate.SeedsWithSQLX(ctx, conn)
	if err != nil {
		return err
	}
	fmt.Println("Seeded database")
	return nil
}
