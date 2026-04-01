package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Loader struct {
	pool *pgxpool.Pool
}

func NewLoader(pool *pgxpool.Pool) *Loader {
	return &Loader{
		pool: pool,
	}
}
