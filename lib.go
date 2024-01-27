package lib

import (
	"context"
	"log"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func newExporter(ctx context.Context) (*stdouttrace.Exporter, error) {
	// TODO: may have to change that
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

func NewTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	log.Printf("Schema URL: %v", semconv.SchemaURL)

	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("ExampleService"),
		),
	)
	if err != nil {
		log.Fatalf("Error merging resources: %v", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}
