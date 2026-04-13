package routes

import (
	"context"
	"mini-exchange/app/handlers/example"
	"mini-exchange/app/handlers/market"
	"mini-exchange/app/handlers/order"
	"mini-exchange/app/handlers/trade"
	"mini-exchange/app/repositories"
	"mini-exchange/app/tasks"
	"mini-exchange/app/ws"

	exampleService "mini-exchange/app/services/example"
	marketService "mini-exchange/app/services/market"
	"mini-exchange/app/services/matching"
	orderService "mini-exchange/app/services/order"
	tradeService "mini-exchange/app/services/trade"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (r *router) api() fiber.Router {

	root := r.app.Group("/")
	api := root.Group("/api")
	apiMarket := api.Group("/market")
	apiOrder := api.Group("/order")
	apiTrade := api.Group("/trade")

	// init hub
	hub := ws.NewHub()

	// init repo
	exampleRepo := repositories.NewExampleRepository(r.goquLibrary.DB())
	orderRepo := repositories.NewOrderRepository()
	orderBookRepo := repositories.NewOrderBookRepository()
	tradeRepo := repositories.NewTradeRepository()

	// init service
	exampleService := exampleService.NewExampleService(exampleRepo)
	tradeService := tradeService.NewTradeService(tradeRepo)
	marketService := marketService.NewMarketService(orderBookRepo, tradeRepo)
	matchingService := matching.NewMatchingService(orderBookRepo, tradeRepo, hub, marketService)
	orderService := orderService.NewOrderService(orderRepo, matchingService)

	// init handler
	exampleHandler := example.NewExampleHandler(exampleService)
	orderHandler := order.NewOrderHandler(orderService)
	tradeHandler := trade.NewTradeHandler(tradeService)
	marketHandler := market.NewMarketHandler(marketService)
	wsHandler := ws.NewWSHandler(hub)

	// init task 
	matchingTask := tasks.NewMatchingTask(matchingService)
	simulationTask := tasks.NewSimulationTask(matchingService, marketService, r.cfg)

	r.startOnce.Do(func() {
		bgContext := context.Background()
		go hub.Run()
		go matchingTask.Begin(bgContext)
		go simulationTask.Run(bgContext)
	})

	// Websocket
	r.app.Get("/ws", websocket.New(wsHandler))

	// Market
	apiMarket.Get("/snapshot", marketHandler.GetSnapshot)

	// Greeting
	api.Get("/greeting", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Konnichiwa !",
		})
	})

	apiOrder.Get("", orderHandler.List)
	apiOrder.Post("", orderHandler.Create)
	apiTrade.Get("", tradeHandler.History)

	// Example
	api.Post("/example", exampleHandler.Create)
	return root
}

