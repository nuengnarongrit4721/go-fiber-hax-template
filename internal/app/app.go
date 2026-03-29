package app

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log/slog"
	"os"
	"strings"

	repoMongo "gofiber-hax/internal/adapters/db/mongo/repository"
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
	signer, err := buildSigner(cfg)
	if err != nil {
		return nil, err
	}

	db, closeDB, err := buildDB(cfg)
	if err != nil {
		return nil, err
	}

	repos, err := buildRepos(db)
	if err != nil {
		return nil, err
	}

	services := buildServices(cfg, repos, signer)
	handlers := buildHandlers(services, logger, signer)

	server := httpadapter.NewServer(cfg.HTTP, handlers.HTTP, buildRouteOptions(cfg, signer), logger)

	return &App{
		HTTP:  server,
		Close: closeDB,
	}, nil
}

/*
NOTE: Build Middleware
*/
func buildRouteOptions(cfg config.Config, signer *jwt.Signer) routes.Options {
	opts := routes.Options{Versions: defaultAPIVersions}
	if cfg.Auth.Enabled {
		opts.Protected = []fiber.Handler{middleware.Auth(cfg.Auth, middleware.WithJWTValidator(buildJWTValidator(cfg, signer)))}
	}
	return opts
}

/*
NOTE: Build Repo
*/
func buildRepos(db *DB) (Repos, error) {
	switch db.Driver {
	case "mongo":
		if db.Mongo == nil {
			return Repos{}, fmt.Errorf("mongo driver selected but mongo is not connected")
		}
		return Repos{
			User: repoMongo.NewUserRepo(db.Mongo.DB),
		}, nil

	case "mysql":
		if db.MySQL == nil {
			return Repos{}, fmt.Errorf("mysql driver selected but mysql is not connected")
		}
		return Repos{
			User: repoMsql.NewUserRepo(db.MySQL.DB),
		}, nil

	case "both", "auto":
		if db.MySQL != nil {
			return Repos{
				User: repoMsql.NewUserRepo(db.MySQL.DB),
			}, nil
		}
		if db.Mongo != nil {
			return Repos{
				User: repoMongo.NewUserRepo(db.Mongo.DB),
			}, nil
		}
	}

	return Repos{}, fmt.Errorf("no database available for user repo")
}

/*
NOTE: Build Service
*/
func buildServices(cfg config.Config, repos Repos, signer *jwt.Signer) Services {
	userService := service.NewUserService(repos.User)
	authService := service.NewAuthService(userService, signer, cfg.Auth)
	return Services{
		User: userService,
		Auth: authService,
	}
}

/*
NOTE: Build Handler
*/
func buildHandlers(services Services, logger *slog.Logger, signer *jwt.Signer) HandlerSet {
	var jwksHandler *handlers.JWKSHandler
	if signer != nil {
		jwksHandler = handlers.NewJwksHandler(signer)
	}
	return HandlerSet{
		HTTP: handlers.VersionedSet{
			V1: handlers.Set{
				User: handlers.NewUserHandler(services.User, logger),
				Auth: handlers.NewAuthHandler(services.Auth),
				JWKS: jwksHandler,
			},
			V2: handlers.Set{
				User: nil,
			},
		},
	}
}

/*
NOTE: Build Singner สำหรับ JWT
*/
func buildSigner(cfg config.Config) (*jwt.Signer, error) {
	if !usesInternalJWT(cfg.Auth.Mode) {
		return nil, nil
	}

	keyPath := strings.TrimSpace(cfg.Auth.JWT.PrivateKeyPath)
	if keyPath == "" {
		return nil, fmt.Errorf("JWT_PRIVATE_KEY_PATH is required for internal JWT mode")
	}

	var (
		priv *rsa.PrivateKey
		pub  *rsa.PublicKey
		err  error
	)

	if isProduction(cfg.App.Env) {
		if _, statErr := os.Stat(keyPath); statErr != nil {
			if os.IsNotExist(statErr) {
				return nil, fmt.Errorf("jwt private key not found at %s; create it before starting production", keyPath)
			}
			return nil, statErr
		}
		priv, pub, err = jwt.LoadKeys(keyPath)
	} else {
		priv, pub, err = jwt.LoadOrGenerateKeys(keyPath)
	}
	if err != nil {
		return nil, err
	}

	return jwt.NewSigner(priv, pub, cfg.Auth.JWT.KeyID), nil
}

/*
NOTE: Build Validator สำหรับ JWT
*/
func buildJWTValidator(cfg config.Config, signer *jwt.Signer) *jwt.Validator {
	mode := normalizeAuthMode(cfg.Auth.Mode)
	switch mode {
	case "jwt":
		if signer == nil {
			return nil
		}
		return jwt.NewValidator(cfg.Auth.JWT, jwt.WithStaticPublicKey(signer.PublicKey()))
	case "jwks", "google":
		return jwt.NewValidator(cfg.Auth.JWT)
	default:
		return nil
	}
}
