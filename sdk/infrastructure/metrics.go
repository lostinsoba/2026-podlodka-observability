package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var buildInfoMetric = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "build_info",
	},
)

type MetricRegistry struct {
	Registry   *prometheus.Registry
	Registerer prometheus.Registerer
}

func NewMetricRegistry(info BuildInfo) (*MetricRegistry, error) {
	registry := prometheus.NewRegistry()
	labels := prometheus.Labels{
		"service":   info.Service,
		"component": info.Component,
		"version":   info.Version,
		"gitCommit": info.GitCommit,
	}
	registerer := prometheus.WrapRegistererWith(labels, registry)
	err := registerer.Register(buildInfoMetric)
	if err != nil {
		return nil, fmt.Errorf("failed to register build info metric: %w", err)
	}
	err = registerer.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	if err != nil {
		return nil, fmt.Errorf("failed to register process metrics: %w", err)
	}
	buildInfoMetric.Set(1)
	return &MetricRegistry{
		Registry:   registry,
		Registerer: registerer,
	}, nil
}

const (
	defaultMetricServerReadTimeout     = 5 * time.Second
	defaultMetricServerShutdownTimeout = 5 * time.Second
)

type MetricsServer struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func NewMetricServer(gatherer prometheus.Gatherer, port int, logger *slog.Logger) *MetricsServer {
	mux := http.NewServeMux()
	opts := promhttp.HandlerOpts{
		ErrorHandling: promhttp.ContinueOnError,
	}
	handler := promhttp.HandlerFor(gatherer, opts)
	mux.Handle("/metrics", handler)
	return &MetricsServer{
		httpServer: &http.Server{
			Addr:        fmt.Sprintf(":%d", port),
			Handler:     mux,
			ReadTimeout: defaultMetricServerReadTimeout,
		},
		logger: logger,
	}
}

func (s *MetricsServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("stopping metric server")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, defaultMetricServerShutdownTimeout)
		defer cancel()
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error("failed to gracefully stop server",
				slog.Any("error", err),
			)
		}
	}()
	s.logger.Info("starting metric server",
		slog.String("addr", s.httpServer.Addr),
	)
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("metric server stopped")
	return nil
}
