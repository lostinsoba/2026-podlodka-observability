package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"external/sdk/middleware"
	"github.com/prometheus/client_golang/prometheus"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"store/internal/controller"
	"store/internal/server/api/dto"
)

type Router struct {
	ctrl   *controller.Controller
	mr     prometheus.Registerer
	tp     *sdktrace.TracerProvider
	logger *slog.Logger
}

func New(ctrl *controller.Controller, mr prometheus.Registerer, tp *sdktrace.TracerProvider, logger *slog.Logger) *Router {
	return &Router{
		ctrl:   ctrl,
		mr:     mr,
		tp:     tp,
		logger: logger,
	}
}

func (r *Router) Route() http.Handler {
	var queryMessagesByCursor http.Handler = r.queryMessagesByCursor()
	queryMessagesByCursor = middleware.WithCallerID(queryMessagesByCursor, "queryMessagesByCursor")
	queryMessagesByCursor = middleware.WithHandlerMetrics(r.mr, queryMessagesByCursor, "queryMessagesByCursor")

	var queryMessagesByOffset http.Handler = r.queryMessagesByOffset()
	queryMessagesByOffset = middleware.WithCallerID(queryMessagesByOffset, "queryMessagesByOffset")
	queryMessagesByOffset = middleware.WithHandlerMetrics(r.mr, queryMessagesByOffset, "queryMessagesByOffset")

	mux := http.NewServeMux()
	mux.Handle("/message/queryByCursor", queryMessagesByCursor)
	mux.Handle("/message/queryByOffset", queryMessagesByOffset)

	tr := r.tp.Tracer("api")
	route := middleware.WithTenantID(mux)
	route = middleware.WithRequestID(route)
	route = middleware.WithTraceID(tr, route)
	return route
}

func renderResponse(writer http.ResponseWriter, response any) {
	data, _ := json.Marshal(response)
	_, _ = writer.Write(data)
}

func renderError(writer http.ResponseWriter, err *dto.ErrorData) {
	writer.WriteHeader(err.StatusCode)
	data, _ := json.Marshal(err)
	_, _ = writer.Write(data)
}
