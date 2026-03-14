package postgres

import (
	"context"
	"database/sql"
)

// TxManager implements transactional execution for PostgreSQL storage.
//
// It starts SQL transaction, injects it into context and commit it
// on successful completion. If fn returns an error, transaction is rolled back.
type TxManager struct {
	db *sql.DB
}

// NewTxManager creates a new PostgreSQL transaction manager.
func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

// WithinTransaction executes fn inside SQL transaction.
func (m *TxManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)

	if err := fn(ctx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

type txKey struct{}

func txFromContext(ctx context.Context) (*sql.Tx, bool) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	return tx, ok
}
