package commands

import (
	"context"
	"github/islamghany/blog/business/data/dbmigrate"
	"github/islamghany/blog/foundation/logger"

	"github.com/jmoiron/sqlx"
)

func Seed(ctx context.Context, log *logger.Logger, conn *sqlx.DB) error {

	err := dbmigrate.SeedsWithSQLX(ctx, conn)
	if err != nil {
		return err
	}
	log.Info(ctx, "Seeded database")
	return nil
}
