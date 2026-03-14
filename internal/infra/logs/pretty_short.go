package logs

import "log/slog"

func DebugPre(msg string, v any) {
	slog.Default().Debug(msg, "data", Pretty(v))
}

func InfoPre(msg string, v any) {
	slog.Default().Info(msg, "data", Pretty(v))
}

func WarnPre(msg string, v any) {
	slog.Default().Warn(msg, "data", Pretty(v))
}

func ErrorPre(msg string, v any) {
	slog.Default().Error(msg, "data", Pretty(v))
}
