package mongo

import (
	"context"

	"gofiber-hax/internal/core/domain"
	coreerrors "gofiber-hax/internal/shared/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	col *mongo.Collection
}

type userDoc struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name"`
	Email string             `bson:"email"`
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{col: db.Collection("users")}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	filter := bson.M{"_id": id}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid}
	}

	var doc userDoc
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.User{}, coreerrors.ErrNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:    doc.ID.Hex(),
		Name:  doc.Name,
		Email: doc.Email,
	}, nil
}
