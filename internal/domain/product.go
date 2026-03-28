package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrProductNotFound    = errors.New("product not found")
	ErrInvalidProductName = errors.New("invalid product name")
	ErrInvalidCategoryID  = errors.New("invalid category id")
	ErrInvalidImageURL    = errors.New("invalid image url")
	ErrInvalidProductID   = errors.New("invalid product id")
)

// Product represents the product.
type Product struct {
	BaseAggregate

	id          int
	categoryID  *int
	name        string
	description string
	imagePath   *string
	variants    []ProductVariant
}

// NewProduct creates a new product.
// The returned product is not persisted yet,
// or error if name empty or invalids
// and category id is invalid.
func NewProduct(
	name string,
	categoryID int,
	description string,
	imagePath string,
) (*Product, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidProductName
	}

	if categoryID <= 0 {
		return nil, ErrInvalidCategoryID
	}

	p := &Product{
		name:        name,
		categoryID:  &categoryID,
		description: description,
		imagePath:   &imagePath,
	}

	p.setInitialVersion(1)
	return p, nil
}

func NewProductFromDB(
	id int,
	categoryID *int,
	name string,
	description string,
	imagePath *string,
	version int,
	variants []ProductVariant,
) *Product {
	p := &Product{
		id:          id,
		categoryID:  categoryID,
		name:        name,
		description: description,
		imagePath:   imagePath,
		variants:    variants,
	}

	p.setInitialVersion(1)
	return p
}

// ---- GETTERS ----

// ID return id of product.
func (p *Product) ID() int {
	return p.id
}

// CategoryID returns identifier category of product.
func (p *Product) CategoryID() *int {
	return p.categoryID
}

// Name return name of product.
func (p *Product) Name() string {
	return p.name
}

// Description returns description of product.
func (p *Product) Description() string {
	return p.description
}

// ImagePath returns imageurl of product.
func (p *Product) ImagePath() *string {
	return p.imagePath
}

// Version returns version of product.
func (p *Product) Version() int {
	return p.version
}

// Variants returns all variants of prodcut.
func (p *Product) Variants() []ProductVariant {
	result := make([]ProductVariant, len(p.variants))
	copy(result, p.variants)

	return result
}

// VariantByID returns variant product of product
// or error if variant was not found.
func (p *Product) VariantByID(id int) (*ProductVariant, error) {
	for i := range p.variants {
		v := &p.variants[i]
		if v.ID() == id && v.IsActive() {
			return v, nil
		}
	}

	return nil, ErrVariantNotFound
}

// ActiveVariants returns of copy the list variants if this active.
func (p *Product) ActiveVariants() []ProductVariant {
	var result []ProductVariant
	for _, v := range p.variants {
		if v.IsActive() {
			result = append(result, v)
		}
	}

	return result
}

// ---- CHANEGERS ----

// Rename renames the product.
func (p *Product) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidProductName
	}

	p.name = name
	p.incrementVersion()
	return nil
}

// ChangeCategory changes category.
func (p *Product) ChangeCategory(categoryID int) error {
	if categoryID <= 0 {
		return ErrInvalidCategoryID
	}
	p.categoryID = &categoryID
	return nil
}

// UpdateDescription updates description the product.
func (p *Product) UpdateDescription(description string) {
	p.description = description
}

// TODO: create other getter-methods.

// ---- SETTERS ----

// SetID is used by repository layer only.
func (p *Product) SetID(id int) {
	p.id = id
}

// SetVariantID sets variant id.
func (p *Product) SetVariantID(index int, id int) {
	if index >= 0 && index < len(p.variants) {
		p.variants[index].SetID(id)
	}
}

func (p *Product) AddVariant(
	packSize string,
	districtID int,
	price int64,
) error {
	for _, v := range p.variants {
		if v.packSize == packSize &&
			v.districtID == districtID &&
			v.IsActive() {
			return ErrVariantAlreadyExists
		}
	}

	v, err := NewProductVariant(packSize, districtID, price)
	if err != nil {
		return err
	}

	p.variants = append(p.variants, *v)
	p.incrementVersion()
	p.addEvent(NewProductVariantAdded(v.ID()))
	return nil
}

func (p *Product) HasVariants() bool {
	return len(p.variants) > 0
}

// ArchiveVariant sets product variant to archive.
func (p *Product) ArchiveVariant(id int, now time.Time) error {
	activeCount := 0
	var target *ProductVariant

	for i := range p.variants {
		if p.variants[i].IsActive() {
			activeCount++
		}
		if p.variants[i].ID() == id && p.variants[i].IsActive() {
			target = &p.variants[i]
		}
	}

	if target == nil {
		return ErrVariantNotFound
	}

	if activeCount <= 1 && target.IsActive() {
		return ErrCannotArchiveLastVariant
	}

	if err := target.Archive(now); err != nil {
		return err
	}
	p.incrementVersion()
	p.addEvent(NewVariantArchived(id))
	return nil
}

// VariantsForUpdate returns ptr of ProductVariant.
func (p *Product) VariantsForUpdate() []*ProductVariant {
	var result []*ProductVariant
	for i := range p.variants {
		result = append(result, &p.variants[i])
	}
	return result
}
