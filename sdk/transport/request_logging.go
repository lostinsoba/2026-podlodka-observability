package transport

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"external/sdk/middleware"
)

type RequestLoggingRoundTripper struct {
	internal http.RoundTripper
	logger   *slog.Logger
}

func NewRequestLoggingRoundTripper(transport http.RoundTripper, logger *slog.Logger) http.RoundTripper {
	return &RequestLoggingRoundTripper{
		internal: transport,
		logger:   logger,
	}
}

func (rt *RequestLoggingRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	var (
		requestID = middleware.GetRequestID(req.Context())
	)
	switch req.Method {
	case http.MethodPost, http.MethodPut:
		r1, r2, _ := drainBody(req.Body)
		payload, _ := readBody(r1)
		req.Body = r2
		rt.logger.Debug("logging request",
			slog.String("request_id", requestID),
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
			slog.String("body", payload),
		)
	default:
		rt.logger.Debug("logging request",
			slog.String("request_id", requestID),
			slog.String("method", req.Method),
			slog.String("url", req.URL.String()),
		)
	}
	return rt.internal.RoundTrip(req)
}

func readBody(b io.ReadCloser) (string, error) {
	var buff bytes.Buffer
	_, err := buff.ReadFrom(b)
	if err != nil {
		return "", err
	}
	s := buff.String()
	return s, nil
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err = b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
