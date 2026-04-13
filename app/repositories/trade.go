package repositories

import (
	"context"
	"mini-exchange/app/entity"
	"sync"
)

type tradeRepo struct {
	trades []*entity.TradeItem
	mu     sync.RWMutex
}

func NewTradeRepository() TradeRepository {
	return &tradeRepo{
		trades: []*entity.TradeItem{},
	}
}

func (r *tradeRepo) Save(ctx context.Context, trade *entity.TradeItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.trades = append(r.trades, trade)
	return nil
}

func (r *tradeRepo) GetAll() []*entity.TradeItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.trades
}