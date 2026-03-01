package services

import "context"

// TxManager defines abstraction over transaction handling.
//
// It allow application services to execute
// business operations atomically without
// depending on concrete database implementation.
type TxManager interface {
	// WithinTransaction execute fn within a transaction.
	//
	// If fn returns error, transaction is rolled back.
	// if fn returns nil, transaction is commited.
	WithinTransaction(
		ctx context.Context,
		fn func(ctx context.Context) error,
	) error
}
