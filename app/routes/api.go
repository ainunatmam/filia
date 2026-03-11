package routes

import (
	"wallet-api/app/handlers/example"
	"wallet-api/app/repositories"

	exampleService "wallet-api/app/services/example"

	"github.com/gofiber/fiber/v2"
)

func (r *router) api() fiber.Router {
	root := r.app.Group("/")
	api := root.Group("/api")

	// init repo
	exampleRepo := repositories.NewExampleRepository(r.goquLibrary.DB())

	// init service
	exampleService := exampleService.NewExampleService(exampleRepo)

	// init handler
	exampleHandler := example.NewExampleHandler(exampleService)

	// Greeting
	api.Get("/greeting", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Konnichiwa !",
		})
	})

	// Example
	api.Post("/example", exampleHandler.Create)
	return root
}
