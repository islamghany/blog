package main

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/data/dbmigrate"
	"github/islamghany/blog/foundation/config"
	"log"
	"net/url"

	"github.com/jackc/pgx/v5"
)

func main() {
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
		log.Fatal(err)
	}
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

	ctx := context.Background()
	dsn := DSN.String()
	conn, err := pgx.Connect(ctx, dsn)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	defer conn.Close(ctx)

	fmt.Println("Migrating database")
	err = dbmigrate.MigrateUp(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Database migrated")

	fmt.Println("Loading seeds")
	err = dbmigrate.Seeds(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Seeds loaded")
	// // fmt.Println("Seeds loaded")
	// // store, err := db.Open(db.Config{
	// // 	Name:         cfg.Name,
	// // 	Host:         cfg.Host,
	// // 	Password:     cfg.Password,
	// // 	DisabelTLS:   cfg.DisabelTLS,
	// // 	User:         cfg.User,
	// // 	MaxOpenConns: cfg.MaxOpenConns,
	// // 	MaxIdleConns: cfg.MaxIdleConns,
	// // });
	// // q := `
	// // 	SELECT * FROM users where id =:id and name =:name and role in (:roles);
	// // `
	// // args := struct {
	// // 	ID    int
	// // 	Name  string
	// // 	Roles []string
	// // }{
	// // 	ID:    1,
	// // 	Name:  "John; select * from users;",
	// // 	Roles: []string{"admin", "user"},
	// // }
	// // query := db.QueryString(q, args)
	// // fmt.Println(query)

	// type User struct {
	// 	ID    uuid.UUID `json:"id" validate:"required,uuid"`
	// 	Name  string    `json:"name" validate:"required,min=3,max=100"`
	// 	Email string    `json:"email" validate:"required,email"`
	// }

	// u := User{
	// 	ID:    uuid.New(),
	// 	Name:  "a",
	// 	Email: "ssss",
	// }

	// err := validate.Check(u)
	// if err != nil {
	// 	fmt.Printf("Error: %v\n", err)

	// }

}
