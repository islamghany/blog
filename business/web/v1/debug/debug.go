package debug

import (
	"expvar"
	"github/islamghany/blog/business/web/v1/metrics"
	"net/http"

	// package provides handlers for exposing runtime profiling data, which is useful for performance analysis.
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Mux registers all the debug routes from the standard library into a new mux
// bypassing the use of the DefaultServerMux. Using the DefaultServerMux would
// be a security risk since a dependency could inject a handler into our service
// without us knowing it.
func DebugMux() http.Handler {
	mux := http.NewServeMux()
	metrics.RegiserPromMetrics()
	mux.Handle("/metrics", promhttp.Handler())            // Exposes the registered metrics via HTTP.
	mux.HandleFunc("/debug/pprof/", pprof.Index)          // Index page for pprof profiles.
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline) // Current command line invocation.
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile) // CPU profile.
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)   // Symbol table.
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)     // A trace of execution of the current program.
	mux.Handle("/debug/vars", expvar.Handler())           // exposes the exported variables registered with expvar, allowing you to monitor the internal state of the application.
	return mux
}
