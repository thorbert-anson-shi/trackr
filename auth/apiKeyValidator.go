package auth

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/jackc/pgx/v5/pgtype"
)

var UserContextKey struct{}

func APIKeyValidator(c fiber.Ctx, key string) (bool, error) {
	queries := generated.New(database.DB)

	hashedKey := hashing.HashSHA256(key)

	user, err := queries.GetUserByApiKey(c, pgtype.Text{String: hashedKey, Valid: true})
	if err != nil {
		logging.GlobalLogger.Println("The provided API key is invalid")
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	c.Locals(UserContextKey, user)

	return true, nil
}
