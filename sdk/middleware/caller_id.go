package middleware

import (
	"context"
	"net/http"
)

type callerIDCtxKey struct{}

func WithCallerID(next http.Handler, callerID string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.WithValue(request.Context(), callerIDCtxKey{}, callerID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetCallerID(ctx context.Context) string {
	callerID, ok := ctx.Value(callerIDCtxKey{}).(string)
	if ok {
		return callerID
	}
	return ""
}
