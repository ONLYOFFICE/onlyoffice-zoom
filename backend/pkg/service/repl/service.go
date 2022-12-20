package repl

import (
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware"
	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/cors"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// NewService Initializes repl service with options.
func NewService(opts ...Option) *http.Server {
	options := newOptions(opts...)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	if options.Debug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

	return &http.Server{
		Addr: options.Address,
		Handler: alice.New(
			chimiddleware.RealIP,
			middleware.NewRateLimiter(200, 1*time.Second, middleware.WithKeyFuncAll),
			chimiddleware.RequestID,
			middleware.Cors(
				cors.WithAllowCredentials(options.CORS.AllowCredentials),
				cors.WithAllowedHeaders(options.CORS.AllowedHeaders...),
				cors.WithAllowedMethods(options.CORS.AllowedMethods...),
				cors.WithAllowedOrigins(options.CORS.AllowedOrigins...),
			),
			middleware.Secure,
			middleware.NoCache,
			middleware.Version(strconv.Itoa(options.Version)),
		).Then(mux),
	}
}
