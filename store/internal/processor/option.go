package processor

import (
	"context"
	"time"

	"store/internal/model"
)

const (
	defaultBatchSize            = 100
	defaultQueueCapacity        = 10000
	defaultBatchProcessInterval = time.Minute
)

type BatchProcessFunc func(ctx context.Context, messages map[string]model.Message) (succeed bool, updatedCount int64)

type MessageStatsFunc func(ctx context.Context) (stats model.MessagesStats, err error)

type Option func(*MessageProcessor)

func OptionBatchSize(batchSize int) func(*MessageProcessor) {
	return func(mp *MessageProcessor) {
		mp.batchSize = batchSize
	}
}

func OptionQueueCapacity(queueCapacity int) func(*MessageProcessor) {
	return func(mp *MessageProcessor) {
		mp.queue = make(chan model.Message, queueCapacity)
	}
}

func OptionBatchProcessInterval(batchProcessInterval time.Duration) func(*MessageProcessor) {
	return func(mp *MessageProcessor) {
		mp.batchProcessInterval = batchProcessInterval
	}
}
