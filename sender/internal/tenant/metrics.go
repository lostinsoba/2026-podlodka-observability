package tenant

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	tenants prometheus.Gauge
}

func newMetrics(registerer prometheus.Registerer) *metrics {
	return &metrics{
		tenants: promauto.With(registerer).NewGauge(prometheus.GaugeOpts{
			Name: "registry_tenants",
		}),
	}
}
