package handlers

import (
	"strings"

	"tobtoby/trackr/auth"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sethvargo/go-diceware/diceware"
)

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required"`
	RegistrationToken string `json:"registration_token" validate:"required"`
}

// CreateUser creates a new user with a generated API key.
// @Summary      Create a new user
// @Description  Admin-only endpoint that creates a new user with a randomly generated API key.
// @Tags         users
// @Accept       json
// @Produce      plain
// @Security     ApiKeyAuth
// @Param        body  body  CreateUserRequest  true  "User details"
// @Success      200   {string}  string  "Generated API key"
// @Failure      403   "Forbidden - user is not admin"
// @Failure      422   "Unprocessable Entity - invalid request body"
// @Failure      500   "Internal Server Error"
// @Router       /api/v1/users [post]
func CreateUser(c fiber.Ctx) error {
	queries := generated.New(database.DB)

	// Fetch user extracted by API key validator
	user := c.Locals(auth.UserContextKey).(generated.User)

	if !user.IsAdmin.Valid || !user.IsAdmin.Bool {
		logging.GlobalLogger.Println("User is not admin, or admin data is invalid")
		return c.SendStatus(403)
	}

	userBuf := new(CreateUserRequest)

	if err := c.Bind().Body(userBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(422)
	}

	apiKey, err := generatePassphrase(4)
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
