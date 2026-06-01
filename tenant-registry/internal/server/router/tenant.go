package router

import (
	"log/slog"
	"net/http"

	"external/sdk/middleware"

	"tenant-registry/internal/server/dto"
)

func (r *Router) listTenantConfigurations() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			requestID = middleware.GetRequestID(request.Context())
		)
		r.logger.Debug("list tenant configurations",
			slog.String("request_id", requestID),
		)
		list := r.ctrl.ListTenantConfigurations(request.Context())
		res := dto.ToTenantConfigurationsData(list)
		renderResponse(writer, res)
	}
}
