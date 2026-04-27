package postgres

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

func openTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := strings.TrimSpace(os.Getenv("TEST_DATABASE_URL"))
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

func TestRepositoryCreate(t *testing.T) {
	t.Parallel()

	pool := openTestPool(t)
	repo := NewRepository(pool)

	got, err := repo.Create(context.Background(), ordersvc.OrderRecord{
		BotID:        "shop-main",
		BotName:      "Shop Main",
		ChatID:       101,
		UserID:       202,
		UserName:     "Алексей",
		UserUsername: "koha90",
		CityID:       "moscow",
		CityName:     "Москва",
		DistrictID:   "center",
		DistrictName: "Центр",
		ProductID:    "rose-box",
		ProductName:  "Rose Box",
		VariantID:    "large",
		VariantName:  "L / 25 шт",
		PriceText:    "5900 ₽",
		Status:       ordersvc.OrderStatusNew,
	})
	require.NoError(t, err)
	require.NotZero(t, got.ID)
	require.Equal(t, ordersvc.OrderStatusNew, got.Status)

	var (
		botID       string
		chatID      int64
		userID      int64
		productName string
		priceText   string
		status      string
	)

	row := pool.QueryRow(context.Background(), `
		select bot_id, chat_id, user_id, product_name, price_text, status
		from orders
		limit 1
	`)
	err = row.Scan(&botID, &chatID, &userID, &productName, &priceText, &status)
	require.NoError(t, err)

	require.Equal(t, "shop-main", botID)
	require.Equal(t, int64(101), chatID)
	require.Equal(t, int64(202), userID)
	require.Equal(t, "Rose Box", productName)
	require.Equal(t, "5900 ₽", priceText)
	require.Equal(t, string(ordersvc.OrderStatusNew), status)
}
