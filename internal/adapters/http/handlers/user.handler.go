package handlers

import (
	"log/slog"

	"gofiber-hax/internal/core/ports/in"
	coreerrors "gofiber-hax/internal/shared/errors"
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

func (h *UserHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.uc.GetByID(c.Context(), id)
	if err != nil {
		if err == coreerrors.ErrNotFound {
			return response.Error(c, fiber.StatusNotFound, "user not found")
		}

		h.log.Error("get user failed", "error", err)
		return response.Error(c, fiber.StatusInternalServerError, "internal error")
	}

	return response.JSON(c, fiber.StatusOK, user)
}
