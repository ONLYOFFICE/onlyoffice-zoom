package trace

import (
	"log"

	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/trace"
)

// NewZipkingExporter creates a new zipkin exporter.
func NewZipkinExporter(url string, opts ...zipkin.Option) trace.SpanExporter {
	exporter, err := zipkin.New(url, opts...)
	if err != nil {
		log.Fatalln("could not initialize a new zipkin exporter: ", err)
	}
	return exporter
}
