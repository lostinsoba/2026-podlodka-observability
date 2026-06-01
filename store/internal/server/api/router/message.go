package router

import (
	"log/slog"
	"net/http"

	"external/sdk/middleware"
	"store/internal/server/api/dto"
	"store/internal/server/api/query"
)

func (r *Router) queryMessagesByCursor() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			requestID = middleware.GetRequestID(request.Context())
			callerID  = middleware.GetCallerID(request.Context())
			tenantID  = middleware.GetTenantID(request.Context())
		)
		if tenantID == "" {
			renderError(writer, dto.InvalidRequestError(requestID, "missing tenant id"))
			return
		}
		q, err := query.ParseMessagesByCursorQuery(request)
		if err != nil {
			r.logger.Error("failed to parse messages by cursor query",
				slog.String("request_id", requestID),
				slog.String("called id", callerID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			renderError(writer, dto.InvalidRequestError(requestID, "failed to parse query messages by cursor request"))
			return
		}
		messagesByCursorQueryRequest := q.ToModel(tenantID)
		messagesByCursorQueryResponse, err := r.ctrl.QueryMessagesByCursor(request.Context(), messagesByCursorQueryRequest)
		if err != nil {
			r.logger.Error("failed to query messages by cursor",
				slog.String("request_id", requestID),
				slog.String("called id", callerID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			renderError(writer, dto.InternalServerError(requestID, "failed to query messages by cursor"))
			return
		}
		renderResponse(writer, dto.ToMessagesByCursorQueryResponseData(messagesByCursorQueryResponse))
	}
}

func (r *Router) queryMessagesByOffset() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var (
			requestID = middleware.GetRequestID(request.Context())
			callerID  = middleware.GetCallerID(request.Context())
			tenantID  = middleware.GetTenantID(request.Context())
		)
		if tenantID == "" {
			renderError(writer, dto.InvalidRequestError(requestID, "missing tenant id"))
			return
		}
		q, err := query.ParseMessagesByOffsetQuery(request)
		if err != nil {
			r.logger.Error("failed to parse messages by offset query",
				slog.String("request_id", requestID),
				slog.String("called id", callerID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			renderError(writer, dto.InvalidRequestError(requestID, "failed to parse messages by offset query request"))
			return
		}
		messagesByOffsetQueryRequest := q.ToModel(tenantID)
		messagesByOffsetQueryResponse, err := r.ctrl.QueryMessagesByOffset(request.Context(), messagesByOffsetQueryRequest)
		if err != nil {
			r.logger.Error("failed query messages by offset",
				slog.String("request_id", requestID),
				slog.String("called id", callerID),
				slog.String("tenant_id", tenantID),
				slog.Any("error", err),
			)
			renderError(writer, dto.InternalServerError(requestID, "failed to query messages by offset"))
			return
		}
		renderResponse(writer, dto.ToMessagesByOffsetQueryResponseData(messagesByOffsetQueryResponse))
	}
}
