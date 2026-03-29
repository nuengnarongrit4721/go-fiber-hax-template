package service

import (
	"context"
	"fmt"
	"strings"

	d "gofiber-hax/internal/core/domain"
	portsin "gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/core/ports/out"
	errs "gofiber-hax/internal/shared/errors"

	"github.com/google/uuid"
)

type UserService struct {
	repo out.UserRepository
}

func NewUserService(repo out.UserRepository) *UserService {
	return &UserService{repo: repo}
}

var _ portsin.UserService = (*UserService)(nil)

func (s *UserService) CreateUserService(ctx context.Context, req *d.Users) error {
	if req == nil {
		return fmt.Errorf("userservice.create error: %w", errs.ErrInvalidInput)
	}
	if strings.TrimSpace(req.AccountID) == "" {
		req.AccountID = uuid.NewString()
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	if err := s.repo.CreateUser(ctx, req); err != nil {
		return fmt.Errorf("userservice.create error: %w", err)
	}
	return nil
}

func (s *UserService) GetByAccountIDService(ctx context.Context, accountID string) (d.Users, error) {
	result, err := s.repo.GetByAccountID(ctx, accountID)
	if err != nil {
		return d.Users{}, fmt.Errorf("userservice.GetByAccountIDService error: %w", err)
	}
	return result, nil
}

func (s *UserService) GetUserByUsernameService(ctx context.Context, username string) (d.Users, error) {
	result, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return d.Users{}, fmt.Errorf("userservice.GetUserByUsernameService error: %w", err)
	}
	return result, nil
}
