package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"botmanager/internal/app/bootstrap"
	"botmanager/internal/app/seed"
	"botmanager/internal/app/tuiapp"
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
	defer func() { _ = appLogger.Close() }()

	appCfg := tuiapp.LoadConfigFromEnv()

	app, err := tuiapp.New(context.Background(), appCfg, &demoRunner{}, appLogger.Logger)
	if err != nil {
		appLogger.Error("build tui app", "err", err)
		os.Exit(1)
	}
	defer app.Close()

	if err := seed.EnsureDemoData(context.Background(), app.BotConfig); err != nil {
		appLogger.Error("failed to ensure demo", "err", err)
		os.Exit(1)
	}

	starter := bootstrap.NewStarter(app.BotConfig, app.Manager)

	results, err := starter.StartEnabled(context.Background())
	for _, result := range results {
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

	if err = tui.Run(app.Manager, app.BotConfig); err != nil {
		appLogger.Error("tui exited with error", "err", err)
		os.Exit(1)
	}
}
