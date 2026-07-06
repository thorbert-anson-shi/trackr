package handlers

import (
	"time"

	"tobtoby/trackr/auth"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostLocationRequest struct {
	Latitude  float32 `json:"latitude" validate:"required"`
	Longitude float32 `json:"longitude" validate:"required"`
	Accuracy  float32 `json:"accuracy" validate:"required"`
}

// PostLocationHandler creates a new location record.
// @Summary      Create a location
// @Description  Records a new location for the authenticated user.
// @Tags         locations
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        body  body  PostLocationRequest  true  "Location data"
// @Success      201   {object}  LocationResponse
// @Failure      422   "Unprocessable Entity - invalid request body"
// @Failure      500   "Internal Server Error"
// @Router       /api/v1/locations [post]
func PostLocationHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)
	locationBuf := new(PostLocationRequest)

	if err := c.Bind().Body(locationBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(422)
	}

	// Fetch user extracted by API key validator
	user := c.Locals(auth.UserContextKey).(generated.User)

	location, err := queries.AddLocation(c, generated.AddLocationParams{
		UserID:    pgtype.Int4{Int32: user.ID, Valid: true},
		Latitude:  pgtype.Float4{Float32: locationBuf.Latitude, Valid: true},
		Longitude: pgtype.Float4{Float32: locationBuf.Longitude, Valid: true},
		Timestamp: pgtype.Timestamp{Time: time.Now().UTC(), Valid: true},
		Accuracy:  pgtype.Float4{Float32: locationBuf.Accuracy, Valid: true},
	})
	if err != nil {
		logging.GlobalLogger.Printf("An error occurred when posting location: %s\n", err.Error())
		return c.SendStatus(500)
	}

	return c.JSON(location)
}
