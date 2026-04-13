package order

import (
	orderService "mini-exchange/app/services/order"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type orderHandler struct {
	orderService orderService.OrderService
}

func NewOrderHandler(orderService orderService.OrderService) OrderHandler {
	return &orderHandler{
		orderService: orderService,
	}
}

type OrderHandler interface {
	Create(c *fiber.Ctx) error
	List(c *fiber.Ctx) error

}

var validate = validator.New()

