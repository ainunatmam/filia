package repositories

import (
	"context"
	"errors"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
	"sync"
)

type orderRepo struct{
	orderItems map[string]*entity.OrderItem
	mu 	  sync.RWMutex
}

func NewOrderRepository() OrderRepository {
	return &orderRepo{
		orderItems: make(map[string]*entity.OrderItem),
	}
}

func (r *orderRepo) Create(ctx context.Context, order *entity.OrderItem) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	libraries.Logger.Debug().Str("order_id", order.OrderId).Msg("[OrderRepository] Creating order")

	if _, exists := r.orderItems[order.OrderId]; exists {
		libraries.Logger.Warn().Str("order_id", order.OrderId).Msg("[OrderRepository] Order already exists")
		return errors.New("order already exists")
	}
	r.orderItems[order.OrderId] = order
	libraries.Logger.Info().Str("order_id", order.OrderId).Msg("[OrderRepository] Order created successfully")
	return nil
}

func (r *orderRepo) GetAll() []*entity.OrderItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	libraries.Logger.Debug().Int("count", len(r.orderItems)).Msg("[OrderRepository] Getting all orders")

	result := make([]*entity.OrderItem, 0, len(r.orderItems))
	for _, o := range r.orderItems {
		result = append(result, o)
	}

	return result
}