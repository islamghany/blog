package main

import (
	"context"
	"expvar"
	"fmt"
	"github/islamghany/blog/apis/blog-api/config"
	"github/islamghany/blog/apis/blog-api/handlers"
	"github/islamghany/blog/business/auth"
	"github/islamghany/blog/business/core/article"
	"github/islamghany/blog/business/core/article/articledb"
	"github/islamghany/blog/business/core/user"
	"github/islamghany/blog/business/core/user/userdb"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	v1 "github/islamghany/blog/business/web/v1"
	"github/islamghany/blog/business/web/v1/debug"
	"github/islamghany/blog/business/web/v1/mid"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	// "github/islamghany/blog/bussiness/v1/debug"
)

var build = "development"
var version = "0.0.1"

func main() {
	// used context
	ctx := context.Background()

	// ========================================================================================
	// Initialize the logger.
	var log *logger.Logger

	logsEvents := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "BLOG: Send an alert")
		},
	}

	tracerIDFunc := func(ctx context.Context) string {
		return web.GetTracerID(ctx)
	}

	var minLevel = logger.LevelInfo
	if build == "development" {
		minLevel = logger.LevelDebug
	}

	log = logger.NewWithEvents(os.Stdout, minLevel, "BLOG-SERVICE", tracerIDFunc, logsEvents)

	// ========================================================================================
	// GOMAXPROCS
	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "build", build)

	// ========================================================================================
	// Initialize the configuration.
	cfg, err := config.LoadConfig(true)
	if err != nil {
		log.Error(ctx, "startup", "msg", err)
		return
	}
	log.Info(ctx, "startup", "config", cfg.Parse())

	// ==========================================================================================
	// Run the application
	if err = run(ctx, log, &cfg); err != nil {
		log.Error(ctx, "startup", "run", err)
	}
}

func run(ctx context.Context, log *logger.Logger, cfg *config.Config) error {
	w := sync.WaitGroup{}
	w.Add(1)
	// ==========================================================================================
	// Starting the application
	log.Info(ctx, "starting", "version", build)
	defer log.Info(ctx, "shutdown complete")

	// ==========================================================================================
	// Starting the DB
	log.Info(ctx, "startup", "status", "db starting", "host", cfg.DBHost)
	db, err := db.Open(
		db.Config{
			Name:         cfg.DBName,
			Host:         cfg.DBHost,
			Password:     cfg.DBPassword,
			User:         cfg.DBUser,
			DisabelTLS:   cfg.DisabelTLS,
			MaxOpenConns: cfg.MaxOpenConns,
			MaxIdleConns: cfg.MaxIdleConns,
		},
	)
	if err != nil {
		return fmt.Errorf("db conection err: %w", err)
	}
	defer func() {
		log.Info(ctx, "shutdown", "status", "db closing", "host", cfg.DBHost)
		db.Close()
	}()
	err = db.Ping()
	if err != nil {
		fmt.Println("err", err)
	}
	// ==========================================================================================
	// Starting debug server
	go func() {
		defer w.Done()
		expvar.NewString("build").Set(build)
		expvar.NewString("version").Set(version)

		log.Info(ctx, "startup", "status", "debug v1 started", "host", cfg.DebugHost)
		if err := http.ListenAndServe(cfg.DebugHost, debug.DebugMux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router cloased", "host", cfg.DebugHost)
		}
	}()
	// ==========================================================================================
	// Auth
	userCore := user.NewCore(log, userdb.NewStore(log, db))
	articleCore := article.NewCore(log, articledb.NewStore(log, db), userCore)
	auth := auth.NewAuth(log, cfg.JWTSecret, userCore)

	// ==========================================================================================
	// Starting the main server
	webMux := v1.WebMux(&v1.WebMuxConfig{
		Log:         log,
		DB:          db,
		Whitelist:   cfg.WHITELIST,
		Auth:        auth,
		ArticleCore: articleCore,
	}, handlers.Routes{})
	srv := http.Server{
		Addr:         cfg.APIHost,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		ErrorLog:     logger.NewStandardLogger(log, logger.LevelError),
		IdleTimeout:  cfg.IdleTimeout,
		Handler:      mid.CORS(cfg.WHITELIST)(webMux),
	}
	shutdownError := make(chan error)
	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", srv.Addr)
		shutdownError <- srv.ListenAndServe()
	}()

	// handling gracefull shutdown
	// here we are using a buffered channel because signal.Notify does't wait for the receiver to be avaiable
	// when sending to the provided channel
	shutdown := make(chan os.Signal, 1)
	// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and
	// relay them to the quit channel. Any other signals will not be caught by
	// signal.Notify() and will retain their default behavior.
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-shutdownError:
		{
			return fmt.Errorf("server error: %w", err)
		}
	case sig := <-shutdown:
		{
			log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig.String())
			defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig.String())
			// create a context that carries the deadline.
			ctx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				srv.Close()
				return fmt.Errorf("could not stop server gracefully: %w", err)
			}
		}

	}
	return nil
}
