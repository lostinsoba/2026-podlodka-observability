package middleware

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func WithHandlerMetrics(mr prometheus.Registerer, handler http.Handler, handlerName string) http.Handler {
	var (
		reqCount = promauto.With(mr).NewCounterVec(
			prometheus.CounterOpts{
				Name: "requests_total",
				Help: "A counter for requests from the wrapped client",
				ConstLabels: prometheus.Labels{
					"handler": handlerName,
				},
			},
			[]string{"code", "method"},
		)
		reqDuration = promauto.With(mr).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "request_duration_seconds",
				Help:    "A histogram of latencies for requests",
				Buckets: prometheus.DefBuckets,
				ConstLabels: prometheus.Labels{
					"handler": handlerName,
				},
			},
			[]string{"code", "method"},
		)
	)
	return promhttp.InstrumentHandlerCounter(reqCount, promhttp.InstrumentHandlerDuration(reqDuration, handler))
}
