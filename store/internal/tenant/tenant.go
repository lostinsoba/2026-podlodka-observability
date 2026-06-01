package tenant

import (
	"sync"
)

type Registry struct {
	tenants map[string]struct{}
	mu      sync.RWMutex
}

func NewRegistry() *Registry {
	return &Registry{
		tenants: make(map[string]struct{}),
	}
}

func (r *Registry) UpdateTenants(tenants ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	clear(r.tenants)
	for _, tenant := range tenants {
		r.tenants[tenant] = struct{}{}
	}
}

func (r *Registry) Lookup(tenant string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.tenants[tenant]
	return ok
}

func (r *Registry) ListTenants() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tenants := make([]string, 0, len(r.tenants))
	for tenant := range r.tenants {
		tenants = append(tenants, tenant)
	}
	return tenants
}

func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tenants)
}
