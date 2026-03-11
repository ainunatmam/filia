package libraries

import (
	"context"
	"database/sql"
)

type TransactionManager interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CommitTx(tx *sql.Tx) error
	RollbackTx(tx *sql.Tx) error
	CommitOrRollback(err error, tx *sql.Tx) error
}

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return t.db.BeginTx(ctx, nil)
}

func (t *transactionManager) CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func (t *transactionManager) RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}

func (t *transactionManager) CommitOrRollback(err error, tx *sql.Tx) error {
	if err != nil {
		tx.Rollback()
	} else {
		err = tx.Commit()
	}
	return err
}