package repository

import (
	"context"
	"errors"

	m "gofiber-hax/internal/adapters/db/mysql/models"
	"gofiber-hax/internal/core/domain"
	coreerrors "gofiber-hax/internal/shared/errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (domain.User, error) {
	var model m.Users
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, coreerrors.ErrNotFound
		}
		return domain.User{}, err
	}

	return domain.User{
		ID:    model.ID,
		Name:  model.Name,
		Email: model.Email,
	}, nil
}
