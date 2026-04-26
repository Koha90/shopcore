package bootstrap

import "github.com/jackc/pgx/v5/pgxpool"

// OrderPoolResolver resolves pgx pool by bot database ID.
type OrderPoolResolver interface {
	Resolve(databaseID string) (*pgxpool.Pool, error)
}

// NewTelegramOrderFactory
