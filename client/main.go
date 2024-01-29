package main

import (
	"context"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	lib "github.com/NullHypothesis/otel-example"
)

var tracer trace.Tracer

type TracingRoundTripper struct {
	origTripper http.RoundTripper
}

func (t TracingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	spanCtx, span := tracer.Start(r.Context(), "RoundTrip")
	defer span.End()

	// Propagate the span context to the server, which is also instrumented for
	// tracing.
	otel.GetTextMapPropagator().Inject(
		spanCtx,
		propagation.HeaderCarrier(r.Header),
	)

	return t.origTripper.RoundTrip(r.WithContext(spanCtx))
}

func main() {
	// Both the client and the server use a tracer provider singleton to create
	// new traces.
	var tp = lib.NewTracerProvider("HttpClient")
	defer tp.Shutdown(context.Background())
	tracer = tp.Tracer("")

	otel.SetTextMapPropagator(propagation.TraceContext{})

	client := http.Client{
		// Use a RoundTripper to inject the span context into the HTTP request.
		Transport: TracingRoundTripper{
			origTripper: &http.Transport{},
		},
	}
	if _, err := client.Get("http://localhost:8080"); err != nil {
		log.Fatalf("Error fetching page: %v", err)
	}
}
