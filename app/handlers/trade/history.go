package trade

import (
	appctx "mini-exchange/app/ctx"
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

	result, err := h.tradeService.History(c.Context())
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Failed(fiber.StatusInternalServerError, "An internal error occurred", requestID))
	}

	return c.Status(fiber.StatusOK).JSON(
		presentation.ResponseBase{}.Success("Success", result, requestID))

}