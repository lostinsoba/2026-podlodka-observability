package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"external/sdk/interceptor"
)

const (
	tenantIDHeaderName = "X-Tenant-ID"
)

type Messages struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	ID        string `json:"id"`
	TenantID  string `json:"tenant_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type MessagesByPage struct {
	Messages   []Message `json:"messages"`
	TotalPages int       `json:"total_pages"`
	TotalCount int       `json:"total_count"`
}

func (c *Client) QueryMessagesByOffset(ctx context.Context, tenantID string, page int, pageSize int) (*MessagesByPage, error) {
	queryParams := url.Values{}
	queryParams.Add("page", strconv.Itoa(page))
	queryParams.Add("page_size", strconv.Itoa(pageSize))
	data, err := c.GET(
		ctx,
		"message/queryByOffset",
		http.StatusOK,
		interceptor.WithRequestHeader(tenantIDHeaderName, tenantID),
		interceptor.WithRequestQueryParams(queryParams),
	)
	if err != nil {
		return nil, err
	}
	var messagesByPage MessagesByPage
	err = json.Unmarshal(data, &messagesByPage)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &messagesByPage, nil
}

type MessagesByCursor struct {
	Messages []Message `json:"messages"`
	Continue string    `json:"continue"`
}

func (c *Client) QueryMessagesByCursor(ctx context.Context, tenantID string, nextContinue string, limit int) (*MessagesByCursor, error) {
	queryParams := url.Values{}
	queryParams.Add("continue", nextContinue)
	queryParams.Add("limit", strconv.Itoa(limit))
	data, err := c.GET(
		ctx,
		"message/queryByCursor",
		http.StatusOK,
		interceptor.WithRequestHeader(tenantIDHeaderName, tenantID),
		interceptor.WithRequestQueryParams(queryParams),
	)
	if err != nil {
		return nil, err
	}
	var messagesByCursor MessagesByCursor
	err = json.Unmarshal(data, &messagesByCursor)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return &messagesByCursor, nil
}
