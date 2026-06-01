package controller

import (
	"context"
	"errors"
	"external/sdk/middleware"
	"log/slog"
	"store/internal/model"
)

var ErrTenantNotFound = errors.New("tenant not found")

func (c *Controller) ScheduleMessagesSave(ctx context.Context, messages []model.Message) error {
	f := func(ctx context.Context) error {
		return c.scheduleMessagesSave(ctx, messages)
	}
	m := newCtrlMiddlewareChain(
		c.withDurationMeasure,
	)
	err := m.apply(ctx, "ScheduleMessageSave", f)
	return err
}

func (c *Controller) scheduleMessagesSave(ctx context.Context, messages []model.Message) error {
	var (
		tenantID = middleware.GetTenantID(ctx)
	)
	if !c.tr.Lookup(tenantID) {
		return ErrTenantNotFound
	}
	return c.mp.Queue(ctx, messages...)
}

func (c *Controller) QueryMessagesByCursor(ctx context.Context, request model.MessagesByCursorQueryRequest) (resp model.MessagesByCursorQueryResponse, err error) {
	f := func(ctx context.Context) error {
		resp, err = c.queryMessagesByCursor(ctx, request)
		return err
	}
	m := newCtrlMiddlewareChain(
		c.withDurationMeasure,
		c.withTracing,
	)
	_ = m.apply(ctx, "QueryMessagesByCursor", f)
	return
}

func (c *Controller) queryMessagesByCursor(ctx context.Context, request model.MessagesByCursorQueryRequest) (resp model.MessagesByCursorQueryResponse, err error) {
	var (
		requestID = middleware.GetRequestID(ctx)
	)
	c.logger.Debug("query messages by cursor",
		slog.String("request_id", requestID),
		slog.String("tenant", request.TenantID),
	)
	return c.d.QueryMessagesByCursor(ctx, request)
}

func (c *Controller) QueryMessagesByOffset(ctx context.Context, request model.MessageByOffsetQueryRequest) (resp model.MessageByOffsetQueryResponse, err error) {
	f := func(ctx context.Context) error {
		resp, err = c.queryMessagesByOffset(ctx, request)
		return err
	}
	m := newCtrlMiddlewareChain(
		c.withDurationMeasure,
		c.withTracing,
	)
	_ = m.apply(ctx, "QueryMessagesByOffset", f)
	return
}

func (c *Controller) queryMessagesByOffset(ctx context.Context, request model.MessageByOffsetQueryRequest) (resp model.MessageByOffsetQueryResponse, err error) {
	var (
		requestID = middleware.GetRequestID(ctx)
	)
	c.logger.Debug("query messages by offset",
		slog.String("request_id", requestID),
		slog.String("tenant", request.TenantID),
	)
	return c.d.QueryMessagesByOffset(ctx, request)
}
