package logs

import (
	"log/slog"
	"os"
	"strings"

	"gofiber-hax/internal/infra/config"

	slogzap "github.com/samber/slog-zap/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	zapLevel := zap.NewAtomicLevelAt(slogLevelToZap(level))

	var encoder zapcore.Encoder

	switch strings.ToLower(cfg.Format) {
	case "pretty":
		encCfg := zap.NewDevelopmentEncoderConfig()
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encCfg)
	case "text", "logfmt", "line":
		encCfg := zap.NewProductionEncoderConfig()
		encoder = zapcore.NewConsoleEncoder(encCfg)
	default:
		encCfg := zap.NewProductionEncoderConfig()
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapLevel)
	zapLogger := zap.New(core)

	logger := slog.New(slogzap.Option{Level: level, Logger: zapLogger}.NewZapHandler())
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

func slogLevelToZap(l slog.Level) zapcore.Level {
	switch l {
	case slog.LevelDebug:
		return zapcore.DebugLevel
	case slog.LevelInfo:
		return zapcore.InfoLevel
	case slog.LevelWarn:
		return zapcore.WarnLevel
	case slog.LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
