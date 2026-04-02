package tuiapp

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/app/bootstrap"
	"github.com/koha90/shopcore/internal/app/pgapp"
	"github.com/koha90/shopcore/internal/app/runtime/telegram"
	botconfigpg "github.com/koha90/shopcore/internal/botconfig/postgres"
	"github.com/koha90/shopcore/internal/manager"
)

// buildRunner assembles Telegram runtim runner with database-aware flow factory.
//
// Wiring path:
//
//	app pool -> database profile repository -> pool registry ->
//	telegram flow factory -> telegram runner
func buildRunner(
	ctx context.Context,
	pool *pgxpool.Pool,
	tgCfg telegram.Config,
	runtimeLog *slog.Logger,
) (manager.Runner, error) {
	const op = "build runner"

	if pool == nil {
		return nil, fmt.Errorf("%s: pool is nil", op)
	}
	if runtimeLog == nil {
		return nil, fmt.Errorf("%s: runtime logger is nil", op)
	}

	store := botconfigpg.NewStore(pool)

	profilesRepo := store.DatabaseProfileRepository()
	poolRegistry := pgapp.NewPoolRegistry(ctx, profilesRepo)

	flowFactory := bootstrap.NewTelegramFlowFactory(poolRegistry)

	runner := telegram.NewRunnerWithFlowFactory(
		tgCfg,
		runtimeLog,
		flowFactory,
	)

	return runner, nil
}
