package tuiapp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/app/pgapp"
	"github.com/koha90/shopcore/internal/botconfig"
	botconfigpg "github.com/koha90/shopcore/internal/botconfig/postgres"
	"github.com/koha90/shopcore/internal/manager"
)

// openPool initialized PostgreSQL connection pool for the TUI application.
func openPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	timeout := cfg.OpenDBTimeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	openCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pool, err := pgapp.OpenPool(openCtx, cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("open postgres pool: %w", err)
	}

	return pool, nil
}

// buildApp assembles storage, configuration service and runtime manager.
//
// This function contains dependency wiring only.
func buildApp(
	pool *pgxpool.Pool,
	runner manager.Runner,
	log *slog.Logger,
) (*App, error) {
	store := botconfigpg.NewStore(pool)

	cfgSvc := botconfig.NewService(
		store.BotRepository(),
		store.DatabaseProfileRepository(),
		log,
	)

	mgr := manager.New(runner)

	return &App{
		Pool:      pool,
		BotConfig: cfgSvc,
		Manager:   mgr,
	}, nil
}
