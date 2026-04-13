package market

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
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

	libraries.Logger.Debug().Str("stock_code", stockCode).Msg("[MarketService] Getting ticker")

	ticker, exists := s.tickers[stockCode]
	if !exists {
		libraries.Logger.Debug().Str("stock_code", stockCode).Msg("[MarketService] No ticker found, returning empty")
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
	libraries.Logger.Debug().Str("stock_code", stockCode).Int("depth", depth).Msg("[MarketService] Getting order book snapshot")

	book, err := s.orderBookRepo.GetBookByStockCode(ctx, stockCode)
	if err != nil {
		libraries.Logger.Error().Str("stock_code", stockCode).Err(err).Msg("[MarketService] Failed to get order book")
		return nil, err
	}

	// Clone and truncate to depth
	snapshot := &entity.OrderBook{
		Buys:  s.truncate(book.Buys, depth),
		Sells: s.truncate(book.Sells, depth),
	}

	libraries.Logger.Debug().Str("stock_code", stockCode).Int("buys", len(snapshot.Buys)).Int("sells", len(snapshot.Sells)).Msg("[MarketService] Order book snapshot retrieved")
	return snapshot, nil
}

func (s *marketService) GetRecentTrades(ctx context.Context, stockCode string, limit int) ([]*entity.TradeItem, error) {
	libraries.Logger.Debug().Str("stock_code", stockCode).Int("limit", limit).Msg("[MarketService] Getting recent trades")

	allTrades := s.tradeRepo.GetAll()
	var filtered []*entity.TradeItem

	// Filter by stock code and take last N
	for i := len(allTrades) - 1; i >= 0; i-- {
		if allTrades[i].StockCode == stockCode {
			filtered = append(filtered, allTrades[i])
		}
		if len(filtered) >= limit {
			break
		}
	}

	libraries.Logger.Debug().Str("stock_code", stockCode).Int("count", len(filtered)).Msg("[MarketService] Recent trades retrieved")
	return filtered, nil
}

func (s *marketService) UpdateTicker(stockCode string, price int, quantity int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	libraries.Logger.Debug().Str("stock_code", stockCode).Int("price", price).Int("quantity", quantity).Msg("[MarketService] Updating ticker")

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
