package tenant

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"external/sdk/client/tenant"
	"external/sdk/transport"

	"sender/internal/model"
)

type Registerer struct {
	registry               *Registry
	interval               time.Duration
	externalRegistryClient *tenant.Client
	metrics                *metrics
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
		metrics:                newMetrics(mr),
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
		r.refreshMetrics()
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
	tenantCfgs := toTenantConfigs(tenantCfgData...)
	r.registry.UpdateTenants(tenantCfgs...)
	return nil
}

func toTenantConfigs(cfgData ...tenant.Configuration) []model.TenantConfig {
	list := make([]model.TenantConfig, 0, len(cfgData))
	for _, cfg := range cfgData {
		list = append(list, model.TenantConfig{
			TenantID:  cfg.TenantID,
			BatchSize: cfg.TenantMaxBatchSize,
		})
	}
	return list
}

func (r *Registerer) refreshMetrics() {
	tenantsCount := r.registry.Count()
	r.metrics.tenants.Set(float64(tenantsCount))
}
