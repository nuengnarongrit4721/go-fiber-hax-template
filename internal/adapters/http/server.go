package http

import (
	"context"
	"log/slog"

	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	"gofiber-hax/internal/adapters/http/routes"
	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	recovermw "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/timeout"
)

type Server struct {
	app  *fiber.App
	addr string
	host string
}

func NewServer(cfg config.HTTPConfig, set handlers.VersionedSet, opts routes.Options, logger *slog.Logger) *Server {
	app := fiber.New(fiber.Config{
		BodyLimit:    cfg.BodyLimit,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
		ErrorHandler: customErrorHandler(logger),
	})

	app.Use(requestid.New())

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: cfg.AllowOrigins,
			AllowHeaders: cfg.AllowHeaders,
			AllowMethods: cfg.AllowMethods,
		},
	))

	app.Use(helmet.New())
	app.Use(recovermw.New())
	app.Use(timeout.NewWithContext(func(c *fiber.Ctx) error {
		return c.Next()
	}, cfg.RequestTTL))

	if cfg.RateLimit.Enabled {
		app.Use(limiter.New(limiter.Config{
			Max:        cfg.RateLimit.Max,
			Expiration: cfg.RateLimit.Window,
			Next: func(c *fiber.Ctx) bool {
				if c.Method() == fiber.MethodOptions {
					return true
				}
				return middleware.SkipBackgroundTasks(c)
			},
			LimitReached: func(c *fiber.Ctx) error {
				return response.Error(c, fiber.StatusTooManyRequests, "too many requests")
			},
		}))
	}

	if cfg.AccessLog.Enabled {
		app.Use(middleware.AccessLogger(cfg.AccessLog))
	}

	routes.Register(app, set, opts)

	return &Server{
		app:  app,
		addr: cfg.Addr,
		host: cfg.Host,
	}
}

func (s *Server) Start() error {
	return s.app.Listen(s.host + ":" + s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func customErrorHandler(logger *slog.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		msg := "internal server error"
		if fiberErr, ok := err.(*fiber.Error); ok {
			code = fiberErr.Code
			msg = fiberErr.Message
		}
		if code >= fiber.StatusInternalServerError && logger != nil {
			logger.Error("http error",
				"method", c.Method(),
				"path", c.Path(),
				"status", code,
				"error", err.Error(),
			)
		}
		return response.Error(c, code, msg)
	}
}
