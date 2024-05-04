package healthcheck

import (
	"context"
	"errors"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"math/rand"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	DB  *sqlx.DB
	Log *logger.Logger
}

// New constructs a Handlers api for the health check.
func New(db *sqlx.DB, log *logger.Logger) *Handlers {
	return &Handlers{
		DB:  db,
		Log: log,
	}
}

func (h *Handlers) HealthCheckHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	rn := rand.Intn(10)

	if rn%2 == 0 {
		return response.NewError(errors.New("big error"), http.StatusBadRequest)
	}
	if rn%3 == 0 {
		panic("a panic happens")
	}
	data := struct {
		Message string `json:"message"`
		Status  int    `json:"status"`
	}{
		Message: "All things are working fine",
		Status:  200,
	}

	return web.Response(ctx, w, data, http.StatusOK)
}
