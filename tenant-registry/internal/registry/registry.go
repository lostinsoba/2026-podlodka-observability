package registry

import (
	"tenant-registry/internal/config"
	"tenant-registry/internal/model"
)

type Registry struct {
	tenantCfgs []model.TenantConfiguration
}

func New(cfg config.TenantRegistry) *Registry {
	tenantCfgs := make([]model.TenantConfiguration, 0, len(cfg.TenantConfigurations.Tenants))
	for _, tenantCfg := range cfg.TenantConfigurations.Tenants {
		tenantCfgs = append(tenantCfgs, model.TenantConfiguration{
			TenantID:           tenantCfg.TenantID,
			TenantMaxBatchSize: tenantCfg.TenantMaxBatchSize,
		})
	}
	return &Registry{
		tenantCfgs: tenantCfgs,
	}
}
