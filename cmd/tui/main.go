package main

import (
	"context"
	"errors"
	"log"
	"time"

	"botmanager/internal/app/bootstrap"
	"botmanager/internal/botconfig"
	"botmanager/internal/botconfig/inmemory"
	"botmanager/internal/manager"
	"botmanager/internal/tui"
)

type demoRunner struct{}

// Run simulates bot runtime lifecycle for local demo.
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
	ctx := context.Background()

	m := manager.New(&demoRunner{})

	botsRepo := inmemory.NewBotRepository()
	dbRepo := inmemory.NewDatabaseProfileRepository()
	cfgSvc := botconfig.NewService(botsRepo, dbRepo, nil)

	_ = cfgSvc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "analytics-db",
		Name:      "Analytics DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-analytics",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token-main",
		DatabaseID: "main-db",
		IsEnabled:  true,
	})

	_ = cfgSvc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "slow-bot",
		Name:       "Slow Bot",
		Token:      "123456:demo-token-slow",
		DatabaseID: "analytics-db",
		IsEnabled:  true,
	})

	_ = cfgSvc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "broken-bot",
		Name:       "Broken Bot",
		Token:      "123456:demo-token-broken",
		DatabaseID: "main-db",
		IsEnabled:  true,
	})

	starter := bootstrap.NewStarter(cfgSvc, m)

	results, err := starter.StartEnabled(ctx)
	for _, result := range results {
		if result.Err != nil {
			log.Printf("bootstrap bot=%s failed: %v", result.ID, result.Err)
			continue
		}

		log.Printf(
			"bootstrap bot=%s registered=%t started=%t",
			result.ID,
			result.Registered,
			result.Started,
		)
	}

	if err != nil {
		log.Printf("bootstrap completed with errors: %v", err)
	}

	if err := tui.Run(m, cfgSvc); err != nil {
		panic(err)
	}
}
