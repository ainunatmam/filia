package order

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/presentation"
	"mini-exchange/app/repositories"
	"mini-exchange/app/services/matching"
)

type OrderService interface {
	Create(ctx context.Context, payload presentation.CreateOrderRequest) error
	List(ctx context.Context) ([]*entity.OrderItem, error)
}

type orderService struct {
	orderRepository repositories.OrderRepository
	matchingService matching.MatchingService
}

func NewOrderService(orderRepository repositories.OrderRepository, matchingService matching.MatchingService) OrderService {
	return &orderService{
		orderRepository: orderRepository,
		matchingService: matchingService,
	}
}
