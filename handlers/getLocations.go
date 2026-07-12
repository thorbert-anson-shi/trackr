package handlers

import (
	"time"

	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
)

type LocationResponse struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"userID"`
	Latitude  float32   `json:"latitude"`
	Longitude float32   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
	Accuracy  float32   `json:"accuracy"`
}

func locationToResponse(l generated.Location) LocationResponse {
	return LocationResponse{
		ID:        l.ID,
		UserID:    l.UserID.Int32,
		Latitude:  l.Latitude,
		Longitude: l.Longitude,
		Timestamp: l.Timestamp.Time,
		Accuracy:  l.Accuracy.Float32,
	}
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

	response := make([]LocationResponse, len(locations))
	for i, l := range locations {
		response[i] = locationToResponse(l)
	}

	return c.JSON(response)
}
