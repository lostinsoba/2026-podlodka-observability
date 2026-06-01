package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"external/sdk/infrastructure"

	"tenant-registry/internal/config"
	"tenant-registry/internal/controller"
	"tenant-registry/internal/registry"
	"tenant-registry/internal/server"
	"tenant-registry/internal/server/router"
)

var (
	service   = "tenant-registry"
	component = "tenant-registry"
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

	r := registry.New(cfg.TenantRegistry)

	ctrl := controller.New(r, logger)

	rt := router.New(ctrl, logger)

	srv, err := server.New(rt.Route(), cfg.API.Port, logger)
	if err != nil {
		logger.Error("failed to create server",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = srv.Run(ctx); err != nil {
		logger.Error("api stopped with error",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	logger.Info("api stopped")
}
