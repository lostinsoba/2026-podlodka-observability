package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"external/sdk/infrastructure"

	"querier/internal/config"
	"querier/internal/querier"
	"querier/internal/tenant"
)

var (
	service   = "querier"
	component = "querier"
	version   = "unknown"
	gitCommit = "unknown"
)

var buildInfo = infrastructure.BuildInfo{
	Service:   service,
	Component: component,
	Version:   version,
	GitCommit: gitCommit,
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	logger := infrastructure.NewLogger(buildInfo, cfg.Infrastructure.LogLevel)
	metricRegistry, err := infrastructure.NewMetricRegistry(buildInfo)
	if err != nil {
		logger.Error("failed to create metric registry",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	g, ctx := errgroup.WithContext(ctx)

	metricsSrv := infrastructure.NewMetricServer(
		metricRegistry.Registry,
		cfg.Infrastructure.MetricPort,
		logger,
	)
	g.Go(func() error {
		if e := metricsSrv.Run(ctx); e != nil {
			return fmt.Errorf("failed to run metric server: %w", e)
		}
		return nil
	})

	tenantRegistry := tenant.NewRegistry()
	tenantRegisterer := tenant.NewRegisterer(
		tenantRegistry,
		tenant.RegistererConfig{
			Interval: cfg.TenantRegistry.TenantRegistryCacheUpdateInterval,
			Endpoint: cfg.TenantRegistry.TenantRegistryEndpoint,
		},
		metricRegistry.Registerer,
		logger,
	)
	g.Go(func() error {
		if e := tenantRegisterer.Run(ctx); e != nil {
			return fmt.Errorf("failed to run tenant registerer: %w", e)
		}
		return nil
	})

	qrCfg := querier.Config{
		QueryEndpoint: cfg.Querier.QueryEndpoint,
		QueryInterval: cfg.Querier.QueryInterval,
	}
	qr := querier.New(qrCfg, tenantRegistry, metricRegistry.Registerer, logger)
	g.Go(func() error {
		if e := qr.Run(ctx); e != nil {
			return fmt.Errorf("failed to run querier: %w", e)
		}
		return nil
	})

	var statusCode int

	err = g.Wait()
	if err != nil {
		logger.Error("error while running querier",
			slog.Any("error", err),
		)
		statusCode = 1
	}

	logger.Info("querier stopped")
	os.Exit(statusCode)
}
