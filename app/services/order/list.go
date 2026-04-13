package order

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
)

func (s *orderService) List(ctx context.Context) ([]*entity.OrderItem, error) {
	libraries.Logger.Debug().Msg("[OrderService] Fetching all orders")
	result := s.orderRepository.GetAll()
	libraries.Logger.Info().Int("count", len(result)).Msg("[OrderService] Orders fetched successfully")
	return result, nil
}