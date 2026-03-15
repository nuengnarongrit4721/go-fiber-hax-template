package logs

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func Error(message interface{}, args ...any) {
	pc := callerPC(3)
	details := parseError(message)

	if logFormat == formatPretty {
		printPrettyError(os.Stdout, pc, details, args...)
		return
	}

	record := slog.NewRecord(time.Now(), slog.LevelError, details.Message, pc)
	if details.Trace != "" {
		record.Add(slog.String("trace", details.Trace))
	}
	record.Add(args...)
	_ = slog.Default().Handler().Handle(context.Background(), record)
}

func Debug(message interface{}, args ...any) {
	logWithLevel(slog.LevelDebug, message, args...)
}

func Info(message interface{}, args ...any) {
	logWithLevel(slog.LevelInfo, message, args...)
}

func Warn(message interface{}, args ...any) {
	logWithLevel(slog.LevelWarn, message, args...)
}

func logWithLevel(level slog.Level, message interface{}, args ...any) {
	ctx := context.Background()
	if !slog.Default().Enabled(ctx, level) {
		return
	}
	pc := callerPC(4)
	details := parseError(message)
	if logFormat == formatPretty && level == slog.LevelDebug {
		printPrettyDebug(debugOutput(), pc, details, message, args...)
		return
	}
	record := slog.NewRecord(time.Now(), level, details.Message, pc)
	if details.Trace != "" {
		record.Add(slog.String("trace", details.Trace))
	}
	if !isSimpleMessage(message) {
		record.Add(slog.Any("data", message))
	}
	record.Add(args...)
	_ = slog.Default().Handler().Handle(ctx, record)
}
