package middleware

import (
	"os"

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

	return logger.New(logger.Config{
		Format:     format,
		TimeFormat: cfg.TimeFormat,
		Output:     os.Stdout,
	})
}
