package controller

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	methodDuration *prometheus.HistogramVec
}

func newMetrics(mr prometheus.Registerer) *metrics {
	return &metrics{
		methodDuration: promauto.With(mr).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "controller_method_duration_seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "caller_id"},
		),
	}
}
