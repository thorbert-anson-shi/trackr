package handlers

import (
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateUserRequest struct {
	Name              string `json:"name" validate:"required"`
	RegistrationToken string `json:"registration_token" validate:"required"`
}

func CreateUser(c fiber.Ctx) error {
	queries := generated.New(database.DB)
	userBuf := new(CreateUserRequest)

	if err := c.Bind().Body(userBuf); err != nil {
		logging.GlobalLogger.Println("Provided request body is invalid")
		return c.SendStatus(422)
	}

	sqlc_user, err := queries.AddUser(c.Context(), generated.AddUserParams{
		Name:              pgtype.Text{String: userBuf.Name},
		RegistrationToken: pgtype.Text{String: userBuf.RegistrationToken},
	})

	if err != nil {
		logging.GlobalLogger.Println("Failed to run query against database")
		return c.SendStatus(500)
	}

	return c.JSON(sqlc_user)
}
