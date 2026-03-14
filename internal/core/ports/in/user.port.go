package in

import (
	"context"

	"gofiber-hax/internal/core/domain"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
}
