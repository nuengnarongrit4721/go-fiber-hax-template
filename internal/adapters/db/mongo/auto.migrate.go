package mongo

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func autoIndex(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexes := map[string][]mongo.IndexModel{
		"users": {
			{
				Keys:    bson.D{{Key: "username", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		// ADD MORE COLLECTIONS AND INDEXES HERE as the project grows
		// "products": { ... },
	}

	for colName, indexModels := range indexes {
		col := db.Collection(colName)
		names, err := col.Indexes().CreateMany(ctx, indexModels)
		if err != nil {
			return fmt.Errorf("failed to create indexes for collection %s: %w", colName, err)
		}

		for _, name := range names {
			slog.Debug("mongo index ensured", "collection", colName, "index", name)
		}
	}

	return nil
}
