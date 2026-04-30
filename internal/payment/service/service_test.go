package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type botPaymentMethodReaderStub struct {
	gotBotID string
	out      []BotPaymentMethod
	err      error
}

func (s *botPaymentMethodReaderStub) ListBotPaymentMethods(
	ctx context.Context,
	botID string,
) ([]BotPaymentMethod, error) {
	s.gotBotID = botID

	return s.out, s.err
}

func TestListBotPaymentMethods(t *testing.T) {
	t.Parallel()

	reader := &botPaymentMethodReaderStub{
		out: []BotPaymentMethod{
			{
				ID:              1,
				Code:            "sbp_main",
				Name:            "СБП",
				Kind:            PaymentKindSBP,
				DisplayName:     "СБП перевод",
				ExtraPercentBPS: 150,
				SortOrder:       10,
			},
		},
	}

	svc := New(reader)

	got, err := svc.ListBotPaymentMethods(context.Background(), " shop-main ")

	require.NoError(t, err)
	require.Len(t, got, 1)

	assert.Equal(t, "shop-main", reader.gotBotID)
	assert.Equal(t, "sbp_main", got[0].Code)
	assert.Equal(t, PaymentKindSBP, got[0].Kind)
	assert.Equal(t, "СБП перевод", got[0].Label())
	assert.Equal(t, 150, got[0].ExtraPercentBPS)
}

func TestListBotPaymentMethodsRejectsEmptyBotID(t *testing.T) {
	t.Parallel()

	svc := New(&botPaymentMethodReaderStub{})

	_, err := svc.ListBotPaymentMethods(context.Background(), "   ")

	require.ErrorIs(t, err, ErrBotIDEmpty)
}

func TestBotPaymentMethodLabelFallsBackToName(t *testing.T) {
	t.Parallel()

	method := BotPaymentMethod{
		Name: "Карта банка",
	}

	assert.Equal(t, "Карта банка", method.Label())
}
