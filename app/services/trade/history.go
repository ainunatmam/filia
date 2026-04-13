package trade

import (
	"context"
	"mini-exchange/app/entity"
)

func (s *tradeService)	History(ctx context.Context) ([]*entity.TradeItem, error) {
	result := s.tradeRepo.GetAll()
	return result, nil
}