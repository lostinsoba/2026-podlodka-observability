package receiver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

func (c *Client) SendMessages(ctx context.Context, tenantID string, messages Messages) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("failed to marshal tenant messages: %w", err)
	}
	_, err = c.POST(
		ctx,
		"message",
		bytes.NewReader(data),
		http.StatusCreated,
		interceptor.WithRequestHeader(tenantIDHeaderName, tenantID),
	)
	return err
}
