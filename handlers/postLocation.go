package handlers

import "github.com/gofiber/fiber/v3"

func PostLocationHandler(c fiber.Ctx) error {
	// The service stores the user's provided location
	// in the database. The server then publishes a "publishLocations"
	// event that the consumers then accept and process.
	return c.SendString("here location")
}
