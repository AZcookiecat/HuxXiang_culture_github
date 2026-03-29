package app

import (
	"log/slog"
	"os"
)

func NewLogger(cfg Config) *slog.Logger {
	var handler slog.Handler
	if cfg.LogJSON {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	}
	return slog.New(handler).With(
		slog.String("service", cfg.AppName),
		slog.String("env", cfg.Environment),
	)
}
