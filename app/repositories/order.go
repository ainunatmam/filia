package repositories

import (
	"context"
	"errors"
	"mini-exchange/app/entity"
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

	if _, exists := r.orderItems[order.OrderId]; exists {
		return errors.New("order already exists")
	}
	r.orderItems[order.OrderId] = order
	return nil
}

func (r *orderRepo) GetAll() []*entity.OrderItem {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*entity.OrderItem, 0, len(r.orderItems))
	for _, o := range r.orderItems {
		result = append(result, o)
	}

	return result
}