package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"botmanager/internal/app/bootstrap"
	"botmanager/internal/app/seed"
	"botmanager/internal/botconfig"
	botconfigpg "botmanager/internal/botconfig/postgres"
	"botmanager/internal/config"
	"botmanager/internal/manager"
	"botmanager/internal/tui"
	"botmanager/pkg/logger"
	"botmanager/pkg/migrator"
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

	dsn := mustPostgresDSN()

	if err = migrator.MigratePostgres(dsn, "./migrations"); err != nil {
		appLogger.Error("failed to migrate database", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		appLogger.Error("failed to create pgx pool", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		appLogger.Error("failed to ping postgres", "err", err)
		os.Exit(1)
	}

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

// mustPostgresDSN builds PostgreSQL DSN from environment variables.
//
// Required environment variables:
//   - DB_HOST
//   - DB_PORT
//   - DB_USER
//   - DB_PASSWORD
//   - DB_DATABASE
//
// Optional:
//   - DB_SSLMODE (defaults to disable)
func mustPostgresDSN() string {
	host := mustEnv("DB_HOST")
	port := mustEnv("DB_PORT")
	user := mustEnv("DB_USER")
	password := mustEnv("DB_PASSWORD")
	database := mustEnv("DB_DATABASE")
	sslmode := envOrDefault("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		database,
		sslmode,
	)
}

// mustEnv returns required environment variable or exits process immediately.
func mustEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("required environment variable %s is not set", key)
	}
	return value
}

// envOrDefault returns environment variable value or fallback if empty.
func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
