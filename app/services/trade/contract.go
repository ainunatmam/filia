package trade

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/repositories"
)

type TradeService interface {
	History(ctx context.Context) ([]*entity.TradeItem, error)
}

type tradeService struct {
	tradeRepo repositories.TradeRepository
}

func NewTradeService(tradeRepository repositories.TradeRepository) TradeService {
	return &tradeService{
		tradeRepo: tradeRepository,
	}
}
