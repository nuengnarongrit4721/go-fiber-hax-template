package http

import (
	"context"
	"log/slog"

	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	"gofiber-hax/internal/adapters/http/routes"
	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
)

type Server struct {
	app  *fiber.App
	addr string
	host string
}

func NewServer(cfg config.HTTPConfig, set handlers.VersionedSet, opts routes.Options, logger *slog.Logger) *Server {
	app := fiber.New()

	app.Use(cors.New(
		cors.Config{
			AllowOrigins: cfg.AlowOrigins,
			AllowHeaders: cfg.AllowHeaders,
			AllowMethods: cfg.AllowMethods,
		},
	))

	app.Use(helmet.New())

	/* //Rate Limit
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return response.Error(c, fiber.StatusTooManyRequests, "Too many requests. Please try again later.")
		},
	}))
	*/

	app.Use(middleware.AccessLogger(cfg.AccessLog))

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
