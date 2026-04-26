package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

// Repository stores orders in PostgreSQL
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates Postrgres order repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Create inserts one confirmed order record.
func (r *Repository) Create(ctx context.Context, record ordersvc.OrderRecord) error {
	const q = `
		insert into orders (
				bot_id,
				bot_name,
				chat_id,
				user_id,
				user_name,
				user_username,
				city_id,
				city_name,
				district_id,
				district_name,
				product_id,
				product_name,
				variant_id,
				variant_name,
				price_text,
				status
		)
		values (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11, $12,
				$13, $14, $15, $16
		)
	`

	_, err := r.pool.Exec(
		ctx,
		q,
		record.BotID,
		record.BotName,
		record.ChatID,
		record.UserID,
		record.UserName,
		record.UserUsername,
		record.CityID,
		record.CityName,
		record.DistrictID,
		record.DistrictName,
		record.ProductID,
		record.ProductName,
		record.VariantID,
		record.VariantName,
		record.PriceText,
		record.Status,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	return nil
}
