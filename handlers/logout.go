package handlers

import (
	"tobtoby/trackr/auth"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

// Logout logs the logout event for the authenticated user.
// @Summary      User logout
// @Description  Logs the logout event for the authenticated user.
// @Tags         auth
// @Security     ApiKeyAuth
// @Success      204
// @Failure      401  "Unauthorized - invalid or missing API key"
// @Router       /api/v1/auth/logout [post]
func Logout(c fiber.Ctx) error {
	// API key validation middleware handles user authentication
	user := c.Locals(auth.UserContextKey).(generated.User)
	logging.GlobalLogger.Printf("User %s has logged out", user.Name)
	return c.SendStatus(fiber.StatusNoContent)
}
