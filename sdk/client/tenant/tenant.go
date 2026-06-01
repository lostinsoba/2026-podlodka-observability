package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Configuration struct {
	TenantID           string `json:"tenant_id"`
	TenantMaxBatchSize int    `json:"tenant_max_batch_size"`
}

type Configurations struct {
	Tenants []Configuration `json:"tenants"`
}

func (c *Client) GetTenantConfigurations(ctx context.Context) ([]Configuration, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}
	var resp *http.Response
	resp, err = c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to process request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid request status code: %d", resp.StatusCode)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var res Configurations
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	return res.Tenants, nil
}
