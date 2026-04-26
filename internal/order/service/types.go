package service

import "context"

const (
	// StatusNew marks freshly created order record.
	StatusNew = "new"
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

// OrderCreator creates persisted order records.
type OrderCreator interface {
	Create(ctx context.Context, params CreateOrderParams) error
}

// Repository stores validate order records.
type Repository interface {
	Create(ctx context.Context, record OrderRecord) error
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
	Status    string
}
