package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidCategoryName = errors.New("invalid category name")

// Category represent category of the product.
type Category struct {
	id          int
	name        string
	description string
	createdAt   time.Time
	updatedAt   time.Time
}

// NewCategory create a new category or
// returns an error if name is empty.
func NewCategory(name string, description string) (*Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidCategoryName
	}

	now := time.Now()

	return &Category{
		name:        name,
		description: description,
		createdAt:   now,
		updatedAt:   now,
	}, nil
}

// ---- SETTERS ----

// Rename renames the category.
// Returns error if the new name is empty.
func (c *Category) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidCategoryName
	}
	c.name = newName
	c.updatedAt = time.Now()
	return nil
}

// UpdateDescription update description of the category.
func (c *Category) UpdateDescription(newDesc string) {
	c.description = newDesc
	c.updatedAt = time.Now()
}

// SetID is used by repository layer only.
func (c *Category) SetID(id int) {
	c.id = id
}

// ---- GETTERS ----

// ID returns id of the category.
func (c *Category) ID() int {
	return c.id
}

// Name returns name of the category.
func (c *Category) Name() string {
	return c.name
}

// Desecription returns description of the category.
func (c *Category) Desecription() string {
	return c.description
}

// CreatedAt returns time where the category was created.
func (c *Category) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns time where the category was updated.
func (c *Category) UpdatedAt() time.Time {
	return c.updatedAt
}
