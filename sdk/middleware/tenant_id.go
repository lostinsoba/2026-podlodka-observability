package middleware

import (
	"context"
	"net/http"
)

const (
	tenantIDHeaderName = "X-Tenant-ID"
)

type tenantIDCtxKey struct{}

func WithTenantID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var tenantID = request.Header.Get(tenantIDHeaderName)
		ctx := context.WithValue(request.Context(), tenantIDCtxKey{}, tenantID)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetTenantID(ctx context.Context) string {
	tenantID, ok := ctx.Value(tenantIDCtxKey{}).(string)
	if ok {
		return tenantID
	}
	return ""
}
