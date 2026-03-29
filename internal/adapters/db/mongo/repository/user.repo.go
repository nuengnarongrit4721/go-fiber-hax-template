package mongo

import (
	"context"
	"fmt"
	"time"

	m "gofiber-hax/internal/adapters/db/mongo/models"
	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/out"
	coreerrors "gofiber-hax/internal/shared/errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	col *mongo.Collection
}

var _ out.UserRepository = (*UserRepo)(nil)

func NewUserRepo(db *mongo.Database) *UserRepo {
	return &UserRepo{col: db.Collection("users")}
}

func (r *UserRepo) CreateUser(ctx context.Context, req *d.Users) error {
	mUsers := ToMongoUser(*req)
	now := time.Now().UTC()
	if mUsers.Id.IsZero() {
		mUsers.Id = primitive.NewObjectID()
	}
	if mUsers.CreatedAt.IsZero() {
		mUsers.CreatedAt = now
	}
	mUsers.UpdatedAt = now
	_, err := r.col.InsertOne(ctx, mUsers)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return coreerrors.ErrConflict
		}
		return fmt.Errorf("mongo.userrepo.create error: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByAccountID(ctx context.Context, AccountID string) (d.Users, error) {
	filter := bson.M{"account_id": AccountID}
	var doc m.Users
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return d.Users{}, coreerrors.ErrNotFound
		}
		return d.Users{}, err
	}

	return ToDomainUser(&doc), nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (d.Users, error) {
	filter := bson.M{"username": username}
	var doc m.Users
	err := r.col.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return d.Users{}, coreerrors.ErrNotFound
		}
		return d.Users{}, err
	}

	return ToDomainUser(&doc), nil
}
