package routes

import (
	"github.com/alpha_batta/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetUp(app *fiber.App) {

	app.Post("/api/users/register", controllers.Register)
	app.Post("/api/users/:id/login", controllers.Login)
	app.Get("/api/users/:id/user", controllers.User)
	app.Post("/api/users/:id/logout", controllers.Logout)

}
