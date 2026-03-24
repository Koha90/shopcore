package tuiapp

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/botconfig"
	"github.com/koha90/shopcore/internal/manager"
)

// App is the fully wired dependency container for cmd/tui.
//
// It exists only at the application boundry and must not be passed into
// domain logic as a service locator.
type App struct {
	// Pool is the shared PostgreSQL connection pool.
	Pool *pgxpool.Pool

	// BotConfig provides configuration use cases for bots and database profiles.
	BotConfig *botconfig.Service

	// Manager controls registered bot runtimes.
	Manager *manager.Manager
}

// Close releases resources owned by the TUI application container.
func (a *App) Close() error {
	if a == nil || a.Pool == nil {
		return nil
	}

	a.Pool.Close()
	return nil
}

// New constructs a ready-to-use TUI application container.
//
// It opens PostgreSQL, builds repositories, assembles configuration services,
// creates runtime manager, and returns the resulting App.
func New(
	ctx context.Context,
	cfg Config,
	runner manager.Runner,
	log *slog.Logger,
) (*App, error) {
	pool, err := openPool(ctx, cfg)
	if err != nil {
		return nil, err
	}

	app, err := buildApp(pool, runner, log)
	if err != nil {
		pool.Close()
		return nil, err
	}

	return app, nil
}
