package repositories

import (
	"context"
	"database/sql"
	"mini-exchange/app/entity"
)

// MockExampleRepository is a mock implementation of ExampleRepository for testing
type MockExampleRepository struct {
	CreateFunc        func(ctx context.Context, data entity.Example) error
	GetFunc           func(ctx context.Context) ([]*entity.Example, error)
	FindFunc          func(ctx context.Context, id uint64) (*entity.Example, error)
	UpdateFunc        func(ctx context.Context, example *entity.Example) (*entity.Example, error)
	DeleteFunc        func(ctx context.Context, id uint64) error
	UpdateTxFunc      func(ctx context.Context, tx *sql.Tx, example *entity.Example) (*entity.Example, error)
	DeleteTxFunc      func(ctx context.Context, tx *sql.Tx, id uint64) error
	FindForUpdateFunc func(ctx context.Context, tx *sql.Tx, id uint64) (*entity.Example, error)
}

func (m *MockExampleRepository) Create(ctx context.Context, data entity.Example) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, data)
	}
	return nil
}

func (m *MockExampleRepository) Get(ctx context.Context) ([]*entity.Example, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx)
	}
	return nil, nil
}

func (m *MockExampleRepository) Find(ctx context.Context, id uint64) (*entity.Example, error) {
	if m.FindFunc != nil {
		return m.FindFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockExampleRepository) Update(ctx context.Context, example *entity.Example) (*entity.Example, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, example)
	}
	return example, nil
}

func (m *MockExampleRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockExampleRepository) UpdateTx(ctx context.Context, tx *sql.Tx, example *entity.Example) (*entity.Example, error) {
	if m.UpdateTxFunc != nil {
		return m.UpdateTxFunc(ctx, tx, example)
	}
	return example, nil
}

func (m *MockExampleRepository) DeleteTx(ctx context.Context, tx *sql.Tx, id uint64) error {
	if m.DeleteTxFunc != nil {
		return m.DeleteTxFunc(ctx, tx, id)
	}
	return nil
}

func (m *MockExampleRepository) FindForUpdate(ctx context.Context, tx *sql.Tx, id uint64) (*entity.Example, error) {
	if m.FindForUpdateFunc != nil {
		return m.FindForUpdateFunc(ctx, tx, id)
	}
	return nil, nil
}
