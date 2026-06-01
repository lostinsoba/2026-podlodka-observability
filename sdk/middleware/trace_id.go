package middleware

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	traceIDHeaderName = "X-Trace-ID"
)

func WithTraceID(tracer trace.Tracer, next http.Handler) http.Handler {
	propagator := propagation.TraceContext{}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var (
			requestID = GetRequestID(request.Context())
			callerID  = GetCallerID(request.Context())
		)
		ctx := propagator.Extract(
			request.Context(),
			propagation.HeaderCarrier(request.Header),
		)
		spanName := buildSpanName(request.Method, request.URL.Path)
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithAttributes(
				attribute.String("caller_id", callerID),
				attribute.String("request_id", requestID),
			),
		)
		defer span.End()
		traceID := span.SpanContext().TraceID().String()
		writer.Header().Set(traceIDHeaderName, traceID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func buildSpanName(method, path string) string {
	return fmt.Sprintf("%s %s", method, path)
}
