package commands

import (
	"context"
	"github/islamghany/blog/business/data/dbmigrate"
	"github/islamghany/blog/foundation/logger"
)

func Migrate(ctx context.Context, log *logger.Logger, dsn string) error {
	err := dbmigrate.MigrateUp(ctx, dsn)
	if err != nil {
		return err
	}
	log.Info(ctx, "Migrated database")
	return nil
}
