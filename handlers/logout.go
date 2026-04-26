package handlers

import (
	"tobtoby/trackr/auth"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

func Logout(c fiber.Ctx) error {
	// API key validation middleware handles user authentication
	user := c.Locals(auth.UserContextKey).(generated.User)
	logging.GlobalLogger.Printf("User %s has logged out", user.Name.String)
	return c.SendStatus(fiber.StatusNoContent)
}
