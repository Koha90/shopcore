// Package memory provides in-memory storage adapters.
package memory

import (
	"context"
	"sync"
)

// TxManager implements transactional execution for in-memory storage.
//
// It simulates a transaction by locking shared storage state
// during execution of a function.
type TxManager struct {
	mu *sync.Mutex
}

// NewTxManager creates a new in-memory transaction manager.
//
// mu must point to shared storage mutex used by all repositoies
// participating in the same logical transaction.
func NewTxManager(mu *sync.Mutex) *TxManager {
	return &TxManager{mu: mu}
}

// WithinTransaction executes fn inside a critical section.
func (t *TxManager) WithinTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return fn(ctx)
}
