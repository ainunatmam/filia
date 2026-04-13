package order

import (
	appctx "mini-exchange/app/ctx"
	"mini-exchange/app/libraries"
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

	libraries.Logger.Info().Str("request_id", requestID).Msg("[OrderHandler] Create order request received")

	var req presentation.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		libraries.Logger.Warn().Str("request_id", requestID).Err(err).Msg("[OrderHandler] Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Invalid request body", requestID))
	}

	if err := validate.Struct(req); err != nil {
		libraries.Logger.Warn().Str("request_id", requestID).Err(err).Msg("[OrderHandler] Validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusBadRequest, "Validation failed: please check your input", requestID))
	}

	libraries.Logger.Debug().Str("request_id", requestID).Str("stock_code", req.StockCode).Str("side", req.Side).Int("price", req.Price).Int("quantity", req.Quantity).Msg("[OrderHandler] Processing create order")

	err := h.orderService.Create(c.Context(), req)
	if err != nil {
		libraries.Logger.Error().Str("request_id", requestID).Err(err).Msg("[OrderHandler] Failed to create order")
		return c.Status(fiber.StatusBadGateway).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	libraries.Logger.Info().Str("request_id", requestID).Msg("[OrderHandler] Order created successfully")
	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", nil, requestID))

}
