package logs

import (
	"log/slog"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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

func normalizeFormat(format string) string {
	return strings.ToLower(strings.TrimSpace(format))
}

func selectEncoder(format string, prodEncCfg zapcore.EncoderConfig) zapcore.Encoder {
	switch format {
	case formatPretty:
		encCfg := zap.NewDevelopmentEncoderConfig()
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		return zapcore.NewConsoleEncoder(encCfg)
	case formatText, formatLogfmt, formatLine:
		return zapcore.NewConsoleEncoder(prodEncCfg)
	default:
		return zapcore.NewJSONEncoder(prodEncCfg)
	}
}
