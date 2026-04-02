package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"

	catalogpg "github.com/koha90/shopcore/internal/catalog/postgres"
	"github.com/koha90/shopcore/internal/flow"
	"github.com/koha90/shopcore/internal/manager"
)

type PoolResolver interface {
	Resolve(databaseID string) (*pgxpool.Pool, error)
}

func NewTelegramFlowFactory(resolver PoolResolver) func(spec manager.BotSpec) *flow.Service {
	return func(spec manager.BotSpec) *flow.Service {
		if resolver == nil {
			return flow.NewService(nil)
		}

		pool, err := resolver.Resolve(spec.DatabaseID)
		if err != nil || pool == nil {
			return flow.NewService(nil)
		}

		loader := catalogpg.NewLoader(pool)
		provider := catalogpg.NewCatalogProvider(loader)

		return flow.NewServiceWithCatalogProvider(nil, provider)
	}
}
