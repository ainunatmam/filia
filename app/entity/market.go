package entity

type Ticker struct {
	StockCode string  `json:"stock_code"`
	LastPrice int     `json:"last_price"`
	Change    float64 `json:"change"`
	Volume    int     `json:"volume"`
}

type MarketSnapshot struct {
	Ticker      *Ticker      `json:"ticker"`
	OrderBook   *OrderBook   `json:"order_book"`
	RecentTrades []*TradeItem `json:"recent_trades"`
}

// WebSocket Message Format
type WSMessage struct {
	Type    string      `json:"type"`    // "subscribe", "unsubscribe", "broadcast"
	Channel string      `json:"channel"` // e.g., "market.ticker:AAPL"
	Data    interface{} `json:"data,omitempty"`
}
