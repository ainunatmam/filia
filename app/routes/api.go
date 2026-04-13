package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func (r *router) api() fiber.Router {
	root := r.app.Group("/")
	api := root.Group("/api")

	r.startBackgroundTasks()
	r.registerWebSocketRoutes()
	r.registerGreetingRoutes(api)
	r.registerMarketRoutes(api)
	r.registerOrderRoutes(api)
	r.registerTradeRoutes(api)

	return root
}

func (r *router) registerWebSocketRoutes() {
	r.app.Get("/ws", websocket.New(r.container.WSHandler))
}

func (r *router) registerGreetingRoutes(api fiber.Router) {
	api.Get("/greeting", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Konnichiwa !",
		})
	})
}

func (r *router) registerMarketRoutes(api fiber.Router) {
	market := api.Group("/market")
	market.Get("/snapshot", r.container.MarketHandler.GetSnapshot)
}

func (r *router) registerOrderRoutes(api fiber.Router) {
	order := api.Group("/order")
	order.Get("", r.container.OrderHandler.List)
	order.Post("", r.container.OrderHandler.Create)
}

func (r *router) registerTradeRoutes(api fiber.Router) {
	trade := api.Group("/trade")
	trade.Get("", r.container.TradeHandler.History)
}
