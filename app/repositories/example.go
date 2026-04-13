package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"mini-exchange/app/entity"

	"github.com/doug-martin/goqu/v9"
)

type exampleRepository struct {
	db *goqu.Database
}

var exampleSelectCols = []interface{}{
	"id",
	"name",
	"created_at",
	"updated_at",
	"deleted_at",
}

func NewExampleRepository(db *goqu.Database) ExampleRepository {
	return &exampleRepository{db: db}
}

func (r *exampleRepository) Create(ctx context.Context, data entity.Example) error {
	sql, args, err := r.db.From(tableExample).Insert().Rows(
		data,
	).Prepared(true).ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to create example: %w", err)
	}
	return nil
}

func (r *exampleRepository) Find(ctx context.Context, id uint64) (*entity.Example, error) {
	var example entity.Example
	found, err := r.db.From(tableExample).
		Select(
			exampleSelectCols...,
		).
		Where(goqu.Ex{"id": id, "deleted_at": nil}).
		ScanStruct(&example)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &example, nil
}

func (r *exampleRepository) Get(ctx context.Context) ([]*entity.Example, error) {
	var examples []*entity.Example
	err := r.db.From(tableExample).
		Select(
			exampleSelectCols...,
		).
		Where(goqu.Ex{"deleted_at": nil}).
		ScanStructs(&examples)
	if err != nil {
		return nil, err
	}
	return examples, nil
}

func (r *exampleRepository) Update(ctx context.Context, example *entity.Example) (*entity.Example, error) {
	sql, args, err := r.db.From(tableExample).
		Update().
		Set(example).
		Where(goqu.Ex{"id": example.ID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update example: %w", err)
	}
	return example, nil
}

func (r *exampleRepository) Delete(ctx context.Context, id uint64) error {
	sql, args, err := r.db.From(tableExample).
		Update().
		Set(goqu.Record{"deleted_at": goqu.L("NOW()")}).
		Where(goqu.Ex{"id": id}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete example: %w", err)
	}
	return nil
}

func (r *exampleRepository) UpdateTx(ctx context.Context, tx *sql.Tx, example *entity.Example) (*entity.Example, error) {
	sql, args, err := r.db.From(tableExample).
		Update().
		Set(example).
		Where(goqu.Ex{"id": example.ID}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update example in transaction: %w", err)
	}
	return example, nil
}

func (r *exampleRepository) DeleteTx(ctx context.Context, tx *sql.Tx, id uint64) error {
	sql, args, err := r.db.From(tableExample).
		Update().
		Set(goqu.Record{"deleted_at": goqu.L("NOW()")}).
		Where(goqu.Ex{"id": id}).
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to delete example in transaction: %w", err)
	}
	return nil
}

func (r *exampleRepository) FindForUpdate(ctx context.Context, tx *sql.Tx, id uint64) (*entity.Example, error) {
	var example entity.Example
	sql, args, err := r.db.From(tableExample).
		Select(
			exampleSelectCols...,
		).
		Where(goqu.Ex{"id": id, "deleted_at": nil}).
		ForUpdate(goqu.Wait).Prepared(true).ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build find for update query: %w", err)
	}

	row := tx.QueryRowContext(ctx, sql, args...)
	err = row.Scan(
		&example.ID,
		&example.Name,
		&example.CreatedAt,
		&example.UpdatedAt,
		&example.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &example, nil
}
