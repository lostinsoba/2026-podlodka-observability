package database

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"store/internal/config"
	"store/internal/database/postgres"
	"store/internal/model"
)

type Database interface {
	Message
	Stats
}

type Message interface {
	UpdateMessages(ctx context.Context, request model.UpdateMessagesRequest) (response model.UpdateMessagesResponse, err error)
	QueryMessagesByCursor(ctx context.Context, request model.MessagesByCursorQueryRequest) (model.MessagesByCursorQueryResponse, error)
	QueryMessagesByOffset(ctx context.Context, request model.MessageByOffsetQueryRequest) (model.MessageByOffsetQueryResponse, error)
}

type Stats interface {
	EvaluateMessagesStats(ctx context.Context) (model.MessagesStats, error)
}

func New(cfg config.Database, mr prometheus.Registerer) (Database, error) {
	return postgres.New(cfg, mr)
}
