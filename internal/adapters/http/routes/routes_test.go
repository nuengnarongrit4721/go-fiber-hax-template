package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gofiber-hax/internal/adapters/http/handlers"
	"gofiber-hax/internal/adapters/http/middleware"
	d "gofiber-hax/internal/core/domain"
	portsin "gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
)

type fakeAuthService struct{}

func (fakeAuthService) RegisterService(ctx context.Context, req *d.RegisterUserInput) error {
	return nil
}

func (fakeAuthService) LoginService(ctx context.Context, req *d.LoginUserInput) (string, error) {
	return "token", nil
}

type fakeUserService struct{}

func (fakeUserService) CreateUserService(ctx context.Context, req *d.Users) error {
	return nil
}

func (fakeUserService) GetByAccountIDService(ctx context.Context, accountID string) (d.Users, error) {
	return d.Users{AccountID: accountID, Email: "user@example.com"}, nil
}

func (fakeUserService) GetUserByUsernameService(ctx context.Context, username string) (d.Users, error) {
	return d.Users{AccountID: "0001", Email: "user@example.com"}, nil
}

func TestRegisterWiresPublicAndProtectedRoutes(t *testing.T) {
	var authSvc portsin.AuthService = fakeAuthService{}
	var userSvc portsin.UserService = fakeUserService{}

	app := fiber.New()
	Register(app, handlers.VersionedSet{
		V1: handlers.Set{
			Auth: handlers.NewAuthHandler(authSvc),
			User: handlers.NewUserHandler(userSvc, nil),
		},
	}, Options{
		Versions: []string{"v1"},
		Protected: []fiber.Handler{
			middleware.Auth(config.AuthConfig{
				Enabled: true,
				Mode:    "token",
				Token:   "secret",
				Header:  "Authorization",
				Scheme:  "Bearer",
			}),
		},
	})

	t.Run("health is public", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("register is public", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"fname":            "John",
			"lname":            "Doe",
			"username":         "johndoe",
			"email":            "john@example.com",
			"phone":            "0812345678",
			"password":         "password123",
			"confirm_password": "password123",
		})
		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("protected route requires auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/0001", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("protected route works with auth", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/0001", nil)
		req.Header.Set("Authorization", "Bearer secret")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusOK {
			t.Fatalf("expected 200, got %d", resp.StatusCode)
		}
	})

	t.Run("placeholder route is not registered", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer secret")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusNotFound {
			t.Fatalf("expected 404, got %d", resp.StatusCode)
		}
	})
}
