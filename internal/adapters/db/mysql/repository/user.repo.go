package repository

import (
	"context"
	"errors"
	"fmt"

	m "gofiber-hax/internal/adapters/db/mysql/models"
	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/out"
	errs "gofiber-hax/internal/shared/errors"

	mysqlerr "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

var _ out.UserRepository = (*UserRepo)(nil)

func (r *UserRepo) CreateUser(ctx context.Context, req *d.Users) error {
	mUsers := ToModelUser(req)
	if err := r.db.WithContext(ctx).Create(&mUsers).Error; err != nil {
		if isDuplicateKeyError(err) {
			return errs.ErrConflict
		}
		return fmt.Errorf("mysql.userrepo.create error: %w", err)
	}
	return nil
}

func (r *UserRepo) GetByAccountID(ctx context.Context, AccountID string) (d.Users, error) {
	var model m.Users
	err := r.db.WithContext(ctx).First(&model, "account_id = ?", AccountID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return d.Users{}, errs.ErrNotFound
		}
		return d.Users{}, err
	}
	return ToDomainUser(&model), nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (d.Users, error) {
	var model m.Users
	err := r.db.WithContext(ctx).First(&model, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return d.Users{}, errs.ErrNotFound
		}
		return d.Users{}, err
	}
	return ToDomainUser(&model), nil
}

func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysqlerr.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
