package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCityName = errors.New("invalid city name")
	ErrInvalidCityID   = errors.New("invalid city id")
)

// City represent the city.
type City struct {
	id        int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewCity create a new city.
// Returns error if name is empty.
func NewCity(name string) (*City, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidCityName
	}

	now := time.Now()
	return &City{
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ---- SETTERS ----

// Rename renames the city
// or returns error if new name of city is empty.
func (c *City) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidCityName
	}
	c.name = newName
	c.updatedAt = time.Now()
	return nil
}

// SetID is used by repository layer only.
func (c *City) SetID(id int) {
	c.id = id
}

// ---- GETTERS ----

// ID returns identifier of the city.
func (c *City) ID() int {
	return c.id
}

// Name returns name of the city.
func (c *City) Name() string {
	return c.name
}

// CreatedAt returns time where created the city.
func (c *City) CreatedAt() time.Time {
	return c.createdAt
}

// UpdatedAt returns time where updated the city.
func (c *City) UpdatedAt() time.Time {
	return c.updatedAt
}
