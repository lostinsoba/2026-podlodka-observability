package processor

import (
	"context"
	"log/slog"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type metrics struct {
	messagesReceivedQueueSize     prometheus.Gauge
	messagesReceivedQueueCapacity prometheus.Gauge
	messagesUpdated               prometheus.Counter
	messagesReceived              *prometheus.CounterVec
	messagesStored                prometheus.Gauge
}

func newMetrics(registerer prometheus.Registerer) *metrics {
	return &metrics{
		messagesReceivedQueueSize: promauto.With(registerer).NewGauge(prometheus.GaugeOpts{
			Name: "message_received_queue_size",
			Help: "Number of currently queued messages",
		}),
		messagesReceivedQueueCapacity: promauto.With(registerer).NewGauge(prometheus.GaugeOpts{
			Name: "messages_received_queue_capacity",
			Help: "Max number of messages in queue",
		}),
		messagesUpdated: promauto.With(registerer).NewCounter(prometheus.CounterOpts{
			Name: "messages_updated_count",
			Help: "Count of messages sent",
		}),
		messagesReceived: promauto.With(registerer).NewCounterVec(prometheus.CounterOpts{
			Name: "messages_received",
			Help: "Count of messages received per tenant id",
		}, []string{"tenant_id"}),
		messagesStored: promauto.With(registerer).NewGauge(prometheus.GaugeOpts{
			Name: "messages_stored",
			Help: "Count of messages stored",
		}),
	}
}

func (mp *MessageProcessor) refreshMessagesReceivedQueueSizeMetric(ctx context.Context) {
	const (
		interval = 30 * time.Second
	)
	mp.logger.Debug("starting refreshing messages received queue size metric",
		slog.Duration("interval", interval),
	)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	mp.metrics.messagesReceivedQueueCapacity.Set(float64(cap(mp.queue)))
	for {
		mp.metrics.messagesReceivedQueueSize.Set(float64(len(mp.queue)))
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}

func (mp *MessageProcessor) refreshMessagesStoredMetric(ctx context.Context) {
	const (
		interval = time.Minute
	)
	mp.logger.Debug("starting refreshing message metrics",
		slog.Duration("interval", interval),
	)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		stats, err := mp.messageStatsFunc(ctx)
		if err != nil {
			mp.logger.Error("failed to refresh messages stored metric",
				slog.Any("error", err),
			)
		} else {
			mp.metrics.messagesStored.Set(float64(stats.MessageCount))
		}
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}
	}
}
