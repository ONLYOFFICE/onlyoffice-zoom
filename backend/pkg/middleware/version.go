package middleware

import (
	"net/http"
)

// Version creates a new X-Version header middleware.
func Version(version string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("X-Api-Version", version)
			next.ServeHTTP(rw, r)
		})
	}
}
