package trace

import (
	"fmt"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type TracerType int

var (
	Default TracerType = 0
	Zipkin  TracerType = 1
)

// NewTracer initializes a new tracer.
func NewTracer(opts ...Option) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	options := NewOptions(opts...)

	if options.Name == "" {
		options.Name = fmt.Sprintf("tracer-%s", uuid.NewString())
	}

	switch options.TracerType {
	case Zipkin:
		if options.Address == "" {
			return nil, ErrTracerInvalidAddressInitialization
		}
		exporter = NewZipkinExporter(options.Address)
	default:
		exporter, _ = stdouttrace.New()
	}

	provider := trace.NewTracerProvider(
		trace.WithSampler(trace.ParentBased(trace.TraceIDRatioBased(options.FractionRatio))),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(options.Name),
		)),
	)

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return provider, nil
}
