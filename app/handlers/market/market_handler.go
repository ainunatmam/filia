package market

import (
	"mini-exchange/app/libraries"
	"mini-exchange/app/services/market"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MarketHandler interface {
	GetSnapshot(c *fiber.Ctx) error
}

type marketHandler struct {
	marketService market.MarketService
}

func NewMarketHandler(marketService market.MarketService) MarketHandler {
	return &marketHandler{
		marketService: marketService,
	}
}

func (h *marketHandler) GetSnapshot(c *fiber.Ctx) error {
	stockCode := c.Query("stock_code")
	if stockCode == "" {
		libraries.Logger.Warn().Msg("[MarketHandler] Missing stock_code parameter")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "stock_code query parameter is required",
		})
	}

	depthStr := c.Query("depth", "20")
	depth, _ := strconv.Atoi(depthStr)

	limitStr := c.Query("trade_limit", "50")
	limit, _ := strconv.Atoi(limitStr)

	libraries.Logger.Info().Str("stock_code", stockCode).Int("depth", depth).Int("trade_limit", limit).Msg("[MarketHandler] GetSnapshot request received")

	ctx := c.Context()

	ticker, err := h.marketService.GetTicker(ctx, stockCode)
	if err != nil {
		libraries.Logger.Error().Str("stock_code", stockCode).Err(err).Msg("[MarketHandler] Failed to get ticker")
		return err
	}

	book, err := h.marketService.GetOrderBookSnapshot(ctx, stockCode, depth)
	if err != nil {
		libraries.Logger.Error().Str("stock_code", stockCode).Err(err).Msg("[MarketHandler] Failed to get order book")
		return err
	}

	trades, err := h.marketService.GetRecentTrades(ctx, stockCode, limit)
	if err != nil {
		libraries.Logger.Error().Str("stock_code", stockCode).Err(err).Msg("[MarketHandler] Failed to get recent trades")
		return err
	}

	libraries.Logger.Info().Str("stock_code", stockCode).Msg("[MarketHandler] GetSnapshot completed successfully")
	return c.JSON(fiber.Map{
		"ticker":        ticker,
		"order_book":    book,
		"recent_trades": trades,
	})
}
