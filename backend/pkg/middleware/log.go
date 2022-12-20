package middleware

import (
	"net/http"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/log"
)

// Log creates a new debug logging middleware.
func Log(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if logger != nil {
				logger.Debugf("calling [%s] route %s", r.Method, r.URL.Path)
			}
			next.ServeHTTP(w, r)
		})
	}
}
