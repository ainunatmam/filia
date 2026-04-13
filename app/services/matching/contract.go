package matching

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/repositories"
	"mini-exchange/app/services/market"
	"mini-exchange/app/ws"
)

type MatchingService interface {
	Submit(ctx context.Context, order *entity.OrderItem) error
	Start(ctx context.Context)
}

type matchingService struct {
	orderCh       chan *entity.OrderItem
	orderBookRepo repositories.OrderBookRepository
	tradeRepo     repositories.TradeRepository
	wsHub         ws.Hub
	marketService market.MarketService
}

func NewMatchingService(
	orderBookRepository repositories.OrderBookRepository,
	tradeRepository repositories.TradeRepository,
	wsHub ws.Hub,
	marketService market.MarketService,
) MatchingService {
	return &matchingService{
		orderCh:       make(chan *entity.OrderItem, 1000),
		orderBookRepo: orderBookRepository,
		tradeRepo:     tradeRepository,
		wsHub:         wsHub,
		marketService: marketService,
	}
}

