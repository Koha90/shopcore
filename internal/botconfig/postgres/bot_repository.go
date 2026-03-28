package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/koha90/shopcore/internal/botconfig"
)

// BotRepository stores bot configs in PostgreSQL.
type BotRepository struct {
	pool *pgxpool.Pool
}

// Save creates or updates bot configuration by ID.
func (r *BotRepository) Save(ctx context.Context, cfg *botconfig.BotConfig) error {
	const q = `
		INSERT INTO bot_configs (
			id, name, token, database_id, start_scenario, is_enabled, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			token = EXCLUDED.token,
			database_id = EXCLUDED.database_id,
			start_scenario = EXCLUDED.start_scenario,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		cfg.ID,
		cfg.Name,
		cfg.Token,
		cfg.DatabaseID,
		cfg.StartScenario,
		cfg.IsEnabled,
		cfg.UpdatedAt,
	)
	return err
}

// ByID returns bot config by ID.
func (r *BotRepository) ByID(ctx context.Context, id string) (*botconfig.BotConfig, error) {
	const q = `
		SELECT id, name, token, database_id, start_scenario, is_enabled, updated_at
		FROM bot_configs
		WHERE id = $1
	`

	var bot botconfig.BotConfig
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&bot.ID,
		&bot.Name,
		&bot.Token,
		&bot.DatabaseID,
		&bot.StartScenario,
		&bot.IsEnabled,
		&bot.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, botconfig.ErrBotNotFound
		}
		return nil, err
	}

	return &bot, nil
}

// List returns all bot configs sorted by ID.
func (r *BotRepository) List(ctx context.Context) ([]botconfig.BotConfig, error) {
	const q = `
		SELECT id, name, token, database_id, start_scenario, is_enabled, updated_at
		FROM bot_configs
		ORDER BY id
	`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []botconfig.BotConfig
	for rows.Next() {
		var bot botconfig.BotConfig
		if err := rows.Scan(
			&bot.ID,
			&bot.Name,
			&bot.Token,
			&bot.DatabaseID,
			&bot.StartScenario,
			&bot.IsEnabled,
			&bot.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, bot)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Delete removes bot config by ID.
func (r *BotRepository) Delete(ctx context.Context, id string) error {
	const q = `
		DELETE FROM bot_configs
		WHERE id = $1
	`

	tag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return botconfig.ErrBotNotFound
	}

	return nil
}
