package repositories

import (
	"context"
	"database/sql"
	"wallet-api/app/entity"
)

const (
	tableExample = "examples"
)

type ExampleRepository interface {
	Create(ctx context.Context, data entity.Example) error
	Get(ctx context.Context) ([]*entity.Example, error)
	Find(ctx context.Context, id uint64) (*entity.Example, error)
	Update(ctx context.Context, example *entity.Example) (*entity.Example, error)
	Delete(ctx context.Context, id uint64) error
	UpdateTx(ctx context.Context, tx *sql.Tx, example *entity.Example) (*entity.Example, error)
	DeleteTx(ctx context.Context, tx *sql.Tx, id uint64) error
	FindForUpdate(ctx context.Context, tx *sql.Tx, id uint64) (*entity.Example, error)
}
