package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	lib "github.com/NullHypothesis/otel-example"
)

var tracer trace.Tracer

func newExporter(ctx context.Context) (*stdouttrace.Exporter, error) {
	return stdouttrace.New(stdouttrace.WithPrettyPrint())
}

type TracingRoundTripper struct {
	origTripper http.RoundTripper
}

func (t TracingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	spanCtx, span := tracer.Start(r.Context(), "RoundTrip")
	//span.SetAttributes(attribute.String("endpoint", r.RequestURI))
	defer span.End()

	otel.GetTextMapPropagator().Inject(
		spanCtx,
		propagation.HeaderCarrier(r.Header),
	)
	//log.Printf("Fields: %v", otel.GetTextMapPropagator().Fields())

	return t.origTripper.RoundTrip(r.WithContext(spanCtx))
}

func main() {
	var (
		tp = lib.NewTracerProvider()
	)
	defer tp.Shutdown(context.Background())
	tracer = tp.Tracer("http-client")

	otel.SetTextMapPropagator(propagation.TraceContext{})

	client := http.Client{
		Transport: TracingRoundTripper{
			origTripper: &http.Transport{},
		},
	}
	resp, err := client.Get("http://localhost:5000")
	if err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
	log.Print(resp)
}
