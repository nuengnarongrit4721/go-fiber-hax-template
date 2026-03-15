package handlers

import (
	"log/slog"

	"gofiber-hax/internal/core/ports/in"
	errs "gofiber-hax/internal/shared/errors"
	"gofiber-hax/internal/shared/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	uc  in.UserService
	log *slog.Logger
}

func NewUserHandler(uc in.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{uc: uc, log: logger}
}

func (h *UserHandler) GetByAccountIDHandler(c *fiber.Ctx) error {
	accountID := c.Params("account_id")
	if accountID == "" {
		return response.Error(c, fiber.StatusBadRequest, "account_id is required")
	}
	user, err := h.uc.GetByAccountIDService(c.Context(), accountID)
	if err != nil {
		if err == errs.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "user not found")
		}

		return response.Error(c, fiber.StatusInternalServerError, errs.ErrInternalServer.Error())
	}

	return response.JSON(c, fiber.StatusOK, ToUserResponse(user))
}
