package domain

import (
	"errors"
	"time"
)

// OrderStatus represents lifecycle of order.
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

var (
	ErrInvalidOrderUserID    = errors.New("invalid order user id")
	ErrOrderEmpty            = errors.New("order must contain items")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrOrderAlreadyCancelled = errors.New("order already cancelled")
	ErrOrderNotPending       = errors.New("order is not pending")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderUpdate           = errors.New("failed to update order")
	ErrOrderSave             = errors.New("failed to save order")
	ErrOrderCancel           = errors.New("failed to cancel order")
)

// Order represents confirmed purchase intent.
//
// Business rules:
//   - Created only with items.
//   - Total is immutable after creation.
//   - Only pending order can be paid or cancelled.
//   - Paid order cannot be cancelled.
//   - Cancelled order cannot be paid.
type Order struct {
	BaseAggregate

	id          int
	userID      int
	items       []OrderItem
	total       int64
	status      OrderStatus
	createdAt   time.Time
	paidAt      *time.Time
	cancelledAt *time.Time
}

// NewOrder creates new pending order.
func NewOrder(userID int, items []OrderItem, createdAt time.Time) (*Order, error) {
	if userID <= 0 {
		return nil, ErrInvalidOrderUserID
	}

	if len(items) == 0 {
		return nil, ErrOrderEmpty
	}

	var total int64
	for _, item := range items {
		total += int64(item.quantity) * item.unitPrice
	}

	o := &Order{
		userID:    userID,
		items:     items,
		total:     total,
		status:    OrderStatusPending,
		createdAt: createdAt,
	}

	o.setInitialVersion(1)

	return o, nil
}

// NewOrderFromDB reconstructs order from persistent storage.
//
// It is intended for repository layer only.
func NewOrderFromDB(
	id int,
	userID int,
	items []OrderItem,
	total int64,
	status OrderStatus,
	varsion int,
	createdAt time.Time,
	paidAt *time.Time,
	cancelledAt *time.Time,
) (*Order, error) {
	if id <= 0 {
		return nil, ErrOrderNotFound
	}
	if userID <= 0 {
		return nil, ErrInvalidOrderUserID
	}
	if len(items) == 0 {
		return nil, ErrOrderEmpty
	}

	o := &Order{
		id:          id,
		userID:      userID,
		items:       items,
		total:       total,
		status:      status,
		createdAt:   createdAt,
		paidAt:      paidAt,
		cancelledAt: cancelledAt,
	}
	o.setInitialVersion(varsion)

	return o, nil
}

// ---- GETTERS ----

// ID returns order id.
func (o *Order) ID() int {
	return o.id
}

// UserID returns user id.
func (o *Order) UserID() int {
	return o.userID
}

// Status returns order status.
func (o *Order) Status() OrderStatus {
	return o.status
}

// Total returns order total amount.
func (o *Order) Total() int64 {
	return o.total
}

// Items returns copy of order items.
func (o *Order) Items() []OrderItem {
	result := make([]OrderItem, len(o.items))
	copy(result, o.items)
	return result
}

// Version returns aggregate version.
func (o *Order) Version() int {
	return o.version
}

// CreatedAt returns order creation time.
func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}

// PaidAt returns payment time if present.
func (o *Order) PaidAt() *time.Time {
	return o.paidAt
}

// CancelledAt returns cancellation time if present.
func (o *Order) CancelledAt() *time.Time {
	return o.cancelledAt
}

// MarkPaid marks order as paid.
//
// Fails if:
//   - already paid
//   - already cancelled
//   - not pending
func (o *Order) MarkPaid(now time.Time) error {
	if o.status == OrderStatusPaid {
		return ErrOrderAlreadyPaid
	}

	if o.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if o.status != OrderStatusPending {
		return ErrOrderNotPending
	}

	o.status = OrderStatusPaid
	o.paidAt = &now

	o.incrementVersion()
	o.addEvent(NewOrderPaid(o.id))

	return nil
}

// Cancel cancels order.
//
// Fails if:
//   - already cancelled
//   - already paid
func (o *Order) Cancel(now time.Time) error {
	if o.status == OrderStatusCancelled {
		return ErrOrderAlreadyCancelled
	}

	if o.status == OrderStatusPaid {
		return ErrOrderAlreadyPaid
	}

	if o.status != OrderStatusPending {
		return ErrOrderNotPending
	}

	o.status = OrderStatusCancelled
	o.cancelledAt = &now

	o.incrementVersion()
	o.addEvent(NewOrderCancelled(o.id))

	return nil
}

// ---- SETTERS ----

// SetID is intended for repository layer only.
func (o *Order) SetID(id int) {
	o.id = id
}
