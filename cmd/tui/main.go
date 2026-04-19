package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"

	"github.com/koha90/shopcore/internal/app/bootstrap"
	"github.com/koha90/shopcore/internal/app/runtime/runtimelog"
	"github.com/koha90/shopcore/internal/app/runtime/telegram"
	"github.com/koha90/shopcore/internal/app/seed"
	"github.com/koha90/shopcore/internal/app/tuiapp"
	"github.com/koha90/shopcore/internal/config"
	"github.com/koha90/shopcore/internal/tui"
	"github.com/koha90/shopcore/pkg/logger"
)

func main() {
	_ = godotenv.Load()

	cfg := config.MustLoad()

	appLogger, err := logger.SetupForTUI(cfg.Env)
	if err != nil {
		log.Fatalf("setup logger: %v", err)
	}
	defer func() { _ = appLogger.Close() }()

	appCfg, err := tuiapp.LoadConfigFromEnv()
	if err != nil {
		appLogger.Error("load tui config", "err", err)
		os.Exit(1)
	}

	runtimeLogs := runtimelog.NewStore(300)

	runtimeWrap := func(next slog.Handler) slog.Handler {
		return runtimelog.NewHandler(next, runtimeLogs)
	}

	runtimeLogger, err := logger.NewFileLoggerWithHandler(
		"logs/runtime.log",
		slog.LevelDebug,
		runtimeWrap,
	)
	if err != nil {
		appLogger.Error("setup runtime logger", "err", err)
		os.Exit(1)
	}

	tgCfg := telegram.LoadConfigFromEnv()

	app, err := tuiapp.New(context.Background(), appCfg, tgCfg, runtimeLogger.Logger, appLogger.Logger)
	if err != nil {
		appLogger.Error("build tui app", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := app.Close(); err != nil {
			appLogger.Error("close app", "err", err)
		}
	}()

	if err := seed.EnsureDemoData(context.Background(), app.BotConfig, seed.DemoDataParams{
		MainDSN:  appCfg.Postgres.DSN(),
		BotID:    "shop-main",
		BotName:  "Shop Main",
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}); err != nil {
		appLogger.Error("failed to ensure demo", "err", err)
		os.Exit(1)
	}

	if err := seed.EnsureCatalogDemoData(context.Background(), app.Pool); err != nil {
		appLogger.Error("failed to ensure catalog demo data", "err", err)
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

	if err = tui.Run(app.Manager, app.BotConfig, runtimeLogs); err != nil {
		appLogger.Error("tui exited with error", "err", err)
		os.Exit(1)
	}
}
