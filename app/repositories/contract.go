package repositories

import (
	"context"
	"database/sql"
	"mini-exchange/app/entity"
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

type OrderRepository interface {
	Create(ctx context.Context, order *entity.OrderItem) error
	GetAll() []*entity.OrderItem 
}

type OrderBookRepository interface {
	GetBookByStockCode(ctx context.Context, code string) (*entity.OrderBook, error)
	Save(ctx context.Context, stockCode string, book *entity.OrderBook) error
}

type TradeRepository interface {
	Save(ctx context.Context, trade *entity.TradeItem) error
	GetAll() []*entity.TradeItem
}