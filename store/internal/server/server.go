package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

type Server struct {
	httpServer *http.Server
	logger     *slog.Logger
}

func New(router http.Handler, port int, logger *slog.Logger) (*Server, error) {
	return &Server{
		httpServer: &http.Server{
			Addr:        fmt.Sprintf(":%d", port),
			Handler:     router,
			ReadTimeout: defaultReadTimeout,
		},
		logger: logger,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		s.logger.Info("stopping server")
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, defaultShutdownTimeout)
		defer cancel()
		err := s.httpServer.Shutdown(shutdownCtx)
		if err != nil {
			s.logger.Error("failed to gracefully stop server",
				slog.Any("error", err),
			)
		}
	}()
	s.logger.Info("starting server",
		slog.String("addr", s.httpServer.Addr),
	)
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("server stopped")
	return nil
}
