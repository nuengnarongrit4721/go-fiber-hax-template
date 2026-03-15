package service

import (
	"context"
	"fmt"

	d "gofiber-hax/internal/core/domain"
	portsin "gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/core/ports/out"
)

type UserService struct {
	repo out.UserRepository
}

func NewUserService(repo out.UserRepository) *UserService {
	return &UserService{repo: repo}
}

var _ portsin.UserService = (*UserService)(nil)

func (s *UserService) GetByID(ctx context.Context, id string) (d.Users, error) {
	if id == "" {
		return d.Users{}, fmt.Errorf("id is required")
	}
	return s.repo.GetByID(ctx, id)
}
