package querier

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"external/sdk/client/store/api"
	"external/sdk/transport"

	"querier/internal/tenant"
)

type Querier struct {
	queryInterval  time.Duration
	storeClient    *api.Client
	tenantRegistry *tenant.Registry
	metrics        *metrics
	logger         *slog.Logger
}

type Config struct {
	QueryEndpoint string
	QueryInterval time.Duration
}

func New(cfg Config, tr *tenant.Registry, mr prometheus.Registerer, logger *slog.Logger) *Querier {
	transportWrapper := transport.Wrapper(func(rt http.RoundTripper) http.RoundTripper {
		return transport.NewRequestMetricsRoundTripper(mr, transport.NewRequestLoggingRoundTripper(rt, logger), "querier")
	})
	storeClient := api.New(
		cfg.QueryEndpoint,
		api.OptionWithTransportWrapper(transportWrapper),
	)
	return &Querier{
		queryInterval:  cfg.QueryInterval,
		storeClient:    storeClient,
		tenantRegistry: tr,
		metrics:        newMetrics(mr),
		logger:         logger,
	}
}

func (q *Querier) Run(ctx context.Context) error {
	q.logger.Debug("starting querier",
		slog.String("interval", q.queryInterval.String()),
	)
	ticker := time.NewTicker(q.queryInterval)
	defer ticker.Stop()

	for {
		q.query(ctx)
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return nil
		}
	}
}

func (q *Querier) query(ctx context.Context) {
	tenants := q.tenantRegistry.ListTenants()
	q.metrics.tenants.Set(float64(len(tenants)))

	for _, tenantID := range tenants {
		err := q.queryTenantMessages(ctx, tenantID)
		if err != nil {
			q.logger.Error("failed to query tenant messages",
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
		}
	}
}

func (q *Querier) queryTenantMessages(ctx context.Context, tenantID string) error {
	err := q.queryTenantMessagesByCursor(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to query messages by cursor: %w", err)
	}
	err = q.queryTenantMessagesByOffset(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("failed to query messages by cursor: %w", err)
	}
	return nil
}

const (
	queryLimit = 10
)

func (q *Querier) queryTenantMessagesByCursor(ctx context.Context, tenantID string) error {
	var (
		nextCursor string
	)
	for {
		resp, err := q.storeClient.QueryMessagesByCursor(
			ctx,
			tenantID,
			nextCursor,
			queryLimit,
		)
		if err != nil {
			return err
		}
		for _, message := range resp.Messages {
			q.logger.Info("queried message",
				slog.String("query_kind", "cursor"),
				slog.String("tenant_id", message.TenantID),
				slog.String("message_id", message.ID),
				slog.Int64("timestamp", message.Timestamp),
			)
		}
		if resp.Continue == "" {
			break
		}
		nextCursor = resp.Continue
	}
	return nil
}

func (q *Querier) queryTenantMessagesByOffset(ctx context.Context, tenantID string) error {
	page := 1
	for {
		resp, err := q.storeClient.QueryMessagesByOffset(
			ctx,
			tenantID,
			page,
			queryLimit,
		)
		if err != nil {
			return err
		}
		for _, message := range resp.Messages {
			q.logger.Info("queried message",
				slog.String("query_kind", "offset"),
				slog.String("tenant_id", message.TenantID),
				slog.String("message_id", message.ID),
				slog.Int64("timestamp", message.Timestamp),
			)
		}
		if page >= resp.TotalPages {
			break
		}
		page++
	}
	return nil
}
