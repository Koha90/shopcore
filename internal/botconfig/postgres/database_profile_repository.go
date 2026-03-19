package postgres

import (
	"context"
	"database/sql"
	"errors"

	"botmanager/internal/botconfig"
)

type DatabaseProfileRepository struct {
	db *sql.DB
}

func (r *DatabaseProfileRepository) Save(ctx context.Context, profile *botconfig.DatabaseProfile) error {
	const q = `
		INSERT INTO database_profiles (
			id, name, driver, dsn, is_enabled, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name
			driver = EXCLUDED.driver
			dsn = EXCLUDED.dsn
			is_enabled = EXCLUDED.is_enabled
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
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

func (r *DatabaseProfileRepository) ByID(ctx context.Context, id string) (*botconfig.DatabaseProfile, error) {
	const q = `
		SELECT id, name, driver, dsn, is_enabled, updated_at
		FROM database_profiles
		WHERE id = $1
	`

	var profile botconfig.DatabaseProfile
	err := r.db.QueryRowContext(ctx, q, id).
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

func (r *DatabaseProfileRepository) List(ctx context.Context) ([]botconfig.DatabaseProfile, error) {
	const q = `
		SELECT id, name, driver, dsn, is_enabled, updated_at
		FROM database_profiles
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []botconfig.DatabaseProfile
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

func (r *DatabaseProfileRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM database_profiles WHERE id = $1`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return botconfig.ErrDatabaseProfileNotFound
	}

	return nil
}
