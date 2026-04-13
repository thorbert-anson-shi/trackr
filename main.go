package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/handlers"
	"tobtoby/trackr/logging"
	"tobtoby/trackr/polling"
	"tobtoby/trackr/validation"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config.InitializeEnv()
	logging.InitializeLogger()
	database.ConnectDB()
	firebase.ConnectFirebase()
	polling.InitializePoller(signalCtx)

	app := fiber.New(fiber.Config{
		StructValidator: &validation.StructValidator{Validator: goValidator.New()},
	})

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")
	v1.Get("/locations", handlers.ListLocationsHandler)
	v1.Post("/locations", handlers.PostLocationHandler)

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
