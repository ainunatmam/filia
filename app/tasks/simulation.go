package tasks

import (
	"context"
	"math/rand"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
	"mini-exchange/app/services/market"
	"mini-exchange/app/services/matching"
	"mini-exchange/app/utilities"
	"mini-exchange/config"
	"time"
)

type SimulationTask interface {
	Run(ctx context.Context)
}

type simulationTask struct {
	matchingService matching.MatchingService
	marketService   market.MarketService
	cfg             *config.Config
}

func NewSimulationTask(ms matching.MatchingService, mkt market.MarketService, cfg *config.Config) SimulationTask {
	return &simulationTask{
		matchingService: ms,
		marketService:   mkt,
		cfg:             cfg,
	}
}

func (t *simulationTask) Run(ctx context.Context) {
	if !t.cfg.Simulation.Enabled {
		libraries.Logger.Info().Msg("Market Simulation is disabled")
		return
	}

	libraries.Logger.Info().Dur("interval", t.cfg.Simulation.Interval).Msg("Starting Market Simulation")

	ticker := time.NewTicker(t.cfg.Simulation.Interval)
	defer ticker.Stop()

	stocks := []string{"IDR", "BTC", "ETH"}
	defaultPrices := map[string]int{
		"IDR": 15000,
		"BTC": 65000,
		"ETH": 3500,
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.generateRandomOrder(ctx, stocks, defaultPrices)
		}
	}
}

func (t *simulationTask) generateRandomOrder(ctx context.Context, stocks []string, defaultPrices map[string]int) {
	stock := stocks[rand.Intn(len(stocks))]

	tickerData, _ := t.marketService.GetTicker(ctx, stock)
	lastPrice := tickerData.LastPrice
	if lastPrice == 0 {
		lastPrice = defaultPrices[stock]
	}

	side := "BUY"
	if rand.Intn(2) == 0 {
		side = "SELL"
	}

	// Price deviation: -2% to +2%
	deviation := (rand.Float64() * 0.04) - 0.02
	price := int(float64(lastPrice) * (1 + deviation))

	// Occasionally favor matching
	if rand.Intn(10) > 7 {
		if side == "BUY" {
			price = lastPrice + 1
		} else {
			price = lastPrice - 1
		}
	}

	if price <= 0 {
		price = 1
	}

	qty := rand.Intn(100) + 1

	order := &entity.OrderItem{
		OrderId:   utilities.GenerateUUID(),
		StockCode: stock,
		Side:      side,
		Price:      price,
		Quantity:   qty,
		Status:     entity.OrderStatusOpen,
		CreatedAt:  time.Now().Format(time.RFC3339),
	}

	libraries.Logger.Debug().
		Str("side", side).
		Str("stock", stock).
		Int("price", price).
		Int("qty", qty).
		Msg("Simulation: Injecting random order")

	_ = t.matchingService.Submit(ctx, order)
}
