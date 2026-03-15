package in

import (
	"context"
	"gofiber-hax/internal/adapters/http/dto"
)

type AuthService interface {
	RegisterService(ctx context.Context, req *dto.RegisterRequest) error
}
