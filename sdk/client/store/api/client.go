package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"external/sdk/interceptor"
	"external/sdk/transport"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(endpoint string, opts ...Option) *Client {
	c := &Client{
		baseURL: endpoint,
		httpClient: &http.Client{
			Transport: buildDefaultTransport(),
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type Option func(*Client)

func OptionWithTransportWrapper(transportWrapper transport.Wrapper) Option {
	return func(c *Client) {
		c.httpClient.Transport = transportWrapper(buildDefaultTransport())
	}
}

func buildDefaultTransport() *http.Transport {
	return &http.Transport{DisableKeepAlives: true}
}

func (c *Client) GET(
	ctx context.Context,
	endpoint string,
	expectedStatusCode int,
	ics ...interceptor.Interceptor,
) ([]byte, error) {
	return c.processRequest(ctx, http.MethodGet, endpoint, nil, expectedStatusCode, ics...)
}

func (c *Client) processRequest(
	ctx context.Context,
	method string,
	endpoint string,
	r io.Reader,
	expectedStatusCode int,
	ics ...interceptor.Interceptor,
) ([]byte, error) {
	u, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, method, u, r)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters: %w", err)
	}
	var do = interceptor.Doer(c.httpClient.Do)
	for _, ic := range ics {
		do = ic(do)
	}
	response, doErr := do(request)
	if doErr != nil {
		return nil, doErr
	}
	defer func() {
		_ = response.Body.Close()
	}()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if response.StatusCode == expectedStatusCode {
		return data, nil
	}
	if message := extractErrorMessage(data); message != "" {
		return nil, fmt.Errorf("invalid status: %s, message: %s", response.Status, message)
	}
	return nil, fmt.Errorf("invalid status: %s", response.Status)
}

func extractErrorMessage(data []byte) (msg string) {
	return string(data)
}
