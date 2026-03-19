package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"botmanager/internal/botconfig"
)

// DatabaseProfileRepository stores database profiles in PostgreSQL.
type DatabaseProfileRepository struct {
	pool *pgxpool.Pool
}

// Save create or updates database profile by ID.
func (r *DatabaseProfileRepository) Save(ctx context.Context, profile *botconfig.DatabaseProfile) error {
	if profile == nil {
		return errors.New("botconfig/postgres: database profile is nil")
	}

	const q = `
		INSERT INTO database_profiles (
			id, name, driver, dsn, is_enabled, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			driver = EXCLUDED.driver,
			dsn = EXCLUDED.dsn,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(
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

// ByID returns database profile by ID.
func (r *DatabaseProfileRepository) ByID(ctx context.Context, id string) (*botconfig.DatabaseProfile, error) {
	const q = `
		SELECT id, name, driver, dsn, is_enabled, updated_at
		FROM database_profiles
		WHERE id = $1
	`

	var profile botconfig.DatabaseProfile

	err := r.pool.QueryRow(ctx, q, id).
		Scan(
			&profile.ID,
			&profile.Name,
			&profile.Driver,
			&profile.DSN,
			&profile.IsEnabled,
			&profile.UpdatedAt,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, botconfig.ErrDatabaseProfileNotFound
		}
		return nil, err
	}

	return &profile, nil
}

// List returns all database profiles sorted by ID.
func (r *DatabaseProfileRepository) List(ctx context.Context) ([]botconfig.DatabaseProfile, error) {
	const q = `
		SELECT id, name, driver, dsn, is_enabled, updated_at
		FROM database_profiles
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]botconfig.DatabaseProfile, 0)
	for rows.Next() {
		var profile botconfig.DatabaseProfile

		if err := rows.Scan(
			&profile.ID, &profile.Name, &profile.Driver, &profile.DSN, &profile.IsEnabled, &profile.UpdatedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes database profile by ID.
func (r *DatabaseProfileRepository) Delete(ctx context.Context, id string) error {
	const q = `
		DELETE FROM database_profiles
		WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return botconfig.ErrDatabaseProfileNotFound
	}

	return nil
}
