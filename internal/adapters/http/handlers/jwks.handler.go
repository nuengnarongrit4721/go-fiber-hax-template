package handlers

import (
	"gofiber-hax/internal/infra/jwt"

	"github.com/gofiber/fiber/v2"
)

type JWKSHandler struct {
	signer *jwt.Signer
}

func NewJwksHandler(signer *jwt.Signer) *JWKSHandler {
	return &JWKSHandler{signer: signer}
}

// GetKeys ทำหน้าที่คืนค่ากุญแจแบบ JSON ส่งกลับไปให้เว็บบราวเซอร์หรือ Middleware อื่นๆ
func (h *JWKSHandler) GetKeysEndpoint(c *fiber.Ctx) error {
	return c.JSON(h.signer.GetJWKSet())
}
