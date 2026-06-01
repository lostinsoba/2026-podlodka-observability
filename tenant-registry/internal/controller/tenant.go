package controller

import (
	"context"
	"log/slog"

	"external/sdk/middleware"

	"tenant-registry/internal/model"
)

func (c *Controller) ListTenantConfigurations(ctx context.Context) []model.TenantConfiguration {
	var (
		requestID = middleware.GetRequestID(ctx)
	)
	c.logger.Debug("list tenant configurations",
		slog.String("request_id", requestID),
	)
	return c.r.ListTenantConfigurations()
}
