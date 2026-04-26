package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type repositoryStub struct {
	record OrderRecord
	err    error
}

func (s *repositoryStub) Create(ctx context.Context, record OrderRecord) error {
	s.record = record
	return s.err
}

func TestServiceCreate(t *testing.T) {
	t.Parallel()

	repo := &repositoryStub{}
	svc := New(repo)

	err := svc.Create(context.Background(), CreateOrderParams{
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
	})
	require.NoError(t, err)

	require.Equal(t, StatusNew, repo.record.Status)
	require.Equal(t, "shop-main", repo.record.BotID)
	require.Equal(t, int64(101), repo.record.ChatID)
	require.Equal(t, "Rose Box", repo.record.ProductName)
	require.Equal(t, "5900 ₽", repo.record.PriceText)
}

func TestBuildRecord_ValidateRequiredFields(t *testing.T) {
	t.Parallel()

	_, err := buildRecord(CreateOrderParams{})
	require.ErrorIs(t, err, ErrBotIDEmpty)
}
