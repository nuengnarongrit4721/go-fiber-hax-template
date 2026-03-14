package logs

import (
	"log/slog"
	"os"
	"strings"

	"gofiber-hax/internal/infra/config"
)

func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	if strings.ToLower(cfg.Format) == "pretty" {
		logger := slog.New(NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
		slog.SetDefault(logger)
		return logger
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)
	return logger
}

func parseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
