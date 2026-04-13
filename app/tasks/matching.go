package tasks

import (
	"context"
	"mini-exchange/app/services/matching"
)

type MatchingTask interface {
	Begin(ctx context.Context)
}

type matchingTask struct {
	matchingService matching.MatchingService
}

func NewMatchingTask(matchingService matching.MatchingService) MatchingTask {
	return &matchingTask{
		matchingService: matchingService,
	}
}


func (t *matchingTask) Begin(ctx context.Context) {
	go t.matchingService.Start(ctx)
}