package handlers

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

func ListLocationsHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)

	locations, err := queries.ListLocations(c.Context())
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching from DB")
		return err
	}

	return c.JSON(locations)
}
