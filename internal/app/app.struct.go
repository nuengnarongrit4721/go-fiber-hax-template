package app

import (
	"gofiber-hax/internal/adapters/db/mongo"
	"gofiber-hax/internal/adapters/db/mysql"
	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/core/ports/out"
	"gofiber-hax/internal/core/service"
)

type DB struct {
	Driver       string
	Mongo        *mongo.Connector
	MySQL        *mysql.Connector
	MySQLReplica *mysql.Connector
}

// NOTE: DB Logic
type Repos struct {
	User out.UserRepository
}

// NOTE: Business Logic
type Services struct {
	User *service.UserService
	Auth *service.AuthService
}

// NOTE: HTTP Logic
type HandlerSet struct {
	HTTP handlers.VersionedSet
}
