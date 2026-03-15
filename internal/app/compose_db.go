package app

import (
	"context"
	"fmt"
	"log/slog"

	"gofiber-hax/internal/adapters/db/mongo"
	"gofiber-hax/internal/adapters/db/mysql"
	"gofiber-hax/internal/infra/config"
)

func buildDB(cfg config.Config, logger *slog.Logger) (*DB, func(ctx context.Context) error, error) {
	mode := normalizeMode(cfg.DB.Driver)
	dbs := &DB{Driver: mode}
	closers := closeGroup{}

	if err := connectMongo(cfg, logger, mode, dbs, &closers); err != nil {
		return nil, nil, err
	}
	if err := connectMySQL(cfg, logger, mode, dbs, &closers); err != nil {
		return nil, nil, err
	}
	if err := connectMySQLReplica(cfg, logger, mode, dbs, &closers); err != nil {
		return nil, nil, err
	}

	if dbs.Mongo == nil && dbs.MySQL == nil {
		return nil, nil, fmt.Errorf("no database configured (set MONGO_URI and/or MYSQL_DSN)")
	}

	return dbs, closers.close, nil
}

func connectMongo(cfg config.Config, logger *slog.Logger, mode string, dbs *DB, closers *closeGroup) error {
	if !wantMongo(mode) {
		return nil
	}
	if cfg.DB.Mongo.URI == "" {
		if mode == "mongo" || mode == "both" {
			return fmt.Errorf("mongo is required but MONGO_URI is empty")
		}
		return nil
	}
	conn, err := mongo.Connect(cfg.DB.Mongo)
	if err != nil {
		return err
	}
	dbs.Mongo = conn
	closers.add(conn.Close)
	logger.Info("Connected to MongoDB")
	return nil
}

func connectMySQL(cfg config.Config, logger *slog.Logger, mode string, dbs *DB, closers *closeGroup) error {
	if !wantMySQL(mode) {
		return nil
	}
	if cfg.DB.MySQL.DSN == "" {
		if mode == "mysql" || mode == "both" {
			return fmt.Errorf("mysql is required but MYSQL_DSN is empty")
		}
		return nil
	}
	conn, err := mysql.Connect(cfg.DB.MySQL)
	if err != nil {
		return err
	}
	dbs.MySQL = conn
	closers.add(conn.Close)
	logger.Info("Connected to MySQL")
	return nil
}

func connectMySQLReplica(cfg config.Config, logger *slog.Logger, mode string, dbs *DB, closers *closeGroup) error {
	if !wantMySQL(mode) || cfg.DB.MySQL.ReplicaDSN == "" {
		return nil
	}
	replicaCfg := cfg.DB.MySQL
	replicaCfg.DSN = cfg.DB.MySQL.ReplicaDSN
	conn, err := mysql.Connect(replicaCfg)
	if err != nil {
		return err
	}
	dbs.MySQLReplica = conn
	closers.add(conn.Close)
	logger.Info("Connected to MySQL Replica")
	return nil
}
