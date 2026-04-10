package seed

import (
	"context"
	"strings"

	"github.com/koha90/shopcore/internal/botconfig"
)

// BotConfigSeedService defines botconfig operations required for demo seed.
type BotConfigSeedService interface {
	ListBots(ctx context.Context) ([]botconfig.BotView, error)
	CreateBot(ctx context.Context, params botconfig.CreateBotParams) error
	CreateDatabaseProfile(ctx context.Context, params botconfig.CreateDatabaseProfileParams) error
}

// DemoDataParams configures local demo seed.
type DemoDataParams struct {
	MainDSN  string
	BotID    string
	BotName  string
	BotToken string
}

// EnsureDemoData populates storage with one local database profile and one demo bot.
//
// Seed is idempotent enough for local development:
// if at least one bot already exists, demo data is not inserted again.
//
// Behavior:
//   - database profile "main-db" is created with provided MainDSN
//   - one demo bot is created and bound to "main-db"
//   - if BotToken is empty, bot is seeded as disabled
func EnsureDemoData(ctx context.Context, svc BotConfigSeedService, params DemoDataParams) error {
	bots, err := svc.ListBots(ctx)
	if err != nil {
		return err
	}
	if len(bots) > 0 {
		return nil
	}

	mainDSN := strings.TrimSpace(params.MainDSN)
	if mainDSN == "" {
		mainDSN = "postgres://postgres:postgres@localhost:5432/shopcore?sslmode=disable"
	}

	botID := strings.TrimSpace(params.BotID)
	if botID == "" {
		botID = "shop-main"
	}

	botName := strings.TrimSpace(params.BotName)
	if botName == "" {
		botName = "Shop Main"
	}

	botToken := strings.TrimSpace(params.BotToken)
	isEnabled := botToken != ""

	if err := svc.CreateDatabaseProfile(ctx, botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       mainDSN,
		IsEnabled: true,
	}); err != nil {
		return err
	}

	if err := svc.CreateBot(ctx, botconfig.CreateBotParams{
		ID:                   botID,
		Name:                 botName,
		Token:                botToken,
		DatabaseID:           "main-db",
		StartScenario:        botconfig.StartScenarioInlineCatalog,
		TelegramAdminUserIDs: []int64{311485249},
		IsEnabled:            isEnabled,
	}); err != nil {
		return err
	}

	return nil
}
