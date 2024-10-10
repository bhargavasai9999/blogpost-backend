package routes

import (
	"blogpost/config"

	"github.com/gofiber/fiber/v2"
)

func SetupBlogRoutes(app *fiber.App) {
	blog := app.Group("/blog")
	blog.Get(config.GetBlogsRoute, func(c *fiber.Ctx) error {

		return c.Status(200).JSON(fiber.Map{"message": "blog get is successful"})

	})

}
