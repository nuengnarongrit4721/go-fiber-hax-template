package middleware

import (
	"os"

	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func AccessLogger(cfg config.AccessLogConfig) fiber.Handler {
	return logger.New(logger.Config{
		Format:     cfg.Format,
		TimeFormat: cfg.TimeFormat,
		Output:     os.Stdout,
	})
}
