package repository

import (
	"context"
	"errors"

	m "gofiber-hax/internal/adapters/db/mysql/models"
	d "gofiber-hax/internal/core/domain"
	coreerrors "gofiber-hax/internal/shared/errors"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (d.Users, error) {
	var model m.Users
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return d.Users{}, coreerrors.ErrNotFound
		}
		return d.Users{}, err
	}

	return toDomain(model), nil
}

func toDomain(m m.Users) d.Users {
	return d.Users{
		BaseDomain: d.BaseDomain{
			ID:        m.ID,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
			DeletedAt: m.DeletedAt,
		},
		AccountID: m.AccountID,
		Fname:     m.Fname,
		Lname:     m.Lname,
		FullName:  m.FullName,
		Username:  m.Username,
		Password:  m.Password,
		Email:     m.Email,
		Phone:     m.Phone,
	}
}
