package tenant

import (
	"sync"

	"sender/internal/model"
)

type Registry struct {
	tenants map[string]model.TenantConfig
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		tenants: make(map[string]model.TenantConfig),
	}
}

func (r *Registry) UpdateTenants(configs ...model.TenantConfig) {
	r.mu.Lock()
	defer r.mu.Unlock()
	clear(r.tenants)
	for _, config := range configs {
		r.tenants[config.TenantID] = config
	}
}

func (r *Registry) Lookup(tenant string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.tenants[tenant]
	return ok
}

func (r *Registry) ListTenantConfigs() []model.TenantConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tenants := make([]model.TenantConfig, 0, len(r.tenants))
	for _, tenant := range r.tenants {
		tenants = append(tenants, tenant)
	}
	return tenants
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tenants)
}
