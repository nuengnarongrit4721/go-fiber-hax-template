package middleware

import (
	"errors"
	"strings"

	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/infra/jwt"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
)

func Auth(cfg config.AuthConfig) fiber.Handler {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = "token"
	}

	var jwtVal *jwt.Validator
	if mode == "jwt" || mode == "google" {
		jwtVal = jwt.NewValidator(cfg.JWT)
	}

	return func(c *fiber.Ctx) error {
		if !cfg.Enabled {
			return c.Next()
		}

		token, err := extractToken(c, cfg)
		if err != nil {
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		switch mode {
		case "token":
			// Token Auth: Validates against a static token defined in config
			if strings.TrimSpace(cfg.Token) == "" || token != strings.TrimSpace(cfg.Token) {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}
		case "jwt", "google":
			if jwtVal == nil {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}

			claims, err := jwtVal.Validate(token)
			if err != nil {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}

			c.Locals("user", claims)
		default:
			return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
		}

		return c.Next()
	}
}

func extractToken(c *fiber.Ctx, cfg config.AuthConfig) (string, error) {
	header := cfg.Header
	if header == "" {
		header = "Authorization"
	}

	value := strings.TrimSpace(c.Get(header))
	if value == "" {
		return "", errors.New("missing auth header")
	}

	scheme := strings.TrimSpace(cfg.Scheme)
	if scheme != "" {
		prefix := scheme + " "
		if len(value) <= len(prefix) || !strings.EqualFold(value[:len(prefix)], prefix) {
			return "", errors.New("invalid scheme")
		}
		value = strings.TrimSpace(value[len(prefix):])
	}

	if value == "" {
		return "", errors.New("empty token")
	}

	return value, nil
}
