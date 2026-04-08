package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	catalogpg "github.com/koha90/shopcore/internal/catalog/postgres"
	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
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
//	bot spec -> database_id -> pg pool ->
//	postgres catalog provider -> flow service
func NewTelegramFlowFactory(resolver PoolResolver) func(spec manager.BotSpec) (*flow.Service, error) {
	return func(spec manager.BotSpec) (*flow.Service, error) {
		const op = "new telegram flow factory"

		if resolver == nil {
			return nil, fmt.Errorf("%s: pool resolver is nil", op)
		}

		pool, err := resolver.Resolve(spec.DatabaseID)
		if err != nil {
			return nil, fmt.Errorf("%s: resolve database %q: %w", op, spec.DatabaseID, err)
		}
		if pool == nil {
			return nil, fmt.Errorf("%s: resolved pool is nil for database %q", op, spec.DatabaseID)
		}

		loader := catalogpg.NewLoader(pool)
		provider := catalogpg.NewCatalogProvider(loader)

		repo := catalogpg.NewRepository(pool)
		catalog := catalogservice.New(repo, repo)
		creator := newFlowCategoryCreator(catalog)

		return flow.NewServiceWithDeps(nil, provider, creator), nil
	}
}
