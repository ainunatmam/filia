package matching

import (
	"context"
	"mini-exchange/app/entity"
)

func (s *matchingService) Submit(ctx context.Context, order *entity.OrderItem) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.orderCh <- order:
		return nil
	}
}