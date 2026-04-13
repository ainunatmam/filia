package entity

import "time"

type TradeItem struct {
	TradeId     string    `json:"trade_id"`
	StockCode   string    `json:"stock_code"`
	Price       int       `json:"price"`
	Quantity    int       `json:"quantity"`
	BuyOrderId  string    `json:"buy_order_id"`
	SellOrderId string    `json:"sell_order_id"`
	Timestamp   time.Time `json:"timestamp"`
}