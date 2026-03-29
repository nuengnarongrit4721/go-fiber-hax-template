package middleware

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/jwt"

	"github.com/gofiber/fiber/v2"
	gojwt "github.com/golang-jwt/jwt/v5"
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

func TestAuthJWTMiddleware(t *testing.T) {
	priv, pub, err := jwt.LoadOrGenerateKeys(filepath.Join(t.TempDir(), "jwt_private.pem"))
	if err != nil {
		t.Fatalf("LoadOrGenerateKeys() error = %v", err)
	}
	signer := jwt.NewSigner(priv, pub, "test-key")
	validator := jwt.NewValidator(config.JWTConfig{
		Issuer:      "gofiber-hax",
		Audience:    "gofiber-hax",
		AllowedAlgs: []string{"RS256"},
		ClockSkew:   time.Second,
	}, jwt.WithStaticPublicKey(signer.PublicKey()))

	token, err := signer.Sign(gojwt.MapClaims{
		"iss": "gofiber-hax",
		"aud": "gofiber-hax",
		"sub": "user-1",
		"exp": time.Now().Add(time.Minute).Unix(),
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	app := fiber.New()
	app.Use(Auth(config.AuthConfig{
		Enabled: true,
		Mode:    "jwt",
		Header:  "Authorization",
		Scheme:  "Bearer",
		JWT: config.JWTConfig{
			Issuer:      "gofiber-hax",
			Audience:    "gofiber-hax",
			AllowedAlgs: []string{"RS256"},
			ClockSkew:   time.Second,
		},
	}, WithJWTValidator(validator)))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if resp.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
}
