package router

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"external/sdk/middleware"

	"tenant-registry/internal/controller"
)

type Router struct {
	ctrl   *controller.Controller
	logger *slog.Logger
}

func New(ctrl *controller.Controller, logger *slog.Logger) *Router {
	return &Router{
		ctrl:   ctrl,
		logger: logger,
	}
}

func (r *Router) Route() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /", r.listTenantConfigurations())

	route := middleware.WithRequestID(mux)
	return route
}

func renderResponse(writer http.ResponseWriter, response any) {
	data, _ := json.Marshal(response)
	_, _ = writer.Write(data)
}
