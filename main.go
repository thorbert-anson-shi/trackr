package main

import (
	"tobtoby/trackr/config"
	"tobtoby/trackr/database"
	"tobtoby/trackr/handlers"
	"tobtoby/trackr/logging"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	logging.InitializeLogger()
	config.InitializeEnv()
	database.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())

	api := app.Group("/api")

	v1 := api.Group("/v1")
	v1.Get("/locations", handlers.ListLocationsHandler)
	v1.Post("/locations", handlers.PostLocationHandler)

	app.Listen(":8000")
}
