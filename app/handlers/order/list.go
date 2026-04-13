package order

import (
	appctx "mini-exchange/app/ctx"
	"mini-exchange/app/libraries"
	"mini-exchange/app/presentation"

	"github.com/gofiber/fiber/v2"
)

func (h *orderHandler) List(c *fiber.Ctx) error {

	requestID := appctx.GetRequestID(c.UserContext())
	if requestID == "" {
		if val := c.Locals("request_id"); val != nil {
			if str, ok := val.(string); ok {
				requestID = str
			}
		}
	}

	libraries.Logger.Info().Str("request_id", requestID).Msg("[OrderHandler] List orders request received")

	result, err := h.orderService.List(c.Context())
	if err != nil {
		libraries.Logger.Error().Str("request_id", requestID).Err(err).Msg("[OrderHandler] Failed to list orders")
		return c.Status(fiber.StatusInternalServerError).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	libraries.Logger.Info().Str("request_id", requestID).Int("count", len(result)).Msg("[OrderHandler] Orders listed successfully")
	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", result, requestID))

}