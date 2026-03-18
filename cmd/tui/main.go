package main

import (
	"context"
	"errors"
	"time"

	"botmanager/internal/botconfig"
	"botmanager/internal/botconfig/inmemory"
	"botmanager/internal/manager"
	"botmanager/internal/tui"
)

type demoRunner struct{}

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
	m := manager.New(&demoRunner{})

	must(m.Register(manager.BotSpec{
		ID:    "shop-main",
		Name:  "Shop Main",
		Token: "token-main",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "slow-bot",
		Name:  "Slow Bot",
		Token: "token-slow",
	}))

	must(m.Register(manager.BotSpec{
		ID:    "broken-bot",
		Name:  "Broken Bot",
		Token: "token-broken",
	}))

	botsRepo := inmemory.NewBotRepository()
	dbRepo := inmemory.NewDatabaseProfileRepository()
	cfgSvc := botconfig.NewService(botsRepo, dbRepo, nil)

	_ = cfgSvc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "analytics-db",
		Name:      "Analytics DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-analytics",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "staging-db",
		Name:      "Staging DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-staging",
		IsEnabled: true,
	})

	_ = cfgSvc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token",
		DatabaseID: "main-db",
		IsEnabled:  true,
	})

	_ = cfgSvc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:         "slow-bot",
		Name:       "Slow Bot",
		Token:      "123456:slow-demo-token",
		DatabaseID: "main-db",
		IsEnabled:  true,
	})
	_ = cfgSvc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:         "broken-bot",
		Name:       "Broken Bot",
		Token:      "123456:broken-demo-token",
		DatabaseID: "main-db",
		IsEnabled:  true,
	})
	if err := tui.Run(m, cfgSvc); err != nil {
		panic(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
