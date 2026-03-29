package mongo

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"gofiber-hax/internal/infra/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Connector struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func Connect(cfg config.MongoConfig) (*Connector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	URI := strings.TrimSpace(cfg.URI)
	if URI == "" {
		credentials := ""
		if cfg.User != "" {
			credentials = url.QueryEscape(cfg.User)
			if cfg.Pass != "" {
				credentials += ":" + url.QueryEscape(cfg.Pass)
			}
			credentials += "@"
		}
		URI = fmt.Sprintf("mongodb://%s%s:%s/%s", credentials, cfg.Host, cfg.Port, cfg.DB)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		return nil, err
	}
	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, err
	}

	db := client.Database(cfg.DB)
	connector := &Connector{
		Client: client,
		DB:     db,
	}

	if err := EnsureIndexes(connector.DB); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("mongo auto index failed: %w", err)
	}

	return connector, nil
}

func (c *Connector) Close(ctx context.Context) error {
	return c.Client.Disconnect(ctx)
}
