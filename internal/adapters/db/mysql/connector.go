package mysql

import (
	"context"
	"time"

	"gofiber-hax/internal/infra/config"

	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Connector struct {
	DB *gorm.DB
}

func Connect(cfg config.MySQLConfig) (*Connector, error) {

	db, err := gorm.Open(mysqlDriver.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if cfg.AutoMigrate {
		if err := autoMigrate(db); err != nil {
			return nil, err
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	return &Connector{DB: db}, nil
}

func (c *Connector) Close(ctx context.Context) error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
