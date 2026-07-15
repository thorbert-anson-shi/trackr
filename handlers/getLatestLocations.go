package handlers

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

// ListLatestLocationsHandler returns all location records.
// @Summary      List every member's latest locations
// @Description  Retrieves all latest location records from the database.
// @Tags         locations
// @Produce      json
// @Security     ApiKeyAuth
// @Success      200  {array}  LocationResponse
// @Failure      500  "Internal Server Error"
// @Router       /api/v1/users/locations/latest [get]
func ListLatestLocationsHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)

	locations, err := queries.ListLatestLocations(c.Context())
	if err != nil {
		logging.GlobalLogger.Println("An error occurred when fetching from DB")
		return err
	}

	response := make([]LocationResponse, len(locations))
	for i, l := range locations {
		response[i] = locationToResponse(l)
	}

	return c.JSON(response)
}
