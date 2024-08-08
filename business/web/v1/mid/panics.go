package mid

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/web/v1/metrics"
	"github/islamghany/blog/foundation/web"
	"net/http"
	"runtime/debug"
)

func Panics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
					metrics.AddPanics(ctx)
					metrics.AddPromRequest(ctx, r.Method, r.URL.Path, "600")
				}
			}()
			return handler(ctx, w, r)

		}
	}
}
