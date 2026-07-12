package handlers

import (
	"tobtoby/trackr/auth"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

// Login logs the login event for the authenticated user.
// @Summary      User login
// @Description  Authenticates the user via API key and logs the login event.
// @Tags         auth
// @Security     ApiKeyAuth
// @Success      200
// @Failure      401  "Unauthorized - invalid or missing API key"
// @Router       /api/v1/auth/login [post]
func Login(c fiber.Ctx) error {
	// It's funny but the apiKeyValidator already
	// does the querying job here.
	user := c.Locals(auth.UserContextKey).(generated.User)
	logging.GlobalLogger.Printf("User %s has logged in", user.Name)
	return c.SendStatus(200)
}
