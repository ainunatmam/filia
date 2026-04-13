package trade

import (
	tradeService "mini-exchange/app/services/trade"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type tradeHandler struct {
	tradeService tradeService.TradeService
}

func NewTradeHandler(tradeService tradeService.TradeService) TradeHandler {
	return &tradeHandler{
		tradeService: tradeService,
	}
}

type TradeHandler interface {
	History(c *fiber.Ctx) error
}

var validate = validator.New()

