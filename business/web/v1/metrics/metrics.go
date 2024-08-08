package metrics

import (
	"context"
	"expvar"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
)

type metrics struct {
	goroutines *expvar.Int // current running goroutines
	requests   *expvar.Int // total requests received.
	errors     *expvar.Int // total errros occurred.
	panics     *expvar.Int // total panics occurred.
}

type PromMetrics struct {
	// Goroutines is the current running goroutines.
	Goroutines prometheus.Gauge
	// Requests is the total requests received.
	Requests prometheus.CounterVec
	// Latency is the request latency.
	RequestLatency prometheus.HistogramVec
}

var m *metrics
var promMetrics *PromMetrics

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		errors:     expvar.NewInt("errors"),
		panics:     expvar.NewInt("panics"),
	}
	promMetrics = &PromMetrics{
		Goroutines: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "blog",
			Name:      "goroutines",
			Help:      "Current running goroutines.",
		}),
		Requests: *prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: "blog",
			Name:      "requests",
			Help:      "Total requests received.",
		}, []string{"method", "path", "code"}),

		RequestLatency: *prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "blog",
			Name:      "request_latency",
			Help:      "Request latency.",
			// the request is by milliseconds
			Buckets: []float64{10, 20, 30, 40, 50, 100, 200, 300, 400, 500, 1000, 2000, 3000, 4000, 5000},
		}, []string{"method", "path"}),
	}
}

type metricsCtxKey int

const key metricsCtxKey = 1
const promKey metricsCtxKey = 2

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

/// Prometheus metrics

func RegiserPromMetrics() {
	register := prometheus.DefaultRegisterer
	register.MustRegister(promMetrics.Goroutines)
	register.MustRegister(promMetrics.Requests)
	register.MustRegister(promMetrics.RequestLatency)
}

// Set the prometheus metrics data into the context
func SetProm(ctx context.Context) context.Context {
	return context.WithValue(ctx, promKey, promMetrics)
}

func AddPromRequest(ctx context.Context, method, path string, code string) {
	if v, ok := ctx.Value(promKey).(*PromMetrics); ok {
		v.Requests.WithLabelValues(method, path, code).Inc()
	}
}
func AddPromLatency(ctx context.Context, method, path string, latency float64) {
	if v, ok := ctx.Value(promKey).(*PromMetrics); ok {
		v.RequestLatency.WithLabelValues(method, path).Observe(latency)
	}
}

func AddPromGoroutines(ctx context.Context) {
	if v, ok := ctx.Value(promKey).(*PromMetrics); ok {
		g := int64(runtime.NumGoroutine())
		v.Goroutines.Set(float64(g))
	}
}
