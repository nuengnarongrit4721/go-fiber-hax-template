package handlers

import (
	"errors"
	"gofiber-hax/internal/adapters/http/dto"
	d "gofiber-hax/internal/core/domain"
	"gofiber-hax/internal/core/ports/in"
	"gofiber-hax/internal/infra/logs"
	"gofiber-hax/internal/infra/validation"
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

func (h *AuthHandler) LoginEndpoint(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		logs.Error(err)
		return response.Error(c, fiber.StatusBadRequest, errs.ErrInvalidInput.Error())
	}

	input := &d.LoginUserInput{
		Username: req.Username,
		Password: req.Password,
	}

	token, err := h.as.LoginService(c.UserContext(), input)
	if err != nil {
		logs.Error(err)
		if errors.Is(err, errs.ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, errs.ErrInvalidInput.Error())
		}
		if errors.Is(err, errs.ErrUnauthorized) {
			return response.Error(c, fiber.StatusUnauthorized, errs.ErrUnauthorized.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, errs.ErrInternalServer.Error())
	}

	return response.JSON(c, fiber.StatusOK, fiber.Map{
		"access_token": token,
	})
}

func (h *AuthHandler) RegisterEndpoint(c *fiber.Ctx) error {
	var req dto.RegisterRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		logs.Error(err)
		return response.Error(c, fiber.StatusBadRequest, errs.ErrInvalidInput.Error())
	}

	input := &d.RegisterUserInput{
		Fname:           req.Fname,
		Lname:           req.Lname,
		Username:        req.Username,
		Email:           req.Email,
		Phone:           req.Phone,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	}

	if err := h.as.RegisterService(c.UserContext(), input); err != nil {
		logs.Error(err)
		if errors.Is(err, errs.ErrInvalidInput) {
			return response.Error(c, fiber.StatusBadRequest, errs.ErrInvalidInput.Error())
		}
		if errors.Is(err, errs.ErrConflict) {
			return response.Error(c, fiber.StatusConflict, errs.ErrConflict.Error())
		}
		return response.Error(c, fiber.StatusInternalServerError, errs.ErrInternalServer.Error())
	}

	return response.JSON(c, fiber.StatusOK, "User registered successfully")
}
