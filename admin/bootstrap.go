// Package admin defines admin bootstrapping logic
package admin

import (
	"context"
	"errors"

	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/generated"
	"tobtoby/trackr/hashing"
	"tobtoby/trackr/logging"

	"github.com/jackc/pgx/v5/pgconn"
)

func BootstrapAdmin(c context.Context) {
	queries := generated.New(database.DB)

	adminName := config.SafeFetchVar("ADMIN_NAME")
	adminAPIKey := config.SafeFetchVarFromFile("ADMIN_API_KEY_FILE")

	hashedAPIKey := hashing.HashSHA256(adminAPIKey)

	_, err := queries.AddUser(c, generated.AddUserParams{
		Name:       adminName,
		ApiKey:     hashedAPIKey,
		FirebaseID: "dummy-firebase-id",
		IsAdmin:    true,
	})
	if err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
			switch pgErr.Code {
			case "23505":
				logging.GlobalLogger.Println("Admin user already exists. Skipping admin bootstrapping")
				return
			default:
				logging.GlobalLogger.Fatalf("Failed to generate admin user: %s; PostgreSQL error code: %s\n", pgErr.Message, pgErr.Code)
			}
		}

		logging.GlobalLogger.Fatalf("An unexpected error occurred: %s\n", err.Error())
	}

	logging.GlobalLogger.Println("Successfully bootstrapped admin user")
}
