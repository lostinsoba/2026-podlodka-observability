package router

import (
	"log/slog"
	"net/http"

	"external/sdk/middleware"
	"store/internal/server/receiver/dto"
)

func (r *Router) processMessages() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			requestID = middleware.GetRequestID(request.Context())
			tenantID  = middleware.GetTenantID(request.Context())
		)
		messageData, err := dto.ParseMessagesData(writer, request)
		if err != nil {
			r.logger.Error("failed to parse message data",
				slog.String("request_id", requestID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			http.Error(writer, requestID, http.StatusBadRequest)
			return
		}
		r.logger.Debug("received message events",
			slog.String("request_id", requestID),
			slog.String("tenant_id", tenantID),
			slog.Int("count", len(messageData.Messages)),
		)
		err = r.ctrl.ScheduleMessagesSave(request.Context(), messageData.ToModel())
		if err != nil {
			r.logger.Error("failed to schedule messages",
				slog.String("request_id", requestID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			http.Error(writer, requestID, http.StatusInternalServerError)
		}
		writer.WriteHeader(http.StatusCreated)
	}
}
