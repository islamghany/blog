package main

import (
	"context"
	"fmt"
	"github/islamghany/blog/apis/tooling/admin/commands"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/business/web/v1/connect"
	"github/islamghany/blog/foundation/logger"
	"net/url"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var build = "development"

type Config struct {
	DBUser       string `mapstructure:"DB_USER" default:"blog"`
	DBPassword   string `mapstructure:"DB_PASSWORD" omit:"true" default:"islamghany"`
	DBHost       string `mapstructure:"DB_HOST" default:"localhost:5432"`
	DBName       string `mapstructure:"DB_NAME" default:"blog_db"`
	MaxIdleConns int    `mapstructure:"MAX_IDLE_CONNS" default:"25"`
	MaxOpenConns int    `mapstructure:"MAX_OPEN_CONNS" default:"25"`
	DisabelTLS   bool   `mapstructure:"DISABLE_TLS" default:"true"`
}

func main() {
	// go run main.go migrate
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	var config Config
	viper.SetDefault("DB_USER", "blog")
	viper.SetDefault("DB_PASSWORD", "islamghany")
	viper.SetDefault("DB_HOST", "localhost:5432")
	viper.SetDefault("DB_NAME", "blog_db")
	viper.SetDefault("MAX_IDLE_CONNS", 25)
	viper.SetDefault("MAX_OPEN_CONNS", 25)
	viper.SetDefault("DISABLE_TLS", true)
	viper.AutomaticEnv()
	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("error unmarshalling config", err)
	}
	log := logger.New(os.Stdout, logger.LevelInfo, "ADMIN", func(ctx context.Context) string {
		return "00000000-0000-0000-0000-000000000000"
	})
	log.Info(context.Background(), "startup", "config", config)
	if err := run(log, &config, cmd); err != nil {
		log.Error(context.Background(), "startup", "error", err)
		os.Exit(1)
	}

}

func run(log *logger.Logger, config *Config, cmd string) error {
	if cmd == "" {
		return fmt.Errorf("no command provided")
	}
	dbConfig := db.Config{
		Host:         config.DBHost,
		User:         config.DBUser,
		Name:         config.DBName,
		Password:     config.DBPassword,
		MaxIdleConns: config.MaxIdleConns,
		MaxOpenConns: config.MaxOpenConns,
		DisabelTLS:   config.DisabelTLS,
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
	conn, err := connect.ConnectWithBackOff(ctx, log, "POSTGRESQL", func() (*sqlx.DB, error) {
		return db.Open(dbConfig)
	}, 10)
	if err != nil {
		return fmt.Errorf("error connecting to database: %w", err)
	}
	log.Info(ctx, "startup", "status", "db connected", "host", config.DBHost)
	defer conn.Close()
	cmds := strings.Split(cmd, ",")
	for _, cmd := range cmds {
		switch cmd {
		case "migrate":
			if err := commands.Migrate(ctx, log, DSN.String()); err != nil {
				return fmt.Errorf("error migrating database: %w", err)
			}
		case "seed":
			if err := commands.Seed(ctx, log, conn); err != nil {
				return fmt.Errorf("error seeding database: %w", err)
			}
		default:
			return fmt.Errorf("unknown command %s", cmd)
		}

	}

	return nil
}
