package entity

type OrderStatus string

const (
	OrderStatusOpen    OrderStatus = "OPEN"
	OrderStatusPartial OrderStatus = "PARTIAL"
	OrderStatusFilled  OrderStatus = "FILLED"
)

type OrderBook struct {
	Buys  []*OrderItem `json:"buys"`
	Sells []*OrderItem `json:"sells"`
}

type OrderStore struct {
	orders map[string]*OrderItem
}

type OrderItem struct {
	OrderId        string      `json:"order_id"`
	Status         OrderStatus `json:"status"`
	StockCode      string      `json:"stock_code"`
	Side           string      `json:"side"`
	Price          int         `json:"price"`
	Quantity       int         `json:"quantity"`
	FilledQuantity int         `json:"filled_quantity"`
	CreatedAt      string      `json:"created_at"`
}
