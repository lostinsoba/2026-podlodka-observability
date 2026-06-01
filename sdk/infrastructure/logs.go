package infrastructure

import (
	"log/slog"
	"os"
)

func NewLogger(buildInfo BuildInfo, level slog.Level) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler).With(
		slog.String("service", buildInfo.Service),
		slog.String("component", buildInfo.Component),
		slog.String("version", buildInfo.Version),
		slog.String("gitCommit", buildInfo.GitCommit),
	)
}
