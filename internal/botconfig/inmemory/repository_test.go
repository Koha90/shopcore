package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/botconfig"
	botconfigmem "github.com/koha90/shopcore/internal/botconfig/inmemory"
)

func TestBotRepository_SaveByIDListDelete(t *testing.T) {
	repo := botconfigmem.NewBotRepository()

	bot := &botconfig.BotConfig{
		ID:         "bot-1",
		Name:       "Bot One",
		Token:      "token",
		DatabaseID: "main-db",
		IsEnabled:  true,
		UpdatedAt:  time.Now(),
	}

	err := repo.Save(context.Background(), bot)
	require.NoError(t, err)

	got, err := repo.ByID(context.Background(), "bot-1")
	require.NoError(t, err)
	require.Equal(t, "Bot One", got.Name)

	list, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)

	err = repo.Delete(context.Background(), "bot-1")
	require.NoError(t, err)

	_, err = repo.ByID(context.Background(), "bot-1")
	require.ErrorIs(t, err, botconfig.ErrBotNotFound)
}

func TestDatabaseProfileRepository_SaveByIDListDelete(t *testing.T) {
	repo := botconfigmem.NewDatabaseProfileRepository()

	profile := &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://main",
		IsEnabled: true,
		UpdatedAt: time.Now(),
	}

	err := repo.Save(context.Background(), profile)
	require.NoError(t, err)

	got, err := repo.ByID(context.Background(), "main-db")
	require.NoError(t, err)
	require.Equal(t, "Main DB", got.Name)

	list, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, list, 1)

	err = repo.Delete(context.Background(), "main-db")
	require.NoError(t, err)

	_, err = repo.ByID(context.Background(), "main-db")
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
}
