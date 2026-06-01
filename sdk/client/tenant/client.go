package tenant

import (
	"external/sdk/transport"
	"net/http"
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
