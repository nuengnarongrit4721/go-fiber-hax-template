package mongo

import (
	"context"

	m "gofiber-hax/internal/adapters/db/mongo/models"
	d "gofiber-hax/internal/core/domain"
	coreerrors "gofiber-hax/internal/shared/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	col *mongo.Collection
}

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{col: db.Collection("users")}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (d.Users, error) {
	filter := bson.M{"_id": id}
	if oid, err := primitive.ObjectIDFromHex(id); err == nil {
		filter = bson.M{"_id": oid}
	}

	var doc m.Users
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return d.Users{}, coreerrors.ErrNotFound
		}
		return d.Users{}, err
	}

	return d.Users{
		AccountID: doc.AccountID,
		Fname:     doc.Fname,
		Lname:     doc.Lname,
		FullName:  doc.FullName,
		Username:  doc.Username,
		Password:  doc.Password,
		Email:     doc.Email,
		Phone:     doc.Phone,
	}, nil
}
