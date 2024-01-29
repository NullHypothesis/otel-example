# Distributed tracing using OpenTelemetry

The OpenTelemetry documentation is complex and often confusing.
This repository contains a simple example of distributed tracing in Go.

The setting is simple:
A Go HTTP client creates a trace and propagates its trace to an
HTTP server, which adds spans to the trace.

The example code is meant to be simple and
contains a number of comments for clarification.
