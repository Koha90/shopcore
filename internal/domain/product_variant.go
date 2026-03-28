package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrVariantNotFound          = errors.New("product variant not found")
	ErrInvalidProductPrice      = errors.New("invalid product price")
	ErrInvalidPackSize          = errors.New("invalid pack size")
	ErrInvalidDistrictID        = errors.New("invalid district id")
	ErrInvalidVariant           = errors.New("invalid variant")
	ErrInvalidVariantID         = errors.New("invalid variant id")
	ErrCannotArchiveLastVariant = errors.New("cannot archive last variant")
	ErrVariantAlreadyExists     = errors.New("variant already exists")
	ErrVariantAlreadyArchived   = errors.New("variant already archived")
)

// ProductVariant repesent a specific purchasable packaging option
// of a product within a specific district.
//
// It is an aggregate root and maintains its own version for
// optimistic concurrency control.
type ProductVariant struct {
	BaseAggregate

	id         int
	packSize   string
	districtID int
	price      int64

	// archivedAt is set when the variant is no longer active.
	// A nil value means the variant is active.
	archivedAt *time.Time
}

// NewProductVariant creates a new product variant.
//
// It validades all business invariants:
//
//   - packSize must not be empty
//   - districtID must be positive
//   - price must be positive
//
// Returns an error if any invariant is violated.
func NewProductVariant(
	packSize string,
	districtID int,
	price int64,
) (*ProductVariant, error) {
	if strings.TrimSpace(packSize) == "" {
		return nil, ErrInvalidPackSize
	}

	if districtID <= 0 {
		return nil, ErrInvalidDistrictID
	}

	if price <= 0 {
		return nil, ErrInvalidProductPrice
	}

	v := &ProductVariant{
		packSize:   packSize,
		districtID: districtID,
		price:      price,
	}

	v.setInitialVersion(1)
	return v, nil
}

// NewProductVariantFromDB reconstruct a ProductVariant
// from persistent storage.
//
// This function must only be used by repository implementations.
func NewProductVariantFromDB(
	id int,
	packSize string,
	distritID int,
	price int64,
	archivedAt *time.Time,
) *ProductVariant {
	v := &ProductVariant{
		id:         id,
		packSize:   packSize,
		districtID: distritID,
		price:      price,
		archivedAt: archivedAt,
	}

	v.setInitialVersion(1)
	return v
}

// ---- SETTERS ----

// ChangePrice changes the price of the variant.
//
// Increment aggregate version.
// Returns ErrInvalidProductPrice if price is non-positive.
func (v *ProductVariant) ChangePrice(price int64) error {
	if price <= 0 {
		return ErrInvalidProductPrice
	}

	v.price = price
	v.incrementVersion()
	return nil
}

// ChangePackSize updates the packaging size.
//
// Increment aggregate version.
// Returns ErrInvalidPackSize if packSize is empty.
func (v *ProductVariant) ChangePackSize(packSize string) error {
	if strings.TrimSpace(packSize) == "" {
		return ErrInvalidPackSize
	}

	v.packSize = packSize
	v.incrementVersion()
	return nil
}

// SetID is used by repository layer only.
func (v *ProductVariant) SetID(id int) {
	v.id = id
}

// ---- GETTERS ----

// ID returns product variant.
func (v *ProductVariant) ID() int {
	return v.id
}

// Price returns price of the variant product.
func (v *ProductVariant) Price() int64 {
	return v.price
}

// DistrictID returns district id of the variant product.
func (v *ProductVariant) DistrictID() int {
	return v.districtID
}

// ArchivedAt returns time were the variant product was archived.
func (v *ProductVariant) ArchivedAt() *time.Time {
	return v.archivedAt
}

// PackSize returns pack size of the variant product.
func (v *ProductVariant) PackSize() string {
	return v.packSize
}

// Version returns version of the variant product.
func (v *ProductVariant) Version() int {
	return v.version
}

// ---- CHANGERS ----

// Archive marks the variant as archived.
//
// After archiving, the variant is considered invactive.
// Increment aggregate version.
// Returns ErrVariantAlreadyArchived if archivedAt is not nil.
func (v *ProductVariant) Archive(now time.Time) error {
	if v.archivedAt != nil {
		return ErrVariantAlreadyArchived
	}
	v.incrementVersion()
	v.archivedAt = &now
	return nil
}

// IsActive returns true if the variant is the not archive.
func (v ProductVariant) IsActive() bool {
	return v.archivedAt == nil
}
