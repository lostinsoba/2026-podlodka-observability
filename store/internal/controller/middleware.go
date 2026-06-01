package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"external/sdk/middleware"
)

func (c *Controller) withDurationMeasure(ctx context.Context, method string, f ctrlFunc) error {
	var (
		requestID = middleware.GetRequestID(ctx)
		callerID  = middleware.GetCallerID(ctx)
	)
	started := time.Now()
	err := f(ctx)
	duration := time.Since(started)
	c.metrics.methodDuration.WithLabelValues(method, callerID).Observe(duration.Seconds())
	c.logger.Debug("measured controller method durations",
		slog.String("method", method),
		slog.String("caller_id", callerID),
		slog.String("request_id", requestID),
		slog.Duration("duration", duration),
	)
	return err
}

func (c *Controller) withTracing(ctx context.Context, method string, f ctrlFunc) error {
	var (
		requestID = middleware.GetRequestID(ctx)
		callerID  = middleware.GetCallerID(ctx)
	)
	spanName := buildSpanName(method)
	ctx, span := c.tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("method", method),
			attribute.String("caller_id", callerID),
			attribute.String("request_id", requestID),
		),
	)
	err := f(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "")
	}
	span.End()
	return err
}

func buildSpanName(method string) string {
	return fmt.Sprintf("ctrl.%s", method)
}
