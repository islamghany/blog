package main

import (
	"context"
	"flag"
	"fmt"
	"github/islamghany/blog/apis/tooling/admin/commands"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	cmd               string
	db_name           string
	db_pass           string
	db_user           string
	db_host           string
	db_max_idle_conns int
	db_max_open_conns int
	db_disable_tls    bool
}

func main() {
	var config Config
	flag.StringVar(&config.db_name, "db_name", "blog_db", "Database name")
	flag.StringVar(&config.db_pass, "db_pass", "islamghany", "Database password")
	flag.StringVar(&config.db_user, "db_user", "blog", "Database user")
	flag.StringVar(&config.db_host, "db_host", "localhost:5432", "Database host")
	flag.IntVar(&config.db_max_idle_conns, "db_max_idle_conns", 10, "Database max idle connections")
	flag.IntVar(&config.db_max_open_conns, "db_max_open_conns", 10, "Database max open connections")
	flag.BoolVar(&config.db_disable_tls, "db_disable_tls", true, "Disable TLS for database connection")
	flag.StringVar(&config.cmd, "cmd", "", "Command to run")
	flag.Parse()

	// log := logger.New(io.Discard, logger.LevelInfo, "ADMIN", func(ctx context.Context) string {
	// 	return "00000000-0000-0000-0000-000000000000"
	// })

	if err := run(&config); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

}

func run(config *Config) error {
	if config.cmd == "" {
		return fmt.Errorf("no command provided")
	}
	dbConfig := db.Config{
		Host:         config.db_host,
		User:         config.db_user,
		Name:         config.db_name,
		Password:     config.db_pass,
		MaxIdleConns: config.db_max_idle_conns,
		MaxOpenConns: config.db_max_open_conns,
		DisabelTLS:   config.db_disable_tls,
	}

	q := make(url.Values)
	ssl := "require"
	if dbConfig.DisabelTLS {
		ssl = "disable"
	}
	q.Set("sslmode", ssl)
	q.Set("timezone", "utc")
	DSN := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbConfig.User, dbConfig.Password),
		Host:     dbConfig.Host,
		Path:     dbConfig.Name,
		RawQuery: q.Encode(),
	}
	ctx := context.Background()
	cmds := strings.Split(config.cmd, ",")
	for _, cmd := range cmds {
		switch cmd {
		case "migrate":
			if err := commands.Migrate(ctx, DSN.String()); err != nil {
				return fmt.Errorf("error migrating database: %w", err)
			}
		case "seed":
			if err := commands.Seed(ctx, dbConfig); err != nil {
				return fmt.Errorf("error seeding database: %w", err)
			}
		default:
			return fmt.Errorf("unknown command %s", config.cmd)
		}

	}

	return nil
}
