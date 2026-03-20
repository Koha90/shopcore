package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"botmanager/internal/botconfig"
)

// TestBotRepository_SaveAndByID verifies that repository can persist
// a bot configuration and load it back by ID.
func TestBotRepository_SaveAndByID(t *testing.T) {
	pool := newTestPool(t)
	repo := &BotRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	bot := &botconfig.BotConfig{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token-main",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  now,
	}

	// Foreign key requires existing database profile.
	err := seedDatabaseProfile(ctx, pool, &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, bot)
	require.NoError(t, err)

	got, err := repo.ByID(ctx, bot.ID)
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, bot.ID, got.ID)
	require.Equal(t, bot.Name, got.Name)
	require.Equal(t, bot.Token, got.Token)
	require.Equal(t, bot.DatabaseID, got.DatabaseID)
	require.Equal(t, bot.IsEnabled, got.IsEnabled)
	require.WithinDuration(t, bot.UpdatedAt, got.UpdatedAt, time.Second)
}

// TestBotRepository_SaveUpdatesExistingBot verifies that Save works as upsert:
// it should update an existing bot row when ID already exists.
func TestBotRepository_SaveUpdatesExistingBot(t *testing.T) {
	pool := newTestPool(t)
	repo := &BotRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	err := seedDatabaseProfile(ctx, pool, &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = seedDatabaseProfile(ctx, pool, &botconfig.DatabaseProfile{
		ID:        "analytics-db",
		Name:      "Analytics DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-analytics",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	bot := &botconfig.BotConfig{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token-main",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  now,
	}

	err = repo.Save(ctx, bot)
	require.NoError(t, err)

	// Update the same entity and save again.
	bot.Name = "Shop Main Renamed"
	bot.Token = "123456:demo-token-updated"
	bot.DatabaseID = "analytics-db"
	bot.IsEnabled = false
	bot.UpdatedAt = now.Add(time.Minute)

	err = repo.Save(ctx, bot)
	require.NoError(t, err)

	got, err := repo.ByID(ctx, bot.ID)
	require.NoError(t, err)
	require.NotNil(t, err)

	require.Equal(t, "Shop Main Renamed", got.Name)
	require.Equal(t, "123456:demo-token-updated", got.Token)
	require.Equal(t, "analytics-db", got.DatabaseID)
	require.False(t, bot.IsEnabled)
	require.WithinDuration(t, bot.UpdatedAt, got.UpdatedAt, time.Second)
}

// TestBotRepository_List verifies that repository returns all bots
// sorted by ID.
func TestBotRepository_List(t *testing.T) {
	pool := newTestPool(t)
	repo := &BotRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	err := seedDatabaseProfile(ctx, pool, &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, &botconfig.BotConfig{
		ID:         "z-bot",
		Name:       "Z Bot",
		Token:      "token-z",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  now,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, &botconfig.BotConfig{
		ID:         "a-bot",
		Name:       "A Bot",
		Token:      "token-a",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  now,
	})
	require.NoError(t, err)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)

	require.Equal(t, "a-bot", list[0].ID)
	require.Equal(t, "z-bot", list[1].ID)
}

// TestBotRepository_Delete verifies that repository deletes a bot
// and that deleted bot can no longer be loaded.
func TestBotRepository_Delete(t *testing.T) {
	pool := newTestPool(t)
	repo := &BotRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	err := seedDatabaseProfile(ctx, pool, &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, &botconfig.BotConfig{
		ID:         "shop-main",
		Name:       "Shop Main",
		Token:      "123456:demo-token-main",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  now,
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, "shop-main")
	require.NoError(t, err)

	got, err := repo.ByID(ctx, "shop-main")
	require.ErrorIs(t, err, botconfig.ErrBotNotFound)
	require.Nil(t, got)
}

// TestBotRepository_DeleteNotFound verifies that Delete returns
// not-found error when target bot does not exist.
func TestBotRepository_DeleteNotFound(t *testing.T) {
	pool := newTestPool(t)
	repo := &BotRepository{pool: pool}

	ctx := context.Background()

	err := repo.Delete(ctx, "missing-bot")
	require.ErrorIs(t, err, botconfig.ErrBotNotFound)
}

// seedDatabaseProfile inserts prereuisite database profile for bot tests.
//
// We keep it in test file instead of production code becouse this is
// test fixture setup, not aplication behavior.
func seedDatabaseProfile(ctx context.Context, pool *pgxpool.Pool, profile *botconfig.DatabaseProfile) error {
	const q = `
		INSERT INTO database_profiles (
			id, name, driver, dsn, is_enabled, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := pool.Exec(
		ctx,
		q,
		profile.ID,
		profile.Name,
		profile.Driver,
		profile.DSN,
		profile.IsEnabled,
		profile.UpdatedAt,
	)

	return err
}
