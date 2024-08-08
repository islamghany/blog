package mid

import (
	"context"
	"fmt"
	"github/islamghany/blog/business/web/v1/metrics"
	"github/islamghany/blog/foundation/web"
	"net/http"
	"time"
)

func Metrics() web.Middleware {
	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx1 := metrics.Set(ctx)
			ctx = metrics.SetProm(ctx1)
			val := web.GetValues(ctx)

			err := handler(ctx, w, r)

			n := metrics.AddRequest(ctx)
			if n%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}
			status := web.GetSetStatusCode(ctx)
			metrics.AddPromRequest(ctx, r.Method, r.URL.Path, fmt.Sprintf("%d", status))
			metrics.AddPromLatency(ctx, r.Method, r.URL.Path, float64(time.Since(val.Time).Milliseconds()))
			metrics.AddGoroutines(ctx)
			if err != nil {
				metrics.AddError(ctx)
			}
			return err

		}
	}
}
