package market

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/repositories"
	"sync"
)

type MarketService interface {
	GetTicker(ctx context.Context, stockCode string) (*entity.Ticker, error)
	GetOrderBookSnapshot(ctx context.Context, stockCode string, depth int) (*entity.OrderBook, error)
	GetRecentTrades(ctx context.Context, stockCode string, limit int) ([]*entity.TradeItem, error)
	UpdateTicker(stockCode string, price int, quantity int)
}

type marketService struct {
	orderBookRepo repositories.OrderBookRepository
	tradeRepo     repositories.TradeRepository
	tickers       map[string]*entity.Ticker
	firstPrices   map[string]int
	mu            sync.RWMutex
}

func NewMarketService(orderBookRepo repositories.OrderBookRepository, tradeRepo repositories.TradeRepository) MarketService {
	return &marketService{
		orderBookRepo: orderBookRepo,
		tradeRepo:     tradeRepo,
		tickers:       make(map[string]*entity.Ticker),
		firstPrices:   make(map[string]int),
	}
}

func (s *marketService) GetTicker(ctx context.Context, stockCode string) (*entity.Ticker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ticker, exists := s.tickers[stockCode]
	if !exists {
		// Return empty ticker if no trades yet
		return &entity.Ticker{
			StockCode: stockCode,
			LastPrice: 0,
			Change:    0,
			Volume:    0,
		}, nil
	}
	return ticker, nil
}

func (s *marketService) GetOrderBookSnapshot(ctx context.Context, stockCode string, depth int) (*entity.OrderBook, error) {
	book, err := s.orderBookRepo.GetBookByStockCode(ctx, stockCode)
	if err != nil {
		return nil, err
	}

	// Clone and truncate to depth
	snapshot := &entity.OrderBook{
		Buys:  s.truncate(book.Buys, depth),
		Sells: s.truncate(book.Sells, depth),
	}

	return snapshot, nil
}

func (s *marketService) GetRecentTrades(ctx context.Context, stockCode string, limit int) ([]*entity.TradeItem, error) {
	allTrades := s.tradeRepo.GetAll()
	var filtered []*entity.TradeItem

	// Filter by stock code and take last N
	// Note: GetAll might return them in chronological order, so we take from the end.
	for i := len(allTrades) - 1; i >= 0; i-- {
		if allTrades[i].StockCode == stockCode {
			filtered = append(filtered, allTrades[i])
		}
		if len(filtered) >= limit {
			break
		}
	}

	return filtered, nil
}

func (s *marketService) UpdateTicker(stockCode string, price int, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.firstPrices[stockCode]; !exists {
		s.firstPrices[stockCode] = price
	}

	ticker, exists := s.tickers[stockCode]
	if !exists {
		ticker = &entity.Ticker{
			StockCode: stockCode,
		}
		s.tickers[stockCode] = ticker
	}

	ticker.LastPrice = price
	ticker.Volume += quantity

	// Calculate change
	firstPrice := s.firstPrices[stockCode]
	if firstPrice > 0 {
		ticker.Change = (float64(price-firstPrice) / float64(firstPrice)) * 100
	}
}

func (s *marketService) truncate(orders []*entity.OrderItem, depth int) []*entity.OrderItem {
	if len(orders) <= depth {
		return orders
	}
	return orders[:depth]
}
