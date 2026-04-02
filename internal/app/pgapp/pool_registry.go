package pgapp

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/botconfig"
)

// DatabaseProfileGetter loads database profile by ID.
type DatabaseProfileGetter interface {
	ByID(ctx context.Context, id string) (*botconfig.DatabaseProfile, error)
}

// PoolRegistry resolves and caches pgx pools by database profile ID.
//
// Multiple bots may use the same database_id, so the registry reuses one pool
// per database profile instead of opening a new pool for every bot runtime.
type PoolRegistry struct {
	ctx   context.Context
	repo  DatabaseProfileGetter
	mu    sync.Mutex
	pools map[string]*pgxpool.Pool
}

// NewPoolRegistry constructs database profile pool registry.
func NewPoolRegistry(ctx context.Context, repo DatabaseProfileGetter) *PoolRegistry {
	if ctx == nil {
		ctx = context.Background()
	}

	return &PoolRegistry{
		ctx:   ctx,
		repo:  repo,
		pools: make(map[string]*pgxpool.Pool),
	}
}

// Resolve returns cached or newly opened pgx pool for database profile ID.
func (r *PoolRegistry) Resolve(databaseID string) (*pgxpool.Pool, error) {
	if r.repo == nil {
		return nil, fmt.Errorf("resolve pool: repository is nil")
	}
	if strings.TrimSpace(databaseID) == "" {
		return nil, fmt.Errorf("resolve pool: empty database id")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if pool, ok := r.pools[databaseID]; ok {
		return pool, nil
	}

	profile, err := r.repo.ByID(r.ctx, databaseID)
	if err != nil {
		return nil, fmt.Errorf("resolve pool: load database profile %q: %w", databaseID, err)
	}
	if profile == nil {
		return nil, fmt.Errorf("resolve pool: database profile %q is nil", databaseID)
	}
	if !profile.IsEnabled {
		return nil, fmt.Errorf("resolve pool: database profile %q is disabled", databaseID)
	}
	if strings.TrimSpace(profile.Driver) != "" && profile.Driver != "postgres" {
		return nil, fmt.Errorf("resolve pool: unsupported driver %q for profile %q", profile.Driver, databaseID)
	}

	pool, err := OpenPoolDSN(r.ctx, profile.DSN)
	if err != nil {
		return nil, fmt.Errorf("resolve pool: open pool for profile %q: %w", databaseID, err)
	}

	r.pools[databaseID] = pool
	return pool, nil
}

// Close closes all cached pools.
func (r *PoolRegistry) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, pool := range r.pools {
		if pool != nil {
			pool.Close()
		}
		delete(r.pools, id)
	}
}
