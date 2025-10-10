package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get users"})
}
