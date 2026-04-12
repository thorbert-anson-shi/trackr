package handlers

import (
	"encoding/json"
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

	marshalled, err := json.Marshal(locations)
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when marshalling location data")
		return err
	}

	c.Response().BodyWriter().Write([]byte(marshalled))

	// When a user requests the other devices' locations,
	// the server publishes a "requestLocation" event to the queue.
	// It expects the subscribers to then send a request to the POST /locations endpoint
	return nil
}
