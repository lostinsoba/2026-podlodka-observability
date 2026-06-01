package transport

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type RequestTracingRoundTripper struct {
	internal http.RoundTripper
	tracer   trace.Tracer
}

func NewRequestTracingRoundTripper(transport http.RoundTripper, tracer trace.Tracer) http.RoundTripper {
	return &RequestTracingRoundTripper{
		internal: transport,
		tracer:   tracer,
	}
}

func (rt *RequestTracingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := rt.tracer.Start(
		req.Context(),
		buildSpanName(req.Method, req.URL.Hostname(), req.URL.Path),
	)
	defer span.End()

	propagator := propagation.TraceContext{}
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	req = req.WithContext(ctx)

	return rt.internal.RoundTrip(req)
}

func buildSpanName(method, host, path string) string {
	return fmt.Sprintf("%s %s %s", method, host, path)
}
