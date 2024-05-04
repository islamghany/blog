package mid

import (
	"context"
	"github/islamghany/blog/business/web/v1/metrics"
	"github/islamghany/blog/foundation/web"
	"net/http"
)

func Metrics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx = metrics.Set(ctx)
			err := handler(ctx, w, r)
			n := metrics.AddRequest(ctx)
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}
			if err != nil {
				metrics.AddError(ctx)
			}
			return err

		}
	}
}
