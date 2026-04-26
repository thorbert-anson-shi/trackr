package handlers

import (
	"tobtoby/trackr/auth"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

func Login(c fiber.Ctx) error {
	// It's funny but the apiKeyValidator already
	// does the querying job here.
	user := c.Locals(auth.UserContextKey).(generated.User)
	logging.GlobalLogger.Printf("User %s has logged in", user.Name.String)
	return c.SendStatus(200)
}
