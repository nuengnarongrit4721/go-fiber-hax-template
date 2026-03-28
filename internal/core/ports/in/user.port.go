package in

import (
	"context"
	d "gofiber-hax/internal/core/domain"
)

type UserService interface {
	CreateUserService(ctx context.Context, req *d.Users) error
	GetByAccountIDService(ctx context.Context, accountID string) (d.Users, error)
	GetUserByUsernameService(ctx context.Context, username string) (d.Users, error)
}
