package database

import (
	"database/sql"
	"embed"

	"tobtoby/trackr/logging"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func MigrateDB(migrationFS *embed.FS) {
	db, err := sql.Open("pgx", DB.Config().ConnString())
	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to migrate DB: %s\n", err.Error())
	}

	goose.SetBaseFS(migrationFS)

	if err := goose.SetDialect("postgres"); err != nil {
		logging.GlobalLogger.Fatalf("Failed to migrate DB: %s\n", err.Error())
	}

	if err := goose.Up(db, "postgresql/schema"); err != nil {
		logging.GlobalLogger.Fatalf("Failed to migrate DB: %s\n", err.Error())
	}
}
