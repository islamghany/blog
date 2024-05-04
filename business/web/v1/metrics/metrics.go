package metrics

import (
	"context"
	"expvar"
	"runtime"
)

type metrics struct {
	goroutines *expvar.Int // current running goroutines
	requests   *expvar.Int // total requests received.
	errors     *expvar.Int // total errros occurred.
	panics     *expvar.Int // total panics occurred.
}

var m *metrics

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("current_running_goroutines"),
		requests:   expvar.NewInt("total_requests_received"),
		errors:     expvar.NewInt("total_errros_occurred"),
		panics:     expvar.NewInt("total_panics_occurred"),
	}
}

type metricsCtxKey int

const key metricsCtxKey = 1

// Set the metrics data into the context
func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

// AddGoroutines refreshes the goroutine metric every 100 requests.
func AddGoroutines(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value()%100 == 0 {
			g := int64(runtime.NumGoroutine())
			v.goroutines.Set(g)
			return g
		}
	}
	return 0
}

func AddRequest(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
		return v.requests.Value()
	}
	return 0
}

func AddError(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
		return v.errors.Value()
	}
	return 0
}

func AddPanics(ctx context.Context) int64 {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
		return v.panics.Value()
	}
	return 0
}
