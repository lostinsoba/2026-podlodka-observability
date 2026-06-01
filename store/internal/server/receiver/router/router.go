package router

import (
	"log/slog"
	"net/http"

	"external/sdk/middleware"
	"github.com/prometheus/client_golang/prometheus"

	"store/internal/controller"
)

type Router struct {
	ctrl   *controller.Controller
	mr     prometheus.Registerer
	logger *slog.Logger
}

func New(ctrl *controller.Controller, mr prometheus.Registerer, logger *slog.Logger) *Router {
	return &Router{
		ctrl:   ctrl,
		mr:     mr,
		logger: logger,
	}
}

func (r *Router) Route() http.Handler {
	var message http.Handler = r.processMessages()
	message = middleware.WithCallerID(message, "message")
	message = middleware.WithHandlerMetrics(r.mr, message, "message")

	mux := http.NewServeMux()
	mux.Handle("/message", message)

	route := middleware.WithTenantID(mux)
	route = middleware.WithRequestID(route)
	return route
}
