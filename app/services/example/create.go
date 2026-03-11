package example

import (
	"context"
	"wallet-api/app/entity"
	"wallet-api/app/presentation"
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
