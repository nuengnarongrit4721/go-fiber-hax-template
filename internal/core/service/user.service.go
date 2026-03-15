package service

import (
	"context"
	"fmt"

	d "gofiber-hax/internal/core/domain"
	portsin "gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/core/ports/out"
	"gofiber-hax/internal/infra/logs"
)

type UserService struct {
	repo out.UserRepository
}

func NewUserService(repo out.UserRepository) *UserService {
	return &UserService{repo: repo}
}

var _ portsin.UserService = (*UserService)(nil)

func (s *UserService) CreateUserService(ctx context.Context, req *d.Users) error {
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

	logs.Debug(result)
	return result, nil
}
