package in

import (
	"context"

	d "gofiber-hax/internal/core/domain"
)

type AuthService interface {
	RegisterService(ctx context.Context, req *d.RegisterUserInput) error
}
