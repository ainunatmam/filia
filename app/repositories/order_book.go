package repositories

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
	"sync"
)

type orderBookRepo struct{
	orderBook map[string]*entity.OrderBook 
	mu    sync.RWMutex
}

func NewOrderBookRepository() OrderBookRepository {
	return &orderBookRepo{
		orderBook: make(map[string]*entity.OrderBook),

	}
}

func (r *orderBookRepo) GetBookByStockCode(ctx context.Context, code string) (*entity.OrderBook, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	libraries.Logger.Debug().Str("stock_code", code).Msg("[OrderBookRepository] Getting order book")

	book, exists := r.orderBook[code]
	if !exists {
		libraries.Logger.Debug().Str("stock_code", code).Msg("[OrderBookRepository] Order book not found, creating new")
		book = &entity.OrderBook{
			Buys:  []*entity.OrderItem{},
			Sells: []*entity.OrderItem{},
		}
		r.orderBook[code] = book
	}
	return book, nil
}

func (r *orderBookRepo) Save(ctx context.Context, stockCode string, book *entity.OrderBook) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	libraries.Logger.Debug().Str("stock_code", stockCode).Int("buys", len(book.Buys)).Int("sells", len(book.Sells)).Msg("[OrderBookRepository] Saving order book")
	r.orderBook[stockCode] = book
	return nil
}