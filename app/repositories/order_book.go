package repositories

import (
	"context"
	"mini-exchange/app/entity"
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

	book, exists := r.orderBook[code]
	if !exists {
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

	r.orderBook[stockCode] = book
	return nil
}