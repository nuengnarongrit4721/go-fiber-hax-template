package out

import (
	"context"
	d "gofiber-hax/internal/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (d.Users, error)
}
