package http

import (
	"context"
	"log/slog"

	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	"gofiber-hax/internal/adapters/http/routes"
	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app  *fiber.App
	addr string
}

func NewServer(cfg config.HTTPConfig, set handlers.VersionedSet, opts routes.Options, logger *slog.Logger) *Server {
	app := fiber.New()

	app.Use(middleware.AccessLogger(cfg.AccessLog))

	routes.Register(app, set, opts)

	return &Server{app: app, addr: cfg.Addr}
}

func (s *Server) Start() error {
	return s.app.Listen(s.addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}
