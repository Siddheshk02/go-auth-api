package main

import (
	"github.com/Siddheshk02/go-auth-api/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/auth", api.Auth)

	app.Post("/user", api.User)

	app.Listen(":3000")
}
