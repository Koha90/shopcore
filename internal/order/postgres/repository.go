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
func (r *Repository) Create(ctx context.Context, record ordersvc.OrderRecord) (ordersvc.CreateResult, error) {
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
		returning id, status
	`

	var result ordersvc.CreateResult
	var status string

	err := r.pool.QueryRow(
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
		string(record.Status),
	).Scan(&result.ID, &status)
	if err != nil {
		return ordersvc.CreateResult{}, fmt.Errorf("insert order: %w", err)
	}

	result.Status = ordersvc.OrderStatus(status)
	return result, nil
}

func (r *Repository) ByID(ctx context.Context, id int64) (ordersvc.Order, error) {
	const q = `
			select
					id,
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
					status,
			from orders
			where id = $1
	`

	var order ordersvc.Order
	var status string

	err := r.pool.QueryRow(ctx, q, id).Scan(
		&order.ID,
		&order.BotID,
		&order.BotName,
		&order.ChatID,
		&order.UserID,
		&order.UserName,
		&order.UserUsername,
		&order.CityID,
		&order.CityName,
		&order.DistrictID,
		&order.DistrictName,
		&order.ProductID,
		&order.ProductName,
		&order.VariantID,
		&order.VariantName,
		&order.PriceText,
		&status,
	)
	if err != nil {
		return ordersvc.Order{}, fmt.Errorf("select order by id: %w", err)
	}

	order.Status = ordersvc.OrderStatus(status)
	return order, nil
}

func (r *Repository) UpdateStatus(ctx context.Context, id int64, status ordersvc.OrderStatus) error {
	const q = `
			update orders
			set status = $2
			where id = $1
	`

	_, err := r.pool.Exec(ctx, q, string(status))
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}
