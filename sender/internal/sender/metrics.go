package sender

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	messagesSent   *prometheus.CounterVec
	messagesFailed *prometheus.CounterVec
}

func newMetrics(registerer prometheus.Registerer) *metrics {
	return &metrics{
		messagesSent: promauto.With(registerer).NewCounterVec(prometheus.CounterOpts{
			Name: "messages_sent_count",
			Help: "Count of messages sent",
		}, []string{"tenant_id"}),
		messagesFailed: promauto.With(registerer).NewCounterVec(prometheus.CounterOpts{
			Name: "messages_failed_count",
			Help: "Count of messages failed to send",
		}, []string{"tenant_id"}),
	}
}
