package auth

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/jackc/pgx/v5/pgtype"
)

func ApiKeyValidator(c fiber.Ctx, key string) (bool, error) {
	queries := generated.New(database.DB)

	hashedKey := hashing.HashSHA256(key)

	_, err := queries.GetUserByApiKey(c, pgtype.Text{String: hashedKey, Valid: true})
	if err != nil {
		return false, keyauth.ErrMissingOrMalformedAPIKey
	}

	return true, nil
}
