package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"

	"gofiber-hax/internal/adapters/db/mongo"
	"gofiber-hax/internal/adapters/db/mysql"
	"gofiber-hax/internal/adapters/db/mysql/repository"
	httpadapter "gofiber-hax/internal/adapters/http"
	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	"gofiber-hax/internal/adapters/http/routes"
	"gofiber-hax/internal/core/ports/out"
	"gofiber-hax/internal/core/service"
	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
)

type App struct {
	HTTP  *httpadapter.Server
	Close func(ctx context.Context) error
}

type DB struct {
	Driver       string
	Mongo        *mongo.Connector
	MySQL        *mysql.Connector
	MySQLReplica *mysql.Connector
}

type Repos struct {
	User out.UserRepository
}

type Services struct {
	User *service.UserService
}

type HandlerSet struct {
	HTTP handlers.VersionedSet
}

func Build(cfg config.Config, logger *slog.Logger) (*App, error) {
	db, closeDB, err := buildDB(cfg)
	if err != nil {
		return nil, err
	}

	repos, err := buildRepos(db)
	if err != nil {
		return nil, err
	}

	services := buildServices(repos)
	handlers := buildHandlers(services, logger)

	server := httpadapter.NewServer(cfg.HTTP, handlers.HTTP, routes.Options{
		Versions:  []string{"v1"},
		Public:    nil,
		Protected: []fiber.Handler{middleware.Auth(cfg.Auth)},
	})

	return &App{
		HTTP:  server,
		Close: closeDB,
	}, nil
}

func buildDB(cfg config.Config) (*DB, func(ctx context.Context) error, error) {
	mode := strings.ToLower(cfg.DB.Driver)
	if mode == "" {
		mode = "auto"
	}

	wantMongo := mode == "mongo" || mode == "both" || mode == "auto"
	wantMySQL := mode == "mysql" || mode == "both" || mode == "auto"

	var (
		dbs      DB
		closeFns []func(ctx context.Context) error
	)
	if wantMongo && cfg.DB.Mongo.URI != "" {
		conn, err := mongo.Connect(cfg.DB.Mongo)
		if err != nil {
			return nil, nil, err
		}
		dbs.Mongo = conn
		closeFns = append(closeFns, conn.Close)
		log.Println("Connected to MongoDB ...")
	} else if mode == "mongo" || mode == "both" {
		return nil, nil, fmt.Errorf("mongo is required but MONGO_URI is empty")
	}

	if wantMySQL && cfg.DB.MySQL.DSN != "" {
		conn, err := mysql.Connect(cfg.DB.MySQL)
		if err != nil {
			return nil, nil, err
		}

		dbs.MySQL = conn
		closeFns = append(closeFns, conn.Close)
		log.Println("Connected to MySQL ...")
	} else if mode == "mysql" || mode == "both" {
		return nil, nil, fmt.Errorf("mysql is required but MYSQL_DSN is empty")
	}

	if wantMySQL && cfg.DB.MySQL.ReplicaDSN != "" {
		replicaCfg := cfg.DB.MySQL
		replicaCfg.DSN = cfg.DB.MySQL.ReplicaDSN
		conn, err := mysql.Connect(replicaCfg)
		if err != nil {
			return nil, nil, err
		}
		dbs.MySQLReplica = conn
		closeFns = append(closeFns, conn.Close)
	}

	if dbs.Mongo == nil && dbs.MySQL == nil {
		return nil, nil, fmt.Errorf("no database configured (set MONGO_URI and/or MYSQL_DSN)")
	}

	closeAll := func(ctx context.Context) error {
		for _, fn := range closeFns {
			_ = fn(ctx)
		}
		return nil
	}

	dbs.Driver = mode
	return &dbs, closeAll, nil
}

func buildRepos(db *DB) (Repos, error) {
	if db.MySQL != nil {
		return Repos{User: repository.NewUserRepo(db.MySQL.DB)}, nil
	}
	if db.Mongo != nil {
		return Repos{User: mongo.NewUserRepo(db.Mongo.DB)}, nil
	}
	return Repos{}, fmt.Errorf("no database available for user repo")
}

func buildServices(repos Repos) Services {
	return Services{
		User: service.NewUserService(repos.User),
	}
}

func buildHandlers(services Services, logger *slog.Logger) HandlerSet {
	return HandlerSet{
		HTTP: handlers.VersionedSet{
			V1: handlers.Set{
				User: handlers.NewUserHandler(services.User, logger),
			},
			V2: handlers.Set{
				User: nil, // template for v2 handler
			},
		},
	}
}
