package response

import "github.com/gofiber/fiber/v2"

type Envelope struct {
	Data  any          `json:"data,omitempty"`
	Error *ErrorDetail `json:"error,omitempty"`
}

type ErrorDetail struct {
	Message string `json:"message"`
}

func JSON(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(Envelope{Data: data})
}

func Error(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(Envelope{Error: &ErrorDetail{Message: msg}})
}
