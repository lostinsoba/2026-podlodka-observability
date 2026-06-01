package transport

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewRequestMetricsRoundTripper(mr prometheus.Registerer, next http.RoundTripper, clientName string) http.RoundTripper {
	var (
		reqCount = promauto.With(mr).NewCounterVec(
			prometheus.CounterOpts{
				Name:        "client_requests_total",
				Help:        "A counter for requests from the wrapped client",
				ConstLabels: map[string]string{"client": clientName},
			},
			[]string{"code", "method"},
		)
		reqDuration = promauto.With(mr).NewHistogramVec(
			prometheus.HistogramOpts{
				Name:        "client_request_duration_seconds",
				Help:        "A histogram of request latencies",
				Buckets:     prometheus.DefBuckets,
				ConstLabels: map[string]string{"client": clientName},
			},
			[]string{"code", "method"},
		)
	)
	return promhttp.InstrumentRoundTripperCounter(reqCount,
		promhttp.InstrumentRoundTripperDuration(reqDuration, next))
}
