package middleware

import (
	"errors"
	"strings"

	"gofiber-hax/internal/infra/config"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
)

func Auth(cfg config.AuthConfig) fiber.Handler {
	mode := strings.ToLower(strings.TrimSpace(cfg.Mode))
	if mode == "" {
		mode = "token"
	}

	var jwtVal *jwtValidator
	if mode == "jwt" || mode == "google" {
		jwtVal = newJWTValidator(cfg.JWT)
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
			if strings.TrimSpace(cfg.Token) == "" || token != strings.TrimSpace(cfg.Token) {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}
		case "jwt", "google":
			if jwtVal == nil {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}
			if err := jwtVal.Validate(token); err != nil {
				return response.Error(c, fiber.StatusUnauthorized, "unauthorized")
			}
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
