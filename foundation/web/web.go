package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
)

type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type App struct {
	*httprouter.Router
	shutdown chan os.Signal
	mw       []Middleware
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp(shutdown chan os.Signal, mw ...Middleware) *App {
	router := httprouter.New()

	return &App{
		Router:   router,
		shutdown: shutdown,
		mw:       mw,
	}
}

// SignalShutdown is used to gracefully shut down the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.shutdown <- syscall.SIGTERM
}

// Handle sets a handler function for a given HTTP method and path pair
// to the application server mux.
func (app *App) Handle(method, group, path string, handler Handler, mw ...Middleware) {
	handler = wrapMiddlewares(handler, mw)
	handler = wrapMiddlewares(handler, app.mw)
	app.handle(method, group, path, handler)
}

func (app *App) handle(method, group, path string, handler Handler) {

	h := func(w http.ResponseWriter, r *http.Request) {
		ctx := SetValues(r.Context(), &Values{
			TracerID: uuid.NewString(),
			Time:     time.Now().UTC(),
		})
		if err := handler(ctx, w, r); err != nil {
			log.Println(err)
		}
	}
	groupedPath := path
	if group != "" {
		groupedPath = "/" + group + path
	}
	app.HandlerFunc(method, groupedPath, h)
}
