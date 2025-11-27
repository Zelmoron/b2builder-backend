package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	fbApp "main/firebase"
)

func (h *Handler) GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Get users"})
}

func (h *Handler) RegisterUser(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "missing authorization header",
		})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid authorization format",
		})
	}

	decodedToken, err := fbApp.AuthClient.VerifyIDToken(context.Background(), token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token",
		})
	}

	fbID := decodedToken.UID
	email := ""
	if emailClaim, ok := decodedToken.Claims["email"].(string); ok {
		email = emailClaim
	}

	user, err := h.service.RegisterUser(fbID, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to register user",
		})
	}

	return c.JSON(user)
}
