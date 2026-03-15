package handlers

import (
	"gofiber-hax/internal/adapters/http/dto"
	"gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/infra/logs"
	errs "gofiber-hax/internal/shared/errors"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	as in.AuthService
}

func NewAuthHandler(as in.AuthService) *AuthHandler {
	return &AuthHandler{as: as}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	return nil
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		logs.Error(err)
		return response.Error(c, fiber.StatusBadRequest, errs.ErrInvalidInput.Error())
	}

	if err := h.as.RegisterService(c.UserContext(), &req); err != nil {
		logs.Error(err)
		return response.Error(c, fiber.StatusInternalServerError, errs.ErrInternalServer.Error())
	}

	return response.JSON(c, fiber.StatusOK, "User registered successfully")
}
