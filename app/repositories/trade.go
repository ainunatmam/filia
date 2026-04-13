package repositories

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
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

	libraries.Logger.Debug().Str("buy_order_id", trade.BuyOrderId).Str("sell_order_id", trade.SellOrderId).Msg("[TradeRepository] Saving trade")
	r.trades = append(r.trades, trade)
	libraries.Logger.Info().Str("stock_code", trade.StockCode).Int("price", trade.Price).Int("quantity", trade.Quantity).Msg("[TradeRepository] Trade saved successfully")
	return nil
}

func (r *tradeRepo) GetAll() []*entity.TradeItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	libraries.Logger.Debug().Int("count", len(r.trades)).Msg("[TradeRepository] Getting all trades")
	return r.trades
}