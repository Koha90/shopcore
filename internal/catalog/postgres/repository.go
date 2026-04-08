package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository implements Postrgres-backed catalog write operations.
//
// Read-side catalog loading lives in Loader.
// Write-side operations are introduced gradually through explicit methods
// used by catalog application service.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository constructs catalog Postrgres repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
