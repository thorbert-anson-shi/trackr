package handlers

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

type LocationResponse struct {
	ID        int32   `json:"id"`
	UserID    int32   `json:"user_id"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Timestamp string  `json:"timestamp"`
	Accuracy  float32 `json:"accuracy"`
}

// ListLocationsHandler returns all location records.
// @Summary      List all locations
// @Description  Retrieves all location records from the database.
// @Tags         locations
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}  LocationResponse
// @Failure      500  "Internal Server Error"
// @Router       /api/v1/locations [get]
func ListLocationsHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)

	locations, err := queries.ListLocations(c.Context())
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching from DB")
		return err
	}

	return c.JSON(locations)
}
