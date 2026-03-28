package app

import (
	"context"
	"fmt"
	"log/slog"

	repoMogo "gofiber-hax/internal/adapters/db/mongo/repository"
	repoMsql "gofiber-hax/internal/adapters/db/mysql/repository"
	httpadapter "gofiber-hax/internal/adapters/http"
	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	"gofiber-hax/internal/adapters/http/routes"
	"gofiber-hax/internal/core/service"
	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/jwt"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	HTTP  *httpadapter.Server
	Close func(ctx context.Context) error
}

var defaultAPIVersions = []string{"v1", "v2"}

/*
NOTE: Build App
*/
func Build(cfg config.Config, logger *slog.Logger) (*App, error) {

	// รันตัวปั๊มกุญแจ หาไฟล์ไม่เจอก็สร้างให้เลย
	priv, pub, err := jwt.LoadOrGenerateKeys("keys/jwt_private.pem")
	if err != nil {
		return nil, err
	}
	signer := jwt.NewSigner(priv, pub)

	db, closeDB, err := buildDB(cfg)
	if err != nil {
		return nil, err
	}

	repos, err := buildRepos(db)
	if err != nil {
		return nil, err
	}

	services := buildServices(repos)
	handlers := buildHandlers(services, logger, signer)

	server := httpadapter.NewServer(cfg.HTTP, handlers.HTTP, buildRouteOptions(cfg), logger)

	return &App{
		HTTP:  server,
		Close: closeDB,
	}, nil
}

/*
NOTE: Build Middleware
*/
func buildRouteOptions(cfg config.Config) routes.Options {
	opts := routes.Options{Versions: defaultAPIVersions}
	if cfg.Auth.Enabled {
		opts.Protected = []fiber.Handler{middleware.Auth(cfg.Auth)}
	}
	return opts
}

/*
NOTE: Build Repo
*/
func buildRepos(db *DB) (Repos, error) {
	if db.MySQL != nil {
		return Repos{
			User: repoMsql.NewUserRepo(db.MySQL.DB),
		}, nil
	}
	if db.Mongo != nil {
		return Repos{
			User: repoMogo.NewUserRepo(db.Mongo.DB),
		}, nil
	}
	return Repos{}, fmt.Errorf("no database available for user repo")
}

/*
NOTE: Build Service
*/
func buildServices(repos Repos) Services {
	userService := service.NewUserService(repos.User)
	authService := service.NewAuthService(userService)
	return Services{
		User: userService,
		Auth: authService,
	}
}

/*
NOTE: Build Handler
*/
func buildHandlers(services Services, logger *slog.Logger, signer *jwt.Signer) HandlerSet {
	return HandlerSet{
		HTTP: handlers.VersionedSet{
			V1: handlers.Set{
				User: handlers.NewUserHandler(services.User, logger),
				Auth: handlers.NewAuthHandler(services.Auth),
				JWKS: handlers.NewJwksHandler(signer),
			},
			V2: handlers.Set{
				User: nil,
			},
		},
	}
}
