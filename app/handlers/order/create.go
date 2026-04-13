package order

import (
	appctx "mini-exchange/app/ctx"
	"mini-exchange/app/presentation"

	"github.com/gofiber/fiber/v2"
)

func (h *orderHandler) Create(c *fiber.Ctx) error {

	requestID := appctx.GetRequestID(c.UserContext())
	if requestID == "" {
		if val := c.Locals("request_id"); val != nil {
			if str, ok := val.(string); ok {
				requestID = str
			}
		}
	}

	var req presentation.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		// Don't expose raw parsing error details
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Invalid request body", requestID))
	}

	if err := validate.Struct(req); err != nil {
		// Return user-friendly validation message without exposing field details
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Validation failed: please check your input", requestID))
	}

	err := h.orderService.Create(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(
		presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", nil, requestID))

}