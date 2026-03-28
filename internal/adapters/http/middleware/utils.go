package middleware

import "github.com/gofiber/fiber/v2"

var SkipBackgroundTasks = func(c *fiber.Ctx) bool {
	if c.Method() == fiber.MethodGet {
		switch c.Path() {
		case "/api/v1/health",
			"/api/v1/ready",
			"/api/v2/health",
			"/api/v2/ready",
			"/api/.well-known/jwks.json":
			return true
		}
	}
	return false
}
