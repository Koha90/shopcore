package botconfig_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/botconfig"
	botconfigmem "github.com/koha90/shopcore/internal/botconfig/inmemory"
)

func TestService_CreateDatabaseProfile(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://user:pass@localhost:5432/app",
		IsEnabled: true,
	})
	require.NoError(t, err)

	got, err := svc.DatabaseProfileByID(context.Background(), "main-db")
	require.NoError(t, err)
	require.Equal(t, "main-db", got.ID)
	require.Equal(t, "Main DB", got.Name)
	require.Equal(t, "postgres", got.Driver)
	require.True(t, got.IsEnabled)
}

func TestService_CreateDatabaseProfile_InvalidParams(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{})
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileIDEmpty)
}

func TestService_CreateBot(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://user:pass@localhost:5432/app",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "123456:abcdef-token",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	got, err := svc.BotByID(context.Background(), "shop-main")
	require.NoError(t, err)
	require.Equal(t, "shop-main", got.ID)
	require.Equal(t, "Shop Main", got.Name)
	require.Equal(t, "main-db", got.DatabaseID)
	require.Equal(t, "Main DB", got.DatabaseName)
	require.Equal(t, "reply_welcome", got.StartScenario)
	require.True(t, got.IsEnabled)
	require.NotEmpty(t, got.TokenMasked)
	require.NotEqual(t, "123456:abcdef-token", got.TokenMasked)
}

func TestService_CreateBot_ProfileNotFound(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "123456:abcdef-token",
		DatabaseID:    "missing-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
}

func TestService_CreateBot_StartScenarioEmpty(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postrgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrBotStartScenarioEmpty)
}

func TestService_CreateBot_StartScenarioInvalid(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Shop Main",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "abracadabra",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrBotStartScenarioInvalid)
}

func TestService_UpdateBot(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "backup-db",
		Name:      "Backup DB",
		Driver:    "postgres",
		DSN:       "postgres://backup",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "123456:abcdef-token",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main Updated",
		DatabaseID:    "backup-db",
		StartScenario: "reply_welcome",
		IsEnabled:     false,
	})
	require.NoError(t, err)

	got, err := svc.BotByID(context.Background(), "shop-main")
	require.NoError(t, err)
	require.Equal(t, "Shop Main Updated", got.Name)
	require.Equal(t, "backup-db", got.DatabaseID)
	require.Equal(t, "Backup DB", got.DatabaseName)
	require.Equal(t, "reply_welcome", got.StartScenario)
	require.False(t, got.IsEnabled)
}

func TestService_UpdateBot_ReplaceToken(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "old-token-123456",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	newToken := "new-token-654321"
	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         &newToken,
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	raw, err := bots.ByID(context.Background(), "shop-main")
	require.NoError(t, err)
	require.Equal(t, newToken, raw.Token)
}

func TestService_UpdateBot_StartScenarioEmpty(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		DatabaseID:    "main-db",
		StartScenario: "",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrBotStartScenarioEmpty)
}

func TestService_UpdateBot_StartScenarioInvalid(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		DatabaseID:    "main-db",
		StartScenario: "not-a-real-scenario",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrBotStartScenarioInvalid)
}

func TestService_ListBots(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "bot-1",
		Name:          "Bot One",
		Token:         "token-one-123",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "bot-2",
		Name:          "Bot Two",
		Token:         "token-two-456",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     false,
	})
	require.NoError(t, err)

	list, err := svc.ListBots(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestService_BotByID(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "123456:abcdef-token",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	got, err := svc.BotByID(context.Background(), "shop-main")
	require.NoError(t, err)
	require.Equal(t, "shop-main", got.ID)
	require.Equal(t, "Shop Main", got.Name)
	require.Equal(t, "Main DB", got.DatabaseName)
	require.Equal(t, "reply_welcome", got.StartScenario)
}

func TestService_UpdateBot_BotNotFound(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "missing-bot",
		Name:          "Missing Bot",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrBotNotFound)
}

func TestService_UpdateBot_ProfileNotFound(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	err = svc.UpdateBot(context.Background(), botconfig.UpdateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		DatabaseID:    "missing-db",
		StartScenario: "reply_welcome",
		IsEnabled:     true,
	})
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
}

func TestService_BotByID_NotFound(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	_, err := svc.BotByID(context.Background(), "missing")
	require.ErrorIs(t, err, botconfig.ErrBotNotFound)
}

func TestService_ListDatabaseProfiles(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "backup-db",
		Name:      "Backup DB",
		Driver:    "postgres",
		DSN:       "postgres://backup",
		IsEnabled: false,
	})
	require.NoError(t, err)

	list, err := svc.ListDatabaseProfiles(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestService_ListEnabledRuntimeBots_PropagatesStartScenario(t *testing.T) {
	bots := botconfigmem.NewBotRepository()
	dbs := botconfigmem.NewDatabaseProfileRepository()
	svc := botconfig.NewService(bots, dbs, nil)

	err := svc.CreateDatabaseProfile(context.Background(), botconfig.CreateDatabaseProfileParams{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
	})
	require.NoError(t, err)

	err = svc.CreateBot(context.Background(), botconfig.CreateBotParams{
		ID:            "shop-main",
		Name:          "Shop Main",
		Token:         "token-123456",
		DatabaseID:    "main-db",
		StartScenario: "inline_catalog",
		IsEnabled:     true,
	})
	require.NoError(t, err)

	list, err := svc.ListEnabledRuntimeBots(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.Equal(t, "inline_catalog", list[0].StartScenario)
}
