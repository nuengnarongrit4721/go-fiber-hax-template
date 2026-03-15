package middleware

import (
	"os"
	"strings"

	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AccessLogger(cfg config.AccessLogConfig) fiber.Handler {
	format := cfg.Format
	if format == "" || format == "pretty" {
		// รูปแบบ CommonFormat แบบที่เห็นในภาพ
		format = "${ip} - - [${time}] \"${method} ${url} ${protocol}\" ${status} ${bytesSent}\n"
	}
	format = ensureTrailingNewline(format)

	return logger.New(logger.Config{
		Format:     format,
		TimeFormat: cfg.TimeFormat,
		Output:     os.Stdout,
	})
}

func ensureTrailingNewline(format string) string {
	if format == "" {
		return format
	}
	if strings.HasSuffix(format, "\n") {
		return format
	}
	return format + "\n"
}
