package out

import (
	"context"

	"gofiber-hax/internal/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (domain.User, error)
}
