package example

import (
	"context"
	"mini-exchange/app/entity"
	"mini-exchange/app/presentation"
)

func (s *exampleService) Create(ctx context.Context, req *presentation.ExampleRequest) error {
	err := s.exampleRepo.Create(ctx, entity.Example{
		Name: req.Name,
	})
	if err != nil {
		return err
	}
	return nil
}
