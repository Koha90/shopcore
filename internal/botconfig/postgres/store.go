package postgres

import "database/sql"

// Store provides PostgreSQL-backend repositories for botconfig.
type Store struct {
	db *sql.DB
}

// NewStore creates a new PostgreSQL store.
func NewStore(db *sql.DB) *Store {
	if db == nil {
		panic("botconfig/postgres: db is nil")
	}

	return &Store{db: db}
}

// BotRepository returns PostgreSQL-backend bat repository.
func (s *Store) BotRepository() *BotRepository {
	return &BotRepository{db: s.db}
}

// DatabaseProfileRepository returns PostgreSQL-backend database profile repository.
func (s *Store) DatabaseProfileRepository() *DatabaseProfileRepository {
	return &DatabaseProfileRepository{db: s.db}
}
