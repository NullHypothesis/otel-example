package lib

import (
	"log"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// NewTracerProvider returns a new provider for tracers. An application should
// call this function only once, as per the documentation:
// https://opentelemetry.io/docs/concepts/signals/traces/#tracer-provider
func NewTracerProvider(serviceName string) *sdktrace.TracerProvider {
	var (
		err error
		exp sdktrace.SpanExporter
		res *resource.Resource
	)

	// Set the service name in the resource.
	if res, err = resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	); err != nil {
		log.Fatal(err)
	}

	// Export traces to stdout for simplicity.
	if exp, err = stdouttrace.New(stdouttrace.WithPrettyPrint()); err != nil {
		log.Fatal(err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
}
