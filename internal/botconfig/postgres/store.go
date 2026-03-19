package postgres

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides PostgreSQL-backend repositories for botconfig.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore creates a new PostgreSQL store.
func NewStore(pool *pgxpool.Pool) *Store {
	if pool == nil {
		panic("botconfig/postgres: pool is nil")
	}

	return &Store{pool: pool}
}

// BotRepository returns PostgreSQL-backend bat repository.
func (s *Store) BotRepository() *BotRepository {
	return &BotRepository{pool: s.pool}
}

// DatabaseProfileRepository returns PostgreSQL-backend database profile repository.
func (s *Store) DatabaseProfileRepository() *DatabaseProfileRepository {
	return &DatabaseProfileRepository{pool: s.pool}
}
