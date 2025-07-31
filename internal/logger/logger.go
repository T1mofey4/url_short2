package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	level := slog.LevelInfo // as default

	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: level}
	handler := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(handler)
}
