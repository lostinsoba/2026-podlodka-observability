package sender

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/prometheus/client_golang/prometheus"

	"external/sdk/client/store/receiver"
	"external/sdk/transport"

	"sender/internal/model"
	"sender/internal/tenant"
)

type Sender struct {
	sendInterval    time.Duration
	sendConcurrency int
	storeClient     *receiver.Client
	tenantRegistry  *tenant.Registry
	metrics         *metrics
	logger          *slog.Logger
}

type Config struct {
	SendEndpoint    string
	SendInterval    time.Duration
	SendConcurrency int
}

func New(cfg Config, tr *tenant.Registry, mr prometheus.Registerer, logger *slog.Logger) *Sender {
	transportWrapper := transport.Wrapper(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewRequestMetricsRoundTripper(mr, transport.NewRequestLoggingRoundTripper(rt, logger), "sender")
	})
	storeClient := receiver.New(
		cfg.SendEndpoint,
		receiver.OptionWithTransportWrapper(transportWrapper),
	)
	return &Sender{
		sendInterval:    cfg.SendInterval,
		sendConcurrency: cfg.SendConcurrency,
		storeClient:     storeClient,
		tenantRegistry:  tr,
		metrics:         newMetrics(mr),
		logger:          logger,
	}
}

func (s *Sender) Run(ctx context.Context) error {
	s.logger.Debug("starting sender",
		slog.String("interval", s.sendInterval.String()),
	)
	ticker := time.NewTicker(s.sendInterval)
	defer ticker.Stop()

	for {
		s.send(ctx)
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
}

func (s *Sender) send(ctx context.Context) {
	tenantCfgs := s.tenantRegistry.ListTenantConfigs()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(s.sendConcurrency)

	for _, tenantCfg := range tenantCfgs {
		cfgCopy := tenantCfg
		g.Go(func() error {
			return s.sendTenantMessages(ctx, cfgCopy)
		})
	}

	err := g.Wait()
	if err != nil {
		s.logger.Error("failed to send messages",
			slog.Any("error", err),
		)
	}
}

func (s *Sender) sendTenantMessages(ctx context.Context, cfg model.TenantConfig) error {
	messages := generateMessages(cfg.TenantID, cfg.BatchSize)
	err := s.storeClient.SendMessages(ctx, cfg.TenantID, messages)
	if err != nil {
		s.logger.Error("failed to send tenant messages",
			slog.String("tenant_id", cfg.TenantID),
			slog.Any("error", err),
		)
		s.metrics.messagesFailed.WithLabelValues(cfg.TenantID).Add(float64(cfg.BatchSize))
		return err
	}
	s.metrics.messagesSent.WithLabelValues(cfg.TenantID).Add(float64(cfg.BatchSize))
	return nil
}

func generateMessages(tenantID string, batchSize int) receiver.Messages {
	messages := make([]receiver.Message, batchSize)
	for i := 0; i < batchSize; i++ {
		messages[i] = receiver.Message{
			ID:        generateID(tenantID),
			TenantID:  tenantID,
			Message:   "test",
			Timestamp: time.Now().UnixMilli(),
		}
	}
	return receiver.Messages{Messages: messages}
}

var seq atomic.Int64

func generateID(tenantID string) string {
	return fmt.Sprintf("%s-%d-%d",
		tenantID,
		time.Now().UnixMilli(),
		seq.Add(1),
	)
}
