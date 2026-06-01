package controller

import (
	"log/slog"

	"tenant-registry/internal/registry"
)

type Controller struct {
	r      *registry.Registry
	logger *slog.Logger
}

func New(r *registry.Registry, logger *slog.Logger) *Controller {
	return &Controller{r: r, logger: logger}
}
