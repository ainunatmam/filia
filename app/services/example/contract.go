package example

import (
	"context"
	"mini-exchange/app/presentation"
	"mini-exchange/app/repositories"
)

type ExampleService interface {
	Create(ctx context.Context, req *presentation.ExampleRequest) error
}

type exampleService struct {
	exampleRepo repositories.ExampleRepository
}

func NewExampleService(exampleRepo repositories.ExampleRepository) ExampleService {
	return &exampleService{
		exampleRepo: exampleRepo,
	}
}
