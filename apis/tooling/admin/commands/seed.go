package commands

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/data/dbmigrate"
	db "github/islamghany/blog/business/data/dbsql/pgx"
)

func Seed(ctx context.Context, cfg db.Config) error {

	conn, err := db.Open(cfg)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = dbmigrate.SeedsWithSQLX(ctx, conn)
	if err != nil {
		return err
	}
	fmt.Println("Seeded database")
	return nil
}
