package pgapp

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// OpenPool creates PostgreSQL pgx pool from app config and verifies connectivity with Ping.
func OpenPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	return OpenPoolDSN(ctx, cfg.DSN())
}

// OpenPoolDSN creates PostgreSQL pgx pool from raw DSN and verifies connectivity with Ping.
func OpenPoolDSN(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, fmt.Errorf("open pgx pool: empty dsn")
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}

	return pool, nil
}
