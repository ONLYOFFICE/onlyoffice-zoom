package wrapper

import (
	"context"
	"strings"

	"github.com/ONLYOFFICE/zoom-onlyoffice/pkg/middleware"
	"go-micro.dev/v4/metadata"
	"go-micro.dev/v4/server"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// TracePropagationHandlerWrapper wraps RPC handlers to trace rpc calls
func TracePropagationHandlerWrapper(hf server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		meta, _ := metadata.FromContext(ctx)
		converted := make(map[string]string)

		for k, v := range meta {
			converted[strings.ToLower(k)] = v
		}

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.MapCarrier(converted))

		ctx, span := otel.GetTracerProvider().Tracer(middleware.InstrumentationName).Start(ctx, req.Endpoint())
		defer span.End()

		if err := hf(ctx, req, rsp); err != nil {
			return err
		}

		return nil
	}
}
