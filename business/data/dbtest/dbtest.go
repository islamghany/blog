package dbtest

import (
	"bytes"
	"context"
	"fmt"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/core/user/userdb"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/foundation/config"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

type Test struct {
	DB       *sqlx.DB
	Log      *logger.Logger
	CoreAPIs CoreAPIs
	TearDown func()
}

// New creates a new Test struct.
func NewTest(t *testing.T) *Test {
	// =========================================================================
	// Configuration
	type Config struct {
		// DB
		User         string `mapstructure:"DB_USER"`
		Password     string `mapstructure:"DB_PASSWORD" omit:"true"`
		Host         string `mapstructure:"DB_HOST"`
		Name         string `mapstructure:"DB_NAME"`
		MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS"`
		MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS"`
		DisabelTLS   bool   `mapstructure:"DISABLE_TLS"`
	}
	var cfg Config
	err := config.LoadConfig(&cfg, ".", "dev", "env")
	if err != nil {
		t.Fatalf("loading config: %v", err)
	}

	// =========================================================================
	// Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Name)
	db, err := db.Open(db.Config{
		Name:         cfg.Name,
		Host:         cfg.Host,
		User:         cfg.User,
		Password:     cfg.Password,
		MaxIdleConns: cfg.MaxIdleConns,
		MaxOpenConns: cfg.MaxOpenConns,
		DisabelTLS:   cfg.DisabelTLS,
	})
	if err != nil {
		t.Fatalf("opening db: %v", err)
	}

	// =========================================================================
	// // Handle Migration and Seed
	// if err := dbmigrate.MigrateUp(ctx, dsn); err != nil {
	// 	t.Fatalf("migrating up: %v", err)
	// }
	// if err := dbmigrate.SeedsWithSQLX(ctx, db); err != nil {
	// 	t.Fatalf("seeding data: %v", err)
	// }

	// =========================================================================
	// Log
	var buf bytes.Buffer
	log := logger.New(&buf, logger.LevelInfo, "TEST", func(context.Context) string { return web.GetTracerID(ctx) })

	// Core APIs
	coreAPIs := newCoreAPIs(log, db)

	t.Log("Ready for testing ...")

	// teardown is the function that should be invoked when the caller is done
	// with the database.
	teardown := func() {
		t.Helper()
		db.Close()

		fmt.Println("******************** LOGS ********************")
		fmt.Print(buf.String())
		fmt.Println("******************** LOGS ********************")
	}
	return &Test{
		DB:       db,
		Log:      log,
		CoreAPIs: coreAPIs,
		TearDown: teardown,
	}

}

// =============================================================================

// CoreAPIs represents all the core api's needed for testing.
type CoreAPIs struct {
	User *user.Core
}

func newCoreAPIs(log *logger.Logger, db *sqlx.DB) CoreAPIs {
	usrCore := user.NewCore(log, userdb.NewStore(log, db))

	return CoreAPIs{
		User: usrCore,
	}
}
