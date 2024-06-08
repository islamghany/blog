package healthcheck

import (
	"context"
	"errors"
	db "github/islamghany/blog/business/data/dbsql/pgx"
	"github/islamghany/blog/business/web/v1/response"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
)

type Handlers struct {
	DB    *sqlx.DB
	Log   *logger.Logger
	Build string
}

// New constructs a Handlers api for the health check.
func New(db *sqlx.DB, log *logger.Logger, build string) *Handlers {
	return &Handlers{
		DB:    db,
		Log:   log,
		Build: build,
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

// LivenessHandler returns a simple status message.
// we use this handler to check if the service is alive. It is used by Kubernetes.
// if the service is not alive, Kubernetes will restart/kill the pod.
func (h *Handlers) LivenessHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	info := make(map[string]interface{})
	host, err := os.Hostname()
	if err != nil {
		host = "unknown"
	}
	info["hostname"] = host
	info["build"] = h.Build
	info["name"] = os.Getenv("KUBERNETES_NAME")
	info["namespace"] = os.Getenv("KUBERNETES_NAMESPACE")
	info["pod_ip"] = os.Getenv("KUBERNETES_POD_IP")
	info["node"] = os.Getenv("KUBERNETES_NODE_NAME")
	info["GOMAXPROCS"] = runtime.GOMAXPROCS(0)

	return web.Response(ctx, w, info, http.StatusOK)
}

// ReadinessHandler returns a simple status message.
// we use this handler to check if the service is ready to accept requests.
// if the service is not ready, Kubernetes will not send traffic to the pod.
func (h *Handlers) ReadinessHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if err := db.StatusCheck(ctx, h.DB); err != nil {
		h.Log.Info(ctx, "readiness check failed", "error", err)
		return response.NewError(err, http.StatusInternalServerError)
	}
	return web.Response(ctx, w, "OK", http.StatusOK)
}
