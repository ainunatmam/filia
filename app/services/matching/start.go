package matching

import (
	"context"
	"mini-exchange/app/entity"
)

func (s *matchingService) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case order := <-s.orderCh:
				s.process(ctx, order)
			}
		}
	}()
}

func (s *matchingService) process(ctx context.Context, order *entity.OrderItem) {

	book, err := s.orderBookRepo.GetBookByStockCode(ctx, order.StockCode)
	if err != nil {
		return
	}

	if order.Side == "BUY" {
		s.matchBuy(ctx, order, book)
	} else {
		s.matchSell(ctx, order, book)
	}

	_ = s.orderBookRepo.Save(ctx, order.StockCode, book)
	s.wsHub.Broadcast("market.orderbook:"+order.StockCode, book)
}

func (s *matchingService) matchBuy(ctx context.Context, order *entity.OrderItem, book *entity.OrderBook) {

	for i := 0; i < len(book.Sells); i++ {
		sell := book.Sells[i]

		if sell.Price > order.Price {
			continue
		}

		matchQty := min(
			order.Quantity-order.FilledQuantity,
			sell.Quantity-sell.FilledQuantity,
		)

		order.FilledQuantity += matchQty
		sell.FilledQuantity += matchQty

		updateStatus(order)
		updateStatus(sell)

		trade := &entity.TradeItem{
			StockCode:   order.StockCode,
			Price:       sell.Price,
			Quantity:    matchQty,
			BuyOrderId:  order.OrderId,
			SellOrderId: sell.OrderId,
		}

		s.tradeRepo.Save(ctx, trade)
		s.marketService.UpdateTicker(order.StockCode, trade.Price, trade.Quantity)

		s.wsHub.Broadcast("market.trade:"+order.StockCode, trade)
		ticker, _ := s.marketService.GetTicker(ctx, order.StockCode)
		s.wsHub.Broadcast("market.ticker:"+order.StockCode, ticker)
		s.wsHub.Broadcast("order.update", order)
		s.wsHub.Broadcast("order.update", sell)

		if sell.FilledQuantity == sell.Quantity {
			book.Sells = removeIndex(book.Sells, i)
			i--
		}

		if order.FilledQuantity == order.Quantity {
			return
		}
	}

	book.Buys = append(book.Buys, order)
}

func (s *matchingService) matchSell(ctx context.Context, order *entity.OrderItem, book *entity.OrderBook) {

	for i := 0; i < len(book.Buys); i++ {
		buy := book.Buys[i]

		if buy.Price < order.Price {
			continue
		}

		matchQty := min(
			order.Quantity-order.FilledQuantity,
			buy.Quantity-buy.FilledQuantity,
		)

		order.FilledQuantity += matchQty
		buy.FilledQuantity += matchQty

		updateStatus(order)
		updateStatus(buy)

		trade := &entity.TradeItem{
			StockCode:   order.StockCode,
			Price:       buy.Price,
			Quantity:    matchQty,
			BuyOrderId:  buy.OrderId,
			SellOrderId: order.OrderId,
		}

		s.tradeRepo.Save(ctx, trade)
		s.marketService.UpdateTicker(order.StockCode, trade.Price, trade.Quantity)

		s.wsHub.Broadcast("market.trade:"+order.StockCode, trade)
		ticker, _ := s.marketService.GetTicker(ctx, order.StockCode)
		s.wsHub.Broadcast("market.ticker:"+order.StockCode, ticker)
		s.wsHub.Broadcast("order.update", order)
		s.wsHub.Broadcast("order.update", buy)

		if buy.FilledQuantity == buy.Quantity {
			book.Buys = removeIndex(book.Buys, i)
			i--
		}

		if order.FilledQuantity == order.Quantity {
			return
		}
	}

	book.Sells = append(book.Sells, order)
}

func updateStatus(o *entity.OrderItem) {
	if o.FilledQuantity == 0 {
		o.Status = entity.OrderStatusOpen
	} else if o.FilledQuantity < o.Quantity {
		o.Status = entity.OrderStatusPartial
	} else {
		o.Status = entity.OrderStatusFilled
	}
}

func removeIndex(slice []*entity.OrderItem, index int) []*entity.OrderItem {
	return append(slice[:index], slice[index+1:]...)
}