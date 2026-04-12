package main

import (
	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/firebase"
	"tobtoby/trackr/handlers"
	"tobtoby/trackr/logging"
	"tobtoby/trackr/validation"

	goValidator "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	logging.InitializeLogger()
	config.InitializeEnv()
	database.ConnectDB()
	firebase.ConnectFirebase()

	app := fiber.New(fiber.Config{
		StructValidator: &validation.StructValidator{Validator: goValidator.New()},
	})

	app.Use(logger.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/locations", handlers.ListLocationsHandler)
	v1.Post("/locations", handlers.PostLocationHandler)

	v1.Post("/tokens", handlers.SetTokenHandler)

	app.Listen(":8000")
}
