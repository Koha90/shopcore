package bootstrap

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/app/runtime/telegram"
	"github.com/koha90/shopcore/internal/manager"
	orderpg "github.com/koha90/shopcore/internal/order/postgres"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

// OrderPoolResolver resolves pgx pool by bot database ID.
type OrderPoolResolver interface {
	Resolve(databaseID string) (*pgxpool.Pool, error)
}

// NewTelegramOrderFactory builds per-bot order creator using database-aware wiring.
//
// Wiring path:
//
//	bot -> spec -> database_id -> pg pool ->
//
// postgres order repository -> order service
func NewTelegramOrderFactory(
	resolver OrderPoolResolver,
) func(spec manager.BotSpec) (telegram.OrderRuntimeService, error) {
	return func(spec manager.BotSpec) (telegram.OrderRuntimeService, error) {
		const op = "new telegram order factory"

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

		repo := orderpg.NewRepository(pool)
		return ordersvc.New(repo), nil
	}
}
