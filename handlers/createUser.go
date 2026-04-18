package handlers

import (
	"strings"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sethvargo/go-diceware/diceware"
)

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required"`
	RegistrationToken string `json:"registration_token" validate:"required"`
}

// CreateUser is an admin-only endpoint that parses the request body,
// generates a random API key, hashes it with SHA256, and stores the
// new user via queries.AddUser().
func CreateUser(c fiber.Ctx) error {
	queries := generated.New(database.DB)

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

	if !(user.IsAdmin.Valid && user.IsAdmin.Bool) {
		logging.GlobalLogger.Println("User is not admin, or admin data is invalid")
		return c.SendStatus(403)
	}

	userBuf := new(CreateUserRequest)

	if err := c.Bind().Body(userBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(422)
	}

	apiKey, err = generatePassphrase(4)
	if err != nil {
		logging.GlobalLogger.Printf("Failed to generate passphrase: %s\n", err.Error())
		return c.SendStatus(500)
	}

	apiKeyHash := hashing.HashSHA256(apiKey)

	_, err = queries.AddUser(c.Context(), generated.AddUserParams{
		Name:              pgtype.Text{String: userBuf.Name, Valid: true},
		ApiKey:            pgtype.Text{String: apiKeyHash, Valid: true},
		RegistrationToken: pgtype.Text{String: userBuf.RegistrationToken, Valid: true},
		IsAdmin:           pgtype.Bool{Bool: false, Valid: true},
	})

	if err != nil {
		logging.GlobalLogger.Println("Failed to run query against database")
		return c.SendStatus(500)
	}

	return c.SendString(apiKey)
}

func generatePassphrase(numWords int) (string, error) {
	parts, err := diceware.Generate(numWords)
	if err != nil {
		return "", err
	}

	return strings.Join(parts, "-"), nil
}
