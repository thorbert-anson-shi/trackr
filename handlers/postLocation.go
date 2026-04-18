package handlers

import (
	"time"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostLocationRequest struct {
	Latitude  float32 `json:"latitude" validate:"required"`
	Longitude float32 `json:"longitude" validate:"required"`
	Accuracy  float32 `json:"accuracy" validate:"required"`
}

func PostLocationHandler(c fiber.Ctx) error {
	queries := generated.New(database.DB)
	locationBuf := new(PostLocationRequest)

	if err := c.Bind().Body(locationBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(422)
	}

	apiKey, err := extractors.FromAuthHeader("Bearer").Extract(c)
	if err != nil {
		logging.GlobalLogger.Println("Cannot extract token from auth header")
		return c.SendStatus(401)
	}

	hashedApiKey := hashing.HashSHA256(apiKey)

	user, err := queries.GetUserByApiKey(c, pgtype.Text{String: hashedApiKey, Valid: true})
	if err != nil {
		logging.GlobalLogger.Println("API key doesn't exist. This may be an attack")
		return c.SendStatus(404)
	}

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
