package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"botmanager/internal/app/bootstrap"
	"botmanager/internal/app/pgapp"
	"botmanager/internal/app/seed"
	"botmanager/internal/botconfig"
	botconfigpg "botmanager/internal/botconfig/postgres"
	"botmanager/internal/config"
	"botmanager/internal/manager"
	"botmanager/internal/tui"
	"botmanager/pkg/logger"
)

// demoRunner simulates bot runtime lifecycle for local development.
//
// It is intentionally simple:
//   - broken-bot fails during startup
//   - slow-bot becomes ready after a delay
//   - other bots become ready immediately
//
// This runner is useful while wiring storage, bootstrap, and TUI together
// before real Telegram runtime is connected.
type demoRunner struct{}

// Run starts demo bot runtime and reports readiness through ready callback.
func (r *demoRunner) Run(ctx context.Context, spec manager.BotSpec, ready func()) error {
	switch spec.ID {
	case "broken-bot":
		time.Sleep(700 * time.Millisecond)
		return errors.New("telegram auth failed")

	case "slow-bot":
		time.Sleep(4 * time.Second)
		ready()
		<-ctx.Done()
		return nil

	default:
		ready()
		<-ctx.Done()
		return nil
	}
}

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()

	appLogger, err := logger.Setup(cfg.Env)
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}
	defer func() {
		_ = appLogger.Close()
	}()

	pgCfg := pgapp.LoadConfigFromEnv()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgapp.OpenPool(ctx, pgCfg)
	if err != nil {
		appLogger.Error("failed to open postgres pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	store := botconfigpg.NewStore(pool)

	cfgSvc := botconfig.NewService(
		store.BotRepository(),
		store.DatabaseProfileRepository(),
		nil,
	)

	if err = seed.EnsureDemoData(context.Background(), cfgSvc); err != nil {
		appLogger.Error("failed to ensure demo data", "err", err)
		os.Exit(1)
	}

	mgr := manager.New(&demoRunner{})

	starter := bootstrap.NewStarter(cfgSvc, mgr)

	bootstrapResults, err := starter.StartEnabled(context.Background())
	for _, result := range bootstrapResults {
		if result.Err != nil {
			appLogger.Error(
				"bootstrap bot failed",
				"bot_id", result.ID,
				"registered", result.Registered,
				"started", result.Started,
				"err", result.Err,
			)
			continue
		}

		appLogger.Info(
			"bootstrap bot ok",
			"bot_id", result.ID,
			"registered", result.Registered,
			"started", result.Started,
		)
	}

	if err != nil {
		appLogger.Error("bootstrap finished with errors", "err", err)
	}

	if err := tui.Run(mgr, cfgSvc); err != nil {
		appLogger.Error("tui exited with error", "err", err)
		os.Exit(1)
	}
}
