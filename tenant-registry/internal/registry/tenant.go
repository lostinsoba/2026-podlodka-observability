package registry

import "tenant-registry/internal/model"

func (r *Registry) ListTenantConfigurations() []model.TenantConfiguration {
	return r.tenantCfgs
}
