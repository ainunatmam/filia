package routes

import (
	"context"
	"mini-exchange/app/handlers/market"
	"mini-exchange/app/handlers/order"
	"mini-exchange/app/handlers/trade"
	"mini-exchange/app/repositories"
	marketService "mini-exchange/app/services/market"
	"mini-exchange/app/services/matching"
	orderService "mini-exchange/app/services/order"
	tradeService "mini-exchange/app/services/trade"
	"mini-exchange/app/tasks"
	"mini-exchange/app/ws"
	"mini-exchange/config"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type routes struct {
	Api fiber.Router
}

type container struct {
	Hub           ws.Hub
	OrderRepo     repositories.OrderRepository
	OrderBookRepo repositories.OrderBookRepository
	TradeRepo     repositories.TradeRepository

	OrderService    orderService.OrderService
	TradeService    tradeService.TradeService
	MarketService   marketService.MarketService
	MatchingService matching.MatchingService

	OrderHandler  order.OrderHandler
	TradeHandler  trade.TradeHandler
	MarketHandler market.MarketHandler
	WSHandler     func(*websocket.Conn)

	MatchingTask   tasks.MatchingTask
	SimulationTask tasks.SimulationTask
}

type router struct {
	app       *fiber.App
	cfg       *config.Config
	startOnce sync.Once
	container *container
}

func NewRouter(app *fiber.App, cfg *config.Config) *router {
	return &router{
		app: app,
		cfg: cfg,
	}
}

func (r *router) initContainer() {
	c := &container{}

	c.Hub = ws.NewHub()
	c.OrderRepo = repositories.NewOrderRepository()
	c.OrderBookRepo = repositories.NewOrderBookRepository()
	c.TradeRepo = repositories.NewTradeRepository()

	c.TradeService = tradeService.NewTradeService(c.TradeRepo)
	c.MarketService = marketService.NewMarketService(c.OrderBookRepo, c.TradeRepo)
	c.MatchingService = matching.NewMatchingService(c.OrderBookRepo, c.TradeRepo, c.Hub, c.MarketService)
	c.OrderService = orderService.NewOrderService(c.OrderRepo, c.MatchingService)

	c.OrderHandler = order.NewOrderHandler(c.OrderService)
	c.TradeHandler = trade.NewTradeHandler(c.TradeService)
	c.MarketHandler = market.NewMarketHandler(c.MarketService)
	c.WSHandler = ws.NewWSHandler(c.Hub)

	c.MatchingTask = tasks.NewMatchingTask(c.MatchingService)
	c.SimulationTask = tasks.NewSimulationTask(c.MatchingService, c.MarketService, r.cfg)

	r.container = c
}

func (r *router) startBackgroundTasks() {
	r.startOnce.Do(func() {
		bgContext := context.Background()
		go r.container.Hub.Run()
		go r.container.MatchingTask.Begin(bgContext)
		go r.container.SimulationTask.Run(bgContext)
	})
}

func (r *router) Init() routes {
	r.initContainer()
	return routes{
		Api: r.api(),
	}
}
