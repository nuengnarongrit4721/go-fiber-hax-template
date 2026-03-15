package out

import (
	"context"
	d "gofiber-hax/internal/core/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *d.Users) error
	GetByAccountID(ctx context.Context, AccountID string) (d.Users, error)
}
