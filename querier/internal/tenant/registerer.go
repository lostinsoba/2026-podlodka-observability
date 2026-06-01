package tenant

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"external/sdk/client/tenant"
	"external/sdk/transport"
)

type Registerer struct {
	registry               *Registry
	interval               time.Duration
	externalRegistryClient *tenant.Client
	logger                 *slog.Logger
}

type RegistererConfig struct {
	Interval time.Duration
	Endpoint string
}

func NewRegisterer(registry *Registry, cfg RegistererConfig, mr prometheus.Registerer, logger *slog.Logger) *Registerer {
	transportWrapper := transport.Wrapper(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewRequestMetricsRoundTripper(mr, transport.NewRequestLoggingRoundTripper(rt, logger), "tenantRegistry")
	})
	externalRegistryClient := tenant.New(
		cfg.Endpoint,
		tenant.OptionWithTransportWrapper(transportWrapper),
	)
	return &Registerer{
		registry:               registry,
		interval:               cfg.Interval,
		externalRegistryClient: externalRegistryClient,
		logger:                 logger,
	}
}

func (r *Registerer) Run(ctx context.Context) error {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		err := r.registerTenants(ctx)
		if err != nil {
			r.logger.Error("failed to register tenants",
				slog.Any("error", err),
			)
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
}

func (r *Registerer) registerTenants(ctx context.Context) error {
	tenantCfgData, err := r.externalRegistryClient.GetTenantConfigurations(ctx)
	if err != nil {
		return err
	}
	tenants := toTenants(tenantCfgData...)
	r.registry.UpdateTenants(tenants...)
	return nil
}

func toTenants(cfgData ...tenant.Configuration) []string {
	list := make([]string, 0, len(cfgData))
	for _, cfg := range cfgData {
		list = append(list, cfg.TenantID)
	}
	return list
}
