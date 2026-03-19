package postgres

import (
	"context"
	"database/sql"
	"errors"

	"botmanager/internal/botconfig"
)

type BotRepository struct {
	db *sql.DB
}

func (r *BotRepository) Save(ctx context.Context, cfg *botconfig.BotConfig) error {
	const q = `
		INSER INTO bot_configs (
			id, name, token, database_id, is_enabled, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			token = EXCLUDED.token,
			database_id = EXCLUDED.is_enabled,
			is_enabled = EXCLUDED.is_enabled,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, q, cfg.ID, cfg.Name, cfg.Token, cfg.DatabaseID, cfg.IsEnabled, cfg.UpdatedAt)
	return err
}

func (r *BotRepository) ByID(ctx context.Context, id string) (*botconfig.BotConfig, error) {
	const q = `
		SELECT id, name, token, database_id, is_enabled, updated_at
		FROM bot_configs
		WHERE id = $1
	`

	var bot botconfig.BotConfig
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&bot.ID,
		&bot.Name,
		&bot.Token,
		&bot.DatabaseID,
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

func (r *BotRepository) List(ctx context.Context) ([]botconfig.BotConfig, error) {
	const q = `
		SELECT id, name, token, database_id, is_enabled, updated_at
		FROM bot_configs
		ORDER BY id
	`

	rows, err := r.db.QueryContext(ctx, q)
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

func (r *BotRepository) Delete(ctx context.Context, id string) error {
	const q = `DELETE FROM bot_configs WHERE id = $1`

	res, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return botconfig.ErrBotNotFound
	}

	return nil
}
