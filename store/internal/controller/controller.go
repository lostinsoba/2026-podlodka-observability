package controller

import (
	"context"
	"log/slog"

	"github.com/prometheus/client_golang/prometheus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"store/internal/database"
	"store/internal/processor"
	"store/internal/tenant"
)

type Controller struct {
	mp      *processor.MessageProcessor
	tr      *tenant.Registry
	d       database.Database
	metrics *metrics
	tracer  trace.Tracer
	logger  *slog.Logger
}

func New(
	d database.Database,
	mp *processor.MessageProcessor,
	tr *tenant.Registry,
	mr prometheus.Registerer,
	tp *sdktrace.TracerProvider,
	logger *slog.Logger,
) *Controller {
	return &Controller{
		d:       d,
		mp:      mp,
		tr:      tr,
		metrics: newMetrics(mr),
		tracer:  tp.Tracer("controller"),
		logger:  logger,
	}
}

type ctrlFunc func(ctx context.Context) error

type ctrlMiddleware func(ctx context.Context, ctrlMethod string, ctrlFunc ctrlFunc) error

type ctrlMiddlewareChain []ctrlMiddleware

func newCtrlMiddlewareChain(middlewares ...ctrlMiddleware) ctrlMiddlewareChain {
	return middlewares
}

func (c ctrlMiddlewareChain) apply(ctx context.Context, ctrlMethod string, f ctrlFunc) error {
	wrapped := f
	for i := len(c) - 1; i >= 0; i-- {
		mw := c[i]
		next := wrapped
		wrapped = func(ctx context.Context) error {
			return mw(ctx, ctrlMethod, next)
		}
	}
	return wrapped(ctx)
}
