package routes

import (
	"blogpost/config"
	"blogpost/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRouter(app *fiber.App) {

	auth := app.Group("/auth")
	auth.Post(config.LoginRoute, handlers.Login)
	auth.Post(config.SignUpRoute, handlers.SignUp)
}
