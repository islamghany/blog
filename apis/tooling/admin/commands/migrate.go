package commands

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/data/dbmigrate"
)

func Migrate(ctx context.Context, dsn string) error {
	err := dbmigrate.MigrateUp(ctx, dsn)
	if err != nil {
		return err
	}
	fmt.Println("Migrated database")
	return nil
}
