package order

import (
	"context"
	"mini-exchange/app/entity"
)

func (s *orderService) List(ctx context.Context) ([]*entity.OrderItem, error) {
	result := s.orderRepository.GetAll()
	return result, nil
}