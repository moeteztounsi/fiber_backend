package main

import (
	"github.com/alpha_batta/database"
	"github.com/alpha_batta/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	database.Connect()
	app := fiber.New()
	app.Use(cors.New(cors.Config{AllowCredentials: true}))

	routes.SetUp(app)

	app.Listen(":3000")
}
