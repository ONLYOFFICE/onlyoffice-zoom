package middleware

import (
	"net/http"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware/cors"
	corsmiddleware "github.com/go-chi/cors"
)

// Cors creates a new CORS middleware.
func Cors(opts ...cors.Option) func(http.Handler) http.Handler {
	options := cors.NewOptions(opts...)
	return corsmiddleware.Handler(corsmiddleware.Options{
		AllowedOrigins:   options.AllowedOrigins,
		AllowedMethods:   options.AllowedMethods,
		AllowedHeaders:   options.AllowedHeaders,
		AllowCredentials: options.AllowCredentials,
	})
}
