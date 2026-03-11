package example

import (
	"context"
	"wallet-api/app/presentation"
	"wallet-api/app/repositories"
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
