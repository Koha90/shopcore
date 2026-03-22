package seed

import (
	"context"

	"botmanager/internal/botconfig"
)

// BotConfigSeedService defines botconfig operations required for demo seed.
type BotConfigSeedService interface {
	ListBots(ctx context.Context) ([]botconfig.BotView, error)
	CreateBot(ctx context.Context, params botconfig.CreateBotParams) error
	CreateDatabaseProfile(ctx context.Context, params botconfig.CreateDatabaseProfileParams) error
}

// EnsureDemoData populates storage with demo database profiles and bots.
//
// Seed is idempotent enough for local development:
// if at least one bot already exists, demo data is not inserted again.
func EnsureDemoData(ctx context.Context, svc BotConfigSeedService) error {
	bots, err := svc.ListBots(ctx)
	if err != nil {
		return err
	}
	if len(bots) > 0 {
		return nil
	}

	if err := svc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
	}); err != nil {
		return err
	}

	if err := svc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "analytics-db",
		Name:      "Analytics DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-analytics",
		IsEnabled: true,
	}); err != nil {
		return err
	}

	if err := svc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "staging-db",
		Name:      "Staging DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-staging",
		IsEnabled: true,
	}); err != nil {
		return err
	}

	if err := svc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token-main",
		DatabaseID: "main-db",
		IsEnabled:  true,
	}); err != nil {
		return err
	}

	if err := svc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "slow-bot",
		Name:       "Slow Bot",
		Token:      "123456:demo-token-slow",
		DatabaseID: "analytics-db",
		IsEnabled:  true,
	}); err != nil {
		return err
	}

	if err := svc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:         "broken-bot",
		Name:       "Broken Bot",
		Token:      "123456:demo-token-broken",
		DatabaseID: "staging-db",
		IsEnabled:  true,
	}); err != nil {
		return err
	}

	return nil
}
