package dto

import "tenant-registry/internal/model"

type TenantConfigurationData struct {
	TenantID           string `json:"tenant_id"`
	TenantMaxBatchSize int    `json:"tenant_max_batch_size"`
}

type TenantConfigurationsData struct {
	Tenants []TenantConfigurationData `json:"tenants"`
}

func ToTenantConfigurationsData(tenantCfgs []model.TenantConfiguration) TenantConfigurationsData {
	cfgs := make([]TenantConfigurationData, 0, len(tenantCfgs))
	for _, tenantCfg := range tenantCfgs {
		cfgs = append(cfgs, TenantConfigurationData{
			TenantID:           tenantCfg.TenantID,
			TenantMaxBatchSize: tenantCfg.TenantMaxBatchSize,
		})
	}
	return TenantConfigurationsData{
		Tenants: cfgs,
	}
}
