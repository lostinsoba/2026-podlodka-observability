package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	requestIDHeaderName = "X-Request-ID"
)

type requestIDCtxKey struct{}

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var requestID = request.Header.Get(requestIDHeaderName)
		if requestID == "" {
			requestID = generateRequestID()
		}
		writer.Header().Set(requestIDHeaderName, requestID)
		ctx := context.WithValue(request.Context(), requestIDCtxKey{}, requestID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func generateRequestID() string {
	return fmt.Sprintf("unknown-%d", time.Now().UnixMilli())
}

func GetRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDCtxKey{}).(string)
	if ok {
		return requestID
	}
	return ""
}
