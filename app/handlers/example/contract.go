package example

import (
	"github.com/gofiber/fiber/v2"
)

type ExampleHandler interface {
	Create(c *fiber.Ctx) error
}
