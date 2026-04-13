package order

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
	"mini-exchange/app/presentation"
	"mini-exchange/app/utilities"
	"time"
)

func (s *orderService) Create(ctx context.Context, payload presentation.CreateOrderRequest) error {
	orderId := utilities.GenerateUUID()
	now := time.Now()
	order := entity.OrderItem{
		StockCode: payload.StockCode,
		Side: payload.Side,
		Price: payload.Price,
		Quantity: payload.Quantity,
		OrderId: orderId,
		Status: entity.OrderStatusOpen,
		CreatedAt: now.Format(time.RFC3339),
	}
	err := s.orderRepository.Create(ctx, &order)
	if err != nil {
		libraries.Logger.Error().Msg("[OrderService] Failed to create order: " + err.Error())
		return err
	}

	err = s.matchingService.Submit(ctx, &order)
	if err != nil {
		libraries.Logger.Error().Msg("[OrderService] Failed to submit order to matching service: " + err.Error())
		return err
	}
	return nil
}