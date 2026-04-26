package main

import (
	"context"
	"embed"
	"os"
	"os/signal"
	"syscall"

	"tobtoby/trackr/admin"
	"tobtoby/trackr/auth"
	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/handlers"
	"tobtoby/trackr/logging"
	"tobtoby/trackr/polling"
	"tobtoby/trackr/validation"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/helmet"
	"github.com/gofiber/fiber/v3/middleware/keyauth"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

//go:embed postgresql/schema/*.sql
var migrationDir embed.FS

func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config.InitializeEnv()
	logging.InitializeLogger()
	database.ConnectDB()
	database.MigrateDB(&migrationDir)
	firebase.ConnectFirebase()
	polling.InitializePoller(signalCtx)

	admin.BootstrapAdmin(signalCtx)

	app := fiber.New(fiber.Config{
		StructValidator: &validation.StructValidator{Validator: goValidator.New()},
	})

	app.Use(logger.New(logger.Config{
		Stream: logging.GlobalLogger.Writer(),
	}))
	app.Use(helmet.New())
	app.Use(keyauth.New(auth.KeyAuthConfig))

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/locations", handlers.ListLocationsHandler)
	v1.Post("/locations", handlers.PostLocationHandler)
	v1.Post("/users", handlers.CreateUser)
	v1.Post("/auth/login", handlers.Login)
	v1.Post("/auth/logout", handlers.Logout)

	app.Get("/health", func(c fiber.Ctx) error { return c.SendStatus(200) })

	go func() {
		if err := app.Listen(":8000"); err != nil {
			logging.GlobalLogger.Printf("Server error: %v\n", err)
		}
	}()

	<-signalCtx.Done()
	logging.GlobalLogger.Println("Shutting down...")

	if err := app.Shutdown(); err != nil {
		logging.GlobalLogger.Printf("Shutdown error: %v\n", err)
	}

	logging.GlobalLogger.Println("Server stopped")
}
