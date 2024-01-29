package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	lib "github.com/NullHypothesis/otel-example"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func printHeaders(ctx context.Context, h http.Header) {
	// Create another child span.
	_, span := tracer.Start(ctx, "printHeaders")
	defer span.End()

	log.Print("HTTP request headers:")
	for key, value := range h {
		log.Printf("%v: %v", key, value)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// The 'Traceparent' HTTP header contains the client's serialized span
	// context. The server creates child spans from this context.
	spanCtx, span := tracer.Start(otel.GetTextMapPropagator().Extract(
		r.Context(),
		propagation.HeaderCarrier(r.Header),
	), "handleIndex")
	defer span.End()

	printHeaders(spanCtx, r.Header)

	fmt.Fprintln(w, "Hello world")
}

func main() {
	// Both the client and the server use a tracer provider singleton to create
	// new traces.
	var tp = lib.NewTracerProvider("HttpService")
	defer tp.Shutdown(context.Background())
	tracer = tp.Tracer("")

	otel.SetTextMapPropagator(propagation.TraceContext{})

	log.Print("Starting Web service.")
	http.ListenAndServe(":8080", http.HandlerFunc(handleIndex))
}
