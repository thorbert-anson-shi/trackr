package database

import (
	"context"
	"fmt"
	"strconv"

	"tobtoby/trackr/config"
	"tobtoby/trackr/logging"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() {
	logging.GlobalLogger.Println("Parsing connection parameters")

	var err error
	p := config.SafeFetchVar("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		logging.GlobalLogger.Fatalln("Failed to parse database port")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.SafeFetchVar("DB_USER"),
		config.SafeFetchVarFromFile("DB_PASSWORD_FILE"),
		config.SafeFetchVar("DB_HOST"),
		port,
		config.SafeFetchVar("DB_NAME"),
	)

	logging.GlobalLogger.Println("Opening connection to DB...")

	DB, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to connect to database: %s\n", err.Error())
	}

	logging.GlobalLogger.Println("Connection Opened to database")
}
