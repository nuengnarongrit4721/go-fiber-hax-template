package in

import (
	"context"
	d "gofiber-hax/internal/core/domain"
)

type UserService interface {
	GetByID(ctx context.Context, id string) (d.Users, error)
}
