package service

import (
	"context"
)

// OrderStatus is persisted operator-visible order state.
type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "new"
	OrderStatusInProgress OrderStatus = "in_progress"
	OrderStatusClosed     OrderStatus = "closed"
)

// CreateOrderParams contains order data required to persist a confirmed order.
//
// It stores a customer-facing snapshot at confirmation time.
// Current step intentionally keeps price as text until numeric quote layer is added.
type CreateOrderParams struct {
	BotID   string
	BotName string

	ChatID int64
	UserID int64

	UserName     string
	UserUsername string

	CityID       string
	CityName     string
	DistrictID   string
	DistrictName string

	ProductID   string
	ProductName string
	VariantID   string
	VariantName string

	PriceText string
}

type Order struct {
	ID int64

	BotID   string
	BotName string

	ChatID int64
	UserID int64

	UserName     string
	UserUsername string

	CityID       string
	CityName     string
	DistrictID   string
	DistrictName string

	ProductID   string
	ProductName string
	VariantID   string
	VariantName string

	PriceText string
	Status    OrderStatus
}

// CreateResult contains created order identifier and current status.
type CreateResult struct {
	ID     int64
	Status OrderStatus
}

// OrderCreator returns created order identity.
type OrderCreator interface {
	Create(ctx context.Context, params CreateOrderParams) (CreateResult, error)
}

// OrderReader loads one order by ID.
type OrderReader interface {
	ByID(ctx context.Context, id int64) (Order, error)
}

// OrderStatusUpdater changes operator-visible order status.
type OrderStatusUpdater interface {
	UpdateStatus(ctx context.Context, id int64, status OrderStatus) error
}

// Repository stores validated order records.
type Repository interface {
	Create(ctx context.Context, record OrderRecord) (CreateResult, error)
	ByID(ctx context.Context, id int64) (Order, error)
	UpdateStatus(ctx context.Context, id int64, status OrderStatus) error
}

// OrderRecord is internal validated record passed to storage layer.
type OrderRecord struct {
	BotID   string
	BotName string

	ChatID int64
	UserID int64

	UserName     string
	UserUsername string

	CityID       string
	CityName     string
	DistrictID   string
	DistrictName string

	ProductID   string
	ProductName string
	VariantID   string
	VariantName string

	PriceText string
	Status    OrderStatus
}

type RuntimeService interface {
	OrderCreator
	OrderReader
	OrderStatusUpdater
}
