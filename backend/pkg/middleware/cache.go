package middleware

import (
	"net/http"
	"time"
)

// NoCache sets no-cache headers.
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		rw.Header().Set("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		rw.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

		next.ServeHTTP(rw, r)
	})
}
