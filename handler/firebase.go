package handler

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"

	fbApp "main/firebase"
)

func (h *Handler) FirebaseAuth(c *fiber.Ctx) error {

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

	cid, err := h.repo.GetUserByFbID(decodedToken.UID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not registered",
		})
	}
	c.Locals("id", cid.ID)
	c.Locals("fbClaims", decodedToken.Claims)

	return c.Next()

}
