package mid

import (
	"context"
	"fmt"
	"github/islamghany/blog/foundation/logger"
	"github/islamghany/blog/foundation/web"
	"net/http"
	"time"
)

func Logger(log *logger.Logger) web.Middleware {
	return func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			val := web.GetValues(ctx)
			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}
			log.Info(ctx, "request started", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr)
			err := handler(ctx, w, r)
			log.Info(ctx, "request completed", "method", r.Method, "path", path, "remoteaddr", r.RemoteAddr, "statusCode", val.SetStatusCode, "since", time.Since(val.Time))
			return err
		}
		return h
	}

}
