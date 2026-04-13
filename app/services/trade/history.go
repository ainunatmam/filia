package trade

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/libraries"
)

func (s *tradeService) History(ctx context.Context) ([]*entity.TradeItem, error) {
	libraries.Logger.Debug().Msg("[TradeService] Fetching trade history")
	result := s.tradeRepo.GetAll()
	libraries.Logger.Info().Int("count", len(result)).Msg("[TradeService] Trade history fetched successfully")
	return result, nil
}