package infrastructure

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

type TraceProviderConfig struct {
	Endpoint               string
	BaseContext            func() context.Context
	MaxSpanQueueSize       int
	MaxSpanExportBatchSize int
	MaxSpanBatchTimeout    time.Duration
	SampleRate             float64
}

type TraceProvider struct {
	Provider *sdktrace.TracerProvider
}

func NewTraceProvider(cfg TraceProviderConfig, info BuildInfo) (*TraceProvider, error) {
	exporter, err := otlptracehttp.New(
		cfg.BaseContext(),
		otlptracehttp.WithEndpoint(cfg.Endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(info.Service),
			semconv.ServiceVersion(info.Version),
			attribute.String("service.gitCommit", info.GitCommit),
			attribute.String("service.component", info.Component),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}
	sampler := sdktrace.ParentBased(
		sdktrace.TraceIDRatioBased(
			cfg.SampleRate,
		),
	)
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithBatcher(
			exporter,
			sdktrace.WithMaxQueueSize(cfg.MaxSpanQueueSize),
			sdktrace.WithMaxExportBatchSize(cfg.MaxSpanExportBatchSize),
			sdktrace.WithBatchTimeout(cfg.MaxSpanBatchTimeout),
		),
		sdktrace.WithResource(res),
	)
	return &TraceProvider{
		Provider: provider,
	}, nil
}

func (t *TraceProvider) Shutdown(ctx context.Context) error {
	return t.Provider.Shutdown(ctx)
}
