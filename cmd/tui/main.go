package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"botmanager/internal/app/bootstrap"
	"botmanager/internal/app/runtime/demo"
	"botmanager/internal/app/seed"
	"botmanager/internal/app/tuiapp"
	"botmanager/internal/config"
	"botmanager/internal/tui"
	"botmanager/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()

	appLogger, err := logger.Setup(cfg.Env)
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}
	defer func() { _ = appLogger.Close() }()

	appCfg := tuiapp.LoadConfigFromEnv()

	app, err := tuiapp.New(context.Background(), appCfg, demo.NewRunner(), appLogger.Logger)
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
