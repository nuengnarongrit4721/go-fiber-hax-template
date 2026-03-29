package app

import (
	"context"

	"gofiber-hax/internal/adapters/db/mongo"
	"gofiber-hax/internal/adapters/db/mysql"
	"gofiber-hax/internal/infra/config"
)

func Migrate(cfg config.Config) error {
	migrationCfg := cfg
	migrationCfg.DB.MySQL.AutoMigrate = false

	dbs, closeDB, err := buildDB(migrationCfg)
	if err != nil {
		return err
	}
	defer func() {
		if closeDB != nil {
			_ = closeDB(context.Background())
		}
	}()

	if dbs.MySQL != nil {
		if err := mysql.Migrate(dbs.MySQL.DB); err != nil {
			return err
		}
	}
	if dbs.Mongo != nil {
		if err := mongo.EnsureIndexes(dbs.Mongo.DB); err != nil {
			return err
		}
	}

	return nil
}
