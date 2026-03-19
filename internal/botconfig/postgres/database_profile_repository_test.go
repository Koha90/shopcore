package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"botmanager/internal/botconfig"
)

func TestDatabaseProfileRepository_SaveAndByID(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	profile := &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: now,
	}

	err := repo.Save(ctx, profile)
	require.NoError(t, err)

	got, err := repo.ByID(ctx, "main-db")
	require.NoError(t, err)
	require.NotNil(t, got)

	require.Equal(t, profile.ID, got.ID)
	require.Equal(t, profile.Name, got.Name)
	require.Equal(t, profile.Driver, got.Driver)
	require.Equal(t, profile.DSN, got.DSN)
	require.Equal(t, profile.IsEnabled, got.IsEnabled)
	require.WithinDuration(t, profile.UpdatedAt, got.UpdatedAt, time.Second)
}

func TestDatabaseProfileRepository_SaveUpdatesExistingProfile(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()

	profile := &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: time.Now().UTC().Round(time.Second),
	}

	err := repo.Save(ctx, profile)
	require.NoError(t, err)

	profile.Name = "Main DB Renamed"
	profile.DSN = "postgres://demo-main-updated"
	profile.IsEnabled = false
	profile.UpdatedAt = time.Now().UTC().Add(time.Minute).Round(time.Second)

	err = repo.Save(ctx, profile)
	require.NoError(t, err)

	got, err := repo.ByID(ctx, "main-db")
	require.NoError(t, err)

	require.Equal(t, "Main DB Renamed", got.Name)
	require.Equal(t, "postgres://demo-main-updated", got.DSN)
	require.Equal(t, false, got.IsEnabled)
}

func TestDatabaseProfileRepository_ByIDNotFound(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()

	got, err := repo.ByID(ctx, "missing-db")
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
	require.Nil(t, got)
}

func TestDatabaseProfileRepository_List(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()
	now := time.Now().UTC().Round(time.Second)

	err := repo.Save(ctx, &botconfig.DatabaseProfile{
		ID:        "z-db",
		Name:      "Z DB",
		Driver:    "postgres",
		DSN:       "postgres://z-db",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, &botconfig.DatabaseProfile{
		ID:        "a-db",
		Name:      "A DB",
		Driver:    "postgres",
		DSN:       "postgres://a-db",
		IsEnabled: true,
		UpdatedAt: now,
	})
	require.NoError(t, err)

	list, err := repo.List(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)
	require.Equal(t, "a-db", list[0].ID)
	require.Equal(t, "z-db", list[1].ID)
}

func TestDatabaseProfileRepository_Delete(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()

	err := repo.Save(ctx, &botconfig.DatabaseProfile{
		ID:        "main-db",
		Name:      "Main DB",
		Driver:    "postgres",
		DSN:       "postgres://demo-main",
		IsEnabled: true,
		UpdatedAt: time.Now().UTC().Round(time.Second),
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, "main-db")
	require.NoError(t, err)

	got, err := repo.ByID(ctx, "main-db")
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
	require.Nil(t, got)
}

func TestDatabaseProfileRepository_DeleteNotFound(t *testing.T) {
	t.Parallel()

	pool := newTestPool(t)
	repo := &DatabaseProfileRepository{pool: pool}

	ctx := context.Background()

	err := repo.Delete(ctx, "missing-db")
	require.ErrorIs(t, err, botconfig.ErrDatabaseProfileNotFound)
}
