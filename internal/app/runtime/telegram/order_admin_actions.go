package telegram

import (
	"strconv"
	"strings"

	"github.com/koha90/shopcore/internal/flow"
	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

const (
	adminOrderActionTakePrefix  = "admin:order:take:"
	adminOrderActionClosePrefix = "admin:order:close:"
)

// buildAdminOrderActionTake builds callback action for taking order into work.
func buildAdminOrderActionTake(orderID int64) flow.ActionID {
	return flow.ActionID(adminOrderActionTakePrefix + strconv.FormatInt(orderID, 10))
}

// buildAdminOrderActionClose build callback action for closing order.
func buildAdminOrderActionClose(orderID int64) flow.ActionID {
	return flow.ActionID(adminOrderActionClosePrefix + strconv.FormatInt(orderID, 10))
}

// parseAdminOrderAction resolves admin order action into order id and target status.
func parseAdminOrderAction(actionID flow.ActionID) (int64, ordersvc.OrderStatus, bool) {
	raw := string(actionID)
	switch {
	case strings.HasPrefix(raw, adminOrderActionTakePrefix):
		id, err := strconv.ParseInt(strings.TrimPrefix(raw, adminOrderActionTakePrefix), 10, 64)
		if err != nil || id <= 0 {
			return 0, "", false
		}
		return id, ordersvc.OrderStatusInProgress, true

	case strings.HasPrefix(raw, adminOrderActionClosePrefix):
		id, err := strconv.ParseInt(strings.TrimPrefix(raw, adminOrderActionClosePrefix), 10, 64)
		if err != nil || id <= 0 {
			return 0, "", false
		}

		return id, ordersvc.OrderStatusClosed, true

	default:
		return 0, "", false
	}
}
