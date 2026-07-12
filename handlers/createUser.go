package handlers

import (
	"errors"
	"strings"
	"time"

	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sethvargo/go-diceware/diceware"
)

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required"`
	RegistrationToken string `json:"registration_token" validate:"required"`
	InviteCode        string `json:"invite_code" validate:"required"`
}

// CreateUser creates a new user with a randomly generated API key.
// @Summary      Create a new user
// @Description  Guest-only endpoint that creates a new user when provided a valid invite code.
// @Tags         users
// @Accept       json
// @Produce      plain
// @Param        body  body  CreateUserRequest  true  "User details"
// @Success      200   {string}  string  "Generated API key"
// @Failure      400   "Bad Request"
// @Failure      403   "Forbidden - the invite link has expired"
// @Failure      403   "Forbidden - the invite link has been used"
// @Failure      422   "Unprocessable Entity - invalid request body"
// @Failure      500   "Internal Server Error"
// @Router       /api/v1/users [post]
func CreateUser(c fiber.Ctx) error {
	queries := generated.New(database.DB)
	userBuf := new(CreateUserRequest)

	if err := c.Bind().Body(userBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	invitation, err := queries.GetInvitationByCode(c, userBuf.InviteCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logging.GlobalLogger.Println("An unrecognized invite code was used")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		logging.GlobalLogger.Printf("An unknown error occurred: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if time.Now().After(invitation.ExpiryDate.Time) {
		_ = c.SendString("The invite link has expired")
		return c.SendStatus(fiber.StatusForbidden)
	}

	if invitation.IsUsed {
		_ = c.SendString("The invite link has been used")
		return c.SendStatus(fiber.StatusForbidden)
	}

	apiKey, err := generatePassphrase(4)
	if err != nil {
		logging.GlobalLogger.Printf("Failed to generate API key: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	apiKeyHash := hashing.HashSHA256(apiKey)

	// INFO: Create user and invalidate invitation within one transaction
	tx, err := database.DB.BeginTx(c, pgx.TxOptions{})
	if err != nil {
		logging.GlobalLogger.Printf("An error occurred when starting transaction: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer tx.Rollback(c)

	qtx := queries.WithTx(tx)

	_, err = qtx.AddUser(c.Context(), generated.AddUserParams{
		Name:              userBuf.Name,
		ApiKey:            apiKeyHash,
		RegistrationToken: pgtype.Text{String: userBuf.RegistrationToken, Valid: true},
		IsAdmin:           false,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			logging.GlobalLogger.Printf("User creation failed: %s; PostgreSQL error code: %s",
				pgErr.Message, pgErr.Code)
		} else {
			logging.GlobalLogger.Printf("Failed to create user: %s\n", err.Error())
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = qtx.InvalidateInvitationById(c, invitation.ID)
	if err != nil {
		logging.GlobalLogger.Println("Failed to invalidate invitation. User creation will be cancelled")
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	err = tx.Commit(c)
	if err != nil {
		logging.GlobalLogger.Printf("An error occurred when committing transaction: %s\n", err.Error())
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).SendString(apiKey)
}

func generatePassphrase(numWords int) (string, error) {
	parts, err := diceware.Generate(numWords)
	if err != nil {
		return "", err
	}

	return strings.Join(parts, "-"), nil
}
