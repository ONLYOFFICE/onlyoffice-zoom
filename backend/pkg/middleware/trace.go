package middleware

import (
	"net/http"

	"go-micro.dev/v4/metadata"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

const (
	InstrumentationName = "github.com/go-micro/plugins/v4/wrapper/trace/opentelemetry"
)

// TracePropagationMiddleware inject previous context into a new request
func TracePropagationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.MapCarrier{}
		propagator := otel.GetTextMapPropagator()
		ctx := propagator.Extract(r.Context(), carrier)

		propagator.Inject(ctx, carrier)
		next.ServeHTTP(w, r.WithContext(metadata.NewContext(ctx, metadata.Metadata(carrier))))
	})
}

// LogTraceMiddleware starts tracing with spans
func LogTraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		carrier := propagation.MapCarrier{}

		ctx, span := otel.GetTracerProvider().Tracer(InstrumentationName).Start(r.Context(), r.URL.Path)
		defer span.End()

		otel.GetTextMapPropagator().Inject(ctx, carrier)
		next.ServeHTTP(w, r.WithContext(metadata.NewContext(ctx, metadata.Metadata(carrier))))
	})
}
