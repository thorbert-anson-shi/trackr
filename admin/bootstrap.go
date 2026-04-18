package admin

import (
	"context"
	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/jackc/pgx/v5/pgtype"
)

func BootstrapAdmin(c context.Context) {
	queries := generated.New(database.DB)

	adminName := config.SafeFetchVar("ADMIN_NAME")
	adminApiKey := config.SafeFetchVar("ADMIN_API_KEY")

	hashedApiKey := hashing.HashSHA256(adminApiKey)

	_, err := queries.AddUser(c, generated.AddUserParams{
		Name:              pgtype.Text{String: adminName, Valid: true},
		ApiKey:            pgtype.Text{String: hashedApiKey, Valid: true},
		RegistrationToken: pgtype.Text{},
		IsAdmin:           pgtype.Bool{Bool: true, Valid: true},
	})

	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to generate admin user: %s\n", err.Error())
	}

	logging.GlobalLogger.Println("Successfully bootstrapped admin user")
}
