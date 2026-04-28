package telegram

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/koha90/shopcore/internal/flow"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

func TestParseAdminOrderAction_Take(t *testing.T) {
	t.Parallel()

	orderID, status, ok := parseAdminOrderAction(buildAdminOrderActionTake(42))
	require.True(t, ok)
	require.Equal(t, int64(42), orderID)
	require.Equal(t, ordersvc.OrderStatusInProgress, status)
}

func TestParseAdminOrderAction_Close(t *testing.T) {
	t.Parallel()

	orderID, status, ok := parseAdminOrderAction(buildAdminOrderActionClose(42))
	require.True(t, ok)
	require.Equal(t, int64(42), orderID)
	require.Equal(t, ordersvc.OrderStatusClosed, status)
}

func TestParseAdminOrderAction_Invalid(t *testing.T) {
	t.Parallel()

	orderID, status, ok := parseAdminOrderAction(flow.ActionID("admin:order:take:abc"))
	require.False(t, ok)
	require.Equal(t, int64(0), orderID)
	require.Equal(t, ordersvc.OrderStatus(""), status)
}
