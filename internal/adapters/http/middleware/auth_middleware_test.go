package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gofiber-hax/internal/infra/config"

	"github.com/gofiber/fiber/v2"
)

func TestAuthTokenMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(config.AuthConfig{
		Enabled: true,
		Mode:    "token",
		Token:   "secret",
		Header:  "Authorization",
		Scheme:  "Bearer",
	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	t.Run("missing token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusUnauthorized {
			t.Fatalf("expected 401, got %d", resp.StatusCode)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer secret")
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("app.Test() error = %v", err)
		}
		if resp.StatusCode != fiber.StatusNoContent {
			t.Fatalf("expected 204, got %d", resp.StatusCode)
		}
	})
}

func TestAuthMiddlewareDisabled(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(config.AuthConfig{Enabled: false}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}
