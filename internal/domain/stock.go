package domain

import "errors"

// Domain errors related to Stock aggregate.
var (
	ErrInvalidQuantity           = errors.New("invalid quantity")
	ErrInsufficientStock         = errors.New("insufficient stock available")
	ErrInvalidReleaseQuantity    = errors.New("invalid release quantity")
	ErrInsufficientReservedStock = errors.New("insufficient reserved stock")
	ErrInvalidWarehouseID        = errors.New("invalid warehouse id")
	ErrInvalidStockVariantID     = errors.New("invalid variant id")
)

// Stock represents physical inventory of a product variant
// stored in a specific warehoues.
//
// This is an aggregate root responsible for maintaining
// inventory invariants:
//
//   - quantity >= 0
//   - reserved >= 0
//   - reserved <= quantity
//
// Business rules:
//   - You cannot reserve more than available.
//   - You cannot release more than reserved.
//   - You cannot decrease more than reserved.
//   - All mutating operations increment version.
type Stock struct {
	BaseAggregate

	id          int
	warehouseID int
	variantID   int
	quantity    int
	reserved    int
}

// NewStock created a new Stock aggregate.
//
// quantity must be >= 0.
// warehouseID and variantID must be valid (> 0).
func NewStock(warehouseID int, variantID int, quantity int) (*Stock, error) {
	if warehouseID <= 0 {
		return nil, ErrInvalidWarehouseID
	}

	if variantID <= 0 {
		return nil, ErrInvalidVariantID
	}

	if quantity < 0 {
		return nil, ErrInvalidQuantity
	}

	s := &Stock{
		warehouseID: warehouseID,
		variantID:   variantID,
		quantity:    quantity,
		reserved:    0,
	}

	s.setInitialVersion(1)
	return s, nil
}

// NewStockFromDB rehydrate Stock from persistence layer.
// Intended to be used only by repository.
func NewStockFromDB(
	id int,
	warehouseID int,
	variantID int,
	quantity int,
	reserved int,
) *Stock {
	s := &Stock{
		id:          id,
		warehouseID: warehouseID,
		variantID:   variantID,
		quantity:    quantity,
		reserved:    reserved,
	}

	s.setInitialVersion(1)
	return s
}

// ---- GETTERS ----

// ID returns stock indentifier.
func (s *Stock) ID() int {
	return s.id
}

// WarehouseID returns warehouse indentifier.
func (s *Stock) WarehouseID() int {
	return s.warehouseID
}

// VariantID returns related product variant indentifier.
func (s *Stock) VariantID() int {
	return s.variantID
}

// Quantity returns total quantity in warehouse.
func (s *Stock) Quantity() int {
	return s.quantity
}

// Reserved returns currently reserved quantity.
func (s *Stock) Reserved() int {
	return s.reserved
}

// Version returns aggregate version for optimistic locking.
func (s *Stock) Version() int {
	return s.version
}

// Available returns available quantity that can be reserved.
func (s *Stock) Available() int {
	return s.quantity - s.reserved
}

// Reserve reserves n units of stock.
//
// Fails if:
//   - n <= 0
//   - available stock is insufficient
func (s *Stock) Reserve(n int) error {
	if n <= 0 {
		return ErrInvalidQuantity
	}

	if s.Available() < n {
		return ErrInsufficientStock
	}

	s.reserved += n
	s.incrementVersion()
	return nil
}

// Release releases previously reserved stock.
//
// Fails if:
//   - n <= 0
//   - reserved stock is insufficient
func (s *Stock) Release(n int) error {
	if n <= 0 {
		return ErrInvalidQuantity
	}

	if s.reserved < n {
		return ErrInvalidReleaseQuantity
	}

	s.reserved -= n
	s.incrementVersion()
	return nil
}

// Decrease permanently removes reserved stock from inventory.
// Should be called after successful payment.
//
// Fails if:
//   - n <= 0
//   - reserved stock is insufficient
func (s *Stock) Decrease(n int) error {
	if n <= 0 {
		return ErrInvalidQuantity
	}

	if s.reserved < n {
		return ErrInsufficientReservedStock
	}

	s.reserved -= n
	s.quantity -= n
	s.incrementVersion()
	return nil
}

// SetID is intended for repository layer only.
func (s *Stock) SetID(id int) {
	s.id = id
}
