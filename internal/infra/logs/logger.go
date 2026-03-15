package logs

import (
	"log/slog"
	"os"

	"gofiber-hax/internal/infra/config"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	zapLevel := zap.NewAtomicLevelAt(slogLevelToZap(level))
	if cfg.Level == "info" {
		cfg.Format = "json"
	}
	logFormat = normalizeFormat(cfg.Format)
	prodEncCfg := zap.NewProductionEncoderConfig()
	prodEncCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := selectEncoder(logFormat, prodEncCfg)
	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
	zapLogger := zap.New(core)

	logger := slog.New(slogzap.Option{
		Level:     level,
		Logger:    zapLogger,
		AddSource: true,
	}.NewZapHandler())

	slog.SetDefault(logger)
	return logger
}
