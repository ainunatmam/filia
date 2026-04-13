package trade

import (
	appctx "mini-exchange/app/ctx"
	"mini-exchange/app/libraries"
	"mini-exchange/app/presentation"

	"github.com/gofiber/fiber/v2"
)

func (h *tradeHandler) 	History(c *fiber.Ctx) error {
	requestID := appctx.GetRequestID(c.UserContext())
	if requestID == "" {
		if val := c.Locals("request_id"); val != nil {
			if str, ok := val.(string); ok {
				requestID = str
			}
		}
	}

	libraries.Logger.Info().Str("request_id", requestID).Msg("[TradeHandler] Trade history request received")

	result, err := h.tradeService.History(c.Context())
	if err != nil {
		libraries.Logger.Error().Str("request_id", requestID).Err(err).Msg("[TradeHandler] Failed to get trade history")
		return c.Status(fiber.StatusOK).JSON(
			presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	libraries.Logger.Info().Str("request_id", requestID).Int("count", len(result)).Msg("[TradeHandler] Trade history fetched successfully")
	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", result, requestID))

}