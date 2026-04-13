package presentation

type ExampleRequest struct {
	Name string `validate:"required" json:"name"`
}

type CreateOrderRequest struct {
	StockCode string `json:"stock_code" validate:"required"`
	Side      string `json:"side" validate:"required,oneof=BUY SELL"`
	Price     int    `json:"price" validate:"required,gt=0"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type ListOrderRequest struct {
	Stock  string `query:"stock" validate:"required"`
	Status string `query:"status" validate:"omitempty,oneof=open closed"`
}