package processor

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"store/internal/database"
	"store/internal/model"
)

var (
	ErrProcessorStopped = errors.New("processor has been stopped")
)

type MessageProcessor struct {
	batchSize            int
	batchProcessInterval time.Duration
	batchProcessFunc     BatchProcessFunc

	messageStatsFunc MessageStatsFunc

	queue chan model.Message

	stop    chan struct{}
	stopped atomic.Bool

	metrics *metrics
	logger  *slog.Logger
}

func New(d database.Database, mr prometheus.Registerer, logger *slog.Logger, options ...Option) *MessageProcessor {
	batchProcessFunc := func(ctx context.Context, messages map[string]model.Message) (succeed bool, updatedCount int64) {
		var (
			resp      model.UpdateMessagesResponse
			updateErr error
		)
		resp, updateErr = d.UpdateMessages(ctx, model.UpdateMessagesRequest{Messages: messages})
		if updateErr != nil {
			logger.Error("failed to update messages",
				slog.Any("error", updateErr),
			)
			return false, 0
		}
		return true, resp.MessagesUpdated
	}
	messageStatsFunc := func(ctx context.Context) (stats model.MessagesStats, err error) {
		return d.EvaluateMessagesStats(ctx)
	}
	mp := &MessageProcessor{
		batchSize:            defaultBatchSize,
		batchProcessInterval: defaultBatchProcessInterval,
		batchProcessFunc:     batchProcessFunc,
		messageStatsFunc:     messageStatsFunc,
		queue:                make(chan model.Message, defaultQueueCapacity),
		stop:                 make(chan struct{}),
		metrics:              newMetrics(mr),
		logger:               logger,
	}
	for _, option := range options {
		option(mp)
	}
	return mp
}

func (mp *MessageProcessor) Queue(ctx context.Context, messages ...model.Message) error {
	if mp.isStopped() {
		return ErrProcessorStopped
	}
	for _, message := range messages {
		mp.queue <- message
		mp.metrics.messagesReceived.WithLabelValues(message.TenantID).Inc()
	}
	return nil
}

func (mp *MessageProcessor) Start(ctx context.Context) error {
	mp.logger.Debug("starting message processor",
		slog.Int("queue capacity", cap(mp.queue)),
		slog.Int("batch size", mp.batchSize),
		slog.Duration("batch process interval", mp.batchProcessInterval),
	)

	go mp.refreshMessagesReceivedQueueSizeMetric(ctx)
	go mp.refreshMessagesStoredMetric(ctx)

	ticker := time.NewTicker(mp.batchProcessInterval)
	defer ticker.Stop()

	messages := map[string]model.Message{}

	onStop := func() {
		for event := range mp.queue {
			messages[event.GetUniqueKey()] = event
		}
		ctx = context.Background()
		mp.batchProcessFunc(ctx, messages)
		close(mp.stop)
	}
	defer onStop()

	for {
		var queue <-chan model.Message
		if len(messages) < mp.batchSize {
			queue = mp.queue
		}
		select {
		case message, ok := <-queue:
			if !ok {
				return nil
			}
			messages[message.GetUniqueKey()] = message
			if len(messages) >= mp.batchSize {
				succeed, updatedCount := mp.batchProcessFunc(ctx, messages)
				if succeed {
					messages = map[string]model.Message{}
					if updatedCount > 0 {
						mp.metrics.messagesUpdated.Add(float64(updatedCount))
					}
				}
			}
		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}
			succeed, updatedCount := mp.batchProcessFunc(ctx, messages)
			if succeed {
				messages = map[string]model.Message{}
				if updatedCount > 0 {
					mp.metrics.messagesUpdated.Add(float64(updatedCount))
				}
			}
		}
	}
}

func (mp *MessageProcessor) Stop(ctx context.Context) error {
	mp.stopped.Store(true)
	close(mp.queue)
	<-mp.stop
	return nil
}

func (mp *MessageProcessor) isStopped() bool {
	return mp.stopped.Load()
}
