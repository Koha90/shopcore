package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	catalogpg "github.com/koha90/shopcore/internal/catalog/postgres"
	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

// PoolResolver resolves pgx pool by bot database ID.
type PoolResolver interface {
	Resolve(databaseID string) (*pgxpool.Pool, error)
}

// NewTelegramFlowFactory builds per-bot flow service using database-aware catalog provider.
//
// Wiring path:
//
//	bot spec -> database_id -> pg pool -> postgres catalog provider -> flow service
func NewTelegramFlowFactory(resolver PoolResolver) func(spec manager.BotSpec) (*flow.Service, error) {
	return func(spec manager.BotSpec) (*flow.Service, error) {
		const op = "build flow service"

		if resolver == nil {
			return nil, fmt.Errorf("%s: pool resolver is nil", op)
		}

		pool, err := resolver.Resolve(spec.DatabaseID)
		if err != nil || pool == nil {
			return nil, fmt.Errorf("%s: resolve fatabase %q: %w", op, spec.DatabaseID, err)
		}

		loader := catalogpg.NewLoader(pool)
		provider := catalogpg.NewCatalogProvider(loader)

		return flow.NewServiceWithCatalogProvider(nil, provider), nil
	}
}
