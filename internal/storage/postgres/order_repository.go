package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/koha90/shopcore/internal/domain"
	"github.com/koha90/shopcore/internal/service"
)

var _ service.OrderRepository = (*OrderRepo)(nil)

// OrderRepo stores orders in PostgreSQL.
type OrderRepo struct {
	db     *sql.DB
	logger *slog.Logger
}

// NewOrderRepo creates a new PostgreSQL order repository.
//
// logger may be nil. In that case slog.Default() is used.
func NewOrderRepo(db *sql.DB, logger *slog.Logger) *OrderRepo {
	if db == nil {
		panic("postgres: db is nil")
	}
	if logger == nil {
		logger = slog.Default()
	}

	return &OrderRepo{
		db:     db,
		logger: logger,
	}
}

// Save persists order state.
//
// New order are inserted, existing orders are updated.
func (r *OrderRepo) Save(ctx context.Context, order *domain.Order) error {
	if order.ID() == 0 {
		return r.insert(ctx, order)
	}
	return r.update(ctx, order)
}

// ByID returns order by its identifier.
func (r *OrderRepo) ByID(ctx context.Context, id int) (*domain.Order, error) {
	const q = `
			SELECT id, user_id, total, status, version, created_at, paid_at, cancelled_at
			FROM orders
			WHERE id = $1
	`

	var (
		orderID     int
		userID      int
		total       int64
		status      string
		version     int
		createdAt   time.Time
		paidAt      sql.NullTime
		cancelledAt sql.NullTime
	)

	row := r.queryRow(ctx, q, id)
	if err := row.Scan(
		&orderID,
		&userID,
		&total,
		&status,
		&version,
		&createdAt,
		&paidAt,
		&cancelledAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	items, err := r.loadItems(ctx, orderID)
	if err != nil {
		return nil, err
	}

	var paidPtr *time.Time
	if paidAt.Valid {
		paidPtr = &paidAt.Time
	}

	var cancelledPtr *time.Time
	if cancelledAt.Valid {
		cancelledPtr = &cancelledAt.Time
	}

	return domain.NewOrderFromDB(
		orderID,
		userID,
		items,
		total,
		domain.OrderStatus(status),
		version,
		createdAt,
		paidPtr,
		cancelledPtr,
	)
}

func (r *OrderRepo) insert(ctx context.Context, order *domain.Order) error {
	const q = `
		INSERT INTO orders (user_id, total, status, version, created_at, paid_at, cancelled_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		REURNING id
	`

	var id int
	if err := r.queryRow(
		ctx,
		q,
		order.UserID(),
		order.Total(),
		string(order.Status()),
		order.Version(),
		order.CreatedAt(),
		order.PaidAt(),
		order.CancelledAt(),
	).Scan(&id); err != nil {
		return err
	}

	order.SetID(id)
	return r.saveItems(ctx, order)
}

func (r *OrderRepo) update(ctx context.Context, order *domain.Order) error {
	const q = `
		UPDATE orders
		SET total = $1,
				status = $2,
				version = $3,
				created_at = $4,
				paid_at = $5,
				cancelled_at = $6
		WHERE id = $7
	`

	_, err := r.exec(
		ctx,
		q,
		order.Total(),
		string(order.Status()),
		order.Version(),
		order.CreatedAt(),
		order.PaidAt(),
		order.CancelledAt(),
		order.ID(),
	)
	if err != nil {
		return err
	}

	const deleteItems = `DELETE FROM order_items WHERE order_id = $1`
	if _, err := r.exec(ctx, deleteItems, order.ID()); err != nil {
		return err
	}

	return r.saveItems(ctx, order)
}

func (r *OrderRepo) saveItems(ctx context.Context, order *domain.Order) error {
	const q = `
		INSERT INTO order_items (order_id, product_id, variant_id, quantity, unit_price)
		VALUES ($1, $2, $3, $4, $5)
	`

	for _, item := range order.Items() {
		if _, err := r.exec(
			ctx,
			q,
			order.ID(),
			item.ProductID(),
			item.VariantID(),
			item.Quantity(),
			item.UnitPrice(),
		); err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderRepo) loadItems(ctx context.Context, orderID int) ([]domain.OrderItem, error) {
	const q = `
		SELECT product_id, variant_id, quantity, unit_price
		FROM order_items
		WHERE order_id = $1
		ORDER BY id
	`

	rows, err := r.query(ctx, q, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.OrderItem

	for rows.Next() {
		var (
			productID int
			variantID int
			quantity  int
			unitPrice int64
		)

		if err := rows.Scan(&productID, &variantID, &quantity, &unitPrice); err != nil {
			return nil, err
		}

		items = append(items, domain.NewOrderItem(productID, variantID, quantity, unitPrice))
	}

	return items, rows.Err()
}

func (r *OrderRepo) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}

func (r *OrderRepo) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}
	return r.db.QueryContext(ctx, query, args...)
}

func (r *OrderRepo) queryRow(ctx context.Context, query string, args ...any) *sql.Row {
	if tx, ok := txFromContext(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}
	return r.db.QueryRowContext(ctx, query, args...)
}
