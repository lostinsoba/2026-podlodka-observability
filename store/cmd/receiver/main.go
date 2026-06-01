package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"external/sdk/infrastructure"

	"store/internal/config"
	"store/internal/controller"
	"store/internal/database"
	"store/internal/processor"
	"store/internal/server"
	"store/internal/server/receiver/router"
	"store/internal/tenant"
)

var (
	service   = "store"
	component = "receiver"
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

	metricsSrv := infrastructure.NewMetricServer(metricRegistry.Registry, cfg.Infrastructure.MetricPort, logger)
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

	d, err := database.New(cfg.Database, metricRegistry.Registerer)
	if err != nil {
		logger.Error("failed to configure storage",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	mp := processor.New(
		d,
		metricRegistry.Registerer,
		logger,
		processor.OptionQueueCapacity(cfg.Receiver.ReceiverQueueCapacity),
		processor.OptionBatchSize(cfg.Receiver.ReceiverBatchSize),
		processor.OptionBatchProcessInterval(cfg.Receiver.ReceiverBatchProcessInterval),
	)
	g.Go(func() error {
		if e := mp.Start(ctx); e != nil {
			return fmt.Errorf("failed to start message processor: %w", e)
		}
		return nil
	})

	traceProviderConfig := infrastructure.TraceProviderConfig{
		Endpoint: cfg.Infrastructure.TraceExportEndpoint,
		BaseContext: func() context.Context {
			return ctx
		},
		MaxSpanQueueSize:       cfg.Infrastructure.TraceSpanMaxQueueSize,
		MaxSpanExportBatchSize: cfg.Infrastructure.TraceSpanMaxExportBatchSize,
		MaxSpanBatchTimeout:    cfg.Infrastructure.TraceSpanBatchTimeout,
		SampleRate:             cfg.Infrastructure.TraceSampleRate,
	}
	traceProvider, err := infrastructure.NewTraceProvider(traceProviderConfig, buildInfo)
	if err != nil {
		logger.Error("failed to init trace provider",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	ctrl := controller.New(d, mp, tenantRegistry, metricRegistry.Registerer, traceProvider.Provider, logger)

	rt := router.New(ctrl, metricRegistry.Registerer, logger)

	srv, err := server.New(rt.Route(), cfg.API.Port, logger)
	if err != nil {
		logger.Error("failed to create server",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	g.Go(func() error {
		if e := srv.Run(ctx); e != nil {
			return fmt.Errorf("failed to run server: %w", e)
		}
		return nil
	})

	messageProcessorShutdown := func(ctx context.Context) {
		if e := mp.Stop(ctx); e != nil {
			logger.Error("failed to stop message processor",
				slog.Any("error", e),
			)
		}
	}
	traceProviderShutdown := func(ctx context.Context) {
		if e := traceProvider.Shutdown(ctx); e != nil {
			logger.Error("failed to shutdown trace provider",
				slog.Any("error", e),
			)
		}
	}

	sc := newShutdownCallbacks(
		messageProcessorShutdown,
		traceProviderShutdown,
	)

	var statusCode int

	err = g.Wait()
	if err != nil {
		logger.Error("error while running receiver",
			slog.Any("error", err),
		)
		statusCode = 1
	}

	sc.shutdown()

	logger.Info("receiver stopped")
	os.Exit(statusCode)
}

const defaultShutdownTimeout = time.Minute

type shutdownCallback func(ctx context.Context)

type shutdownCallbacks []shutdownCallback

func newShutdownCallbacks(sc ...shutdownCallback) shutdownCallbacks {
	return sc
}

func (scs shutdownCallbacks) shutdown() {
	shutdownCtx := context.Background()
	shutdownCtx, cancel := context.WithTimeout(shutdownCtx, defaultShutdownTimeout)
	defer cancel()
	for _, sc := range scs {
		sc(shutdownCtx)
	}
}
