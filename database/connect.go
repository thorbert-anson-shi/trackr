package database

import (
	"context"
	"fmt"
	"strconv"
	"tobtoby/trackr/config"
	"tobtoby/trackr/logging"

	"github.com/jackc/pgx/v5"
)

// ConnectDB connect to db
func ConnectDB() {
	logging.GlobalLogger.Println("Parsing connection parameters")

	var err error
	p := config.Config["DB_PORT"]
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		logging.GlobalLogger.Fatalln("Failed to parse database port")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.Config["DB_USER"],
		config.Config["DB_PASSWORD"],
		config.Config["DB_HOST"],
		port,
		config.Config["DB_NAME"],
	)

	logging.GlobalLogger.Println("Opening connection to DB...")

	DB, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		logging.GlobalLogger.Fatalln("Failed to connect to database")
	}

	logging.GlobalLogger.Println("Connection Opened to database")

	// DB.AutoMigrate(&model.Product{}, &model.User{})
	// fmt.Println("Database Migrated")

}
