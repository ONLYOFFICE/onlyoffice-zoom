package middleware

import (
	"net/http"
	"time"
)

func Timeout(timeout time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(http.HandlerFunc(next.ServeHTTP), timeout, "request timeout")
	}
}
