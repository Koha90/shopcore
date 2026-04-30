package postgres

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository provides Postgres-backend payment storage operations.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository construct payment Postgres repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) ensureReady(op string) error {
	if r == nil {
		return fmt.Errorf("%s: repository is nil", op)
	}
	if r.pool == nil {
		return fmt.Errorf("%s: pool is nil", op)
	}

	return nil
}
