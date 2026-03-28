package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrInvalidDistrictName = errors.New("invalid district name")

// District represent the district of the city.
type District struct {
	id        int
	cityID    int
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewDistrict creates a new district of city.
// Returns errror if the city id is wrong or name of
// the district is empty.
func NewDistrict(cityID int, name string) (*District, error) {
	if cityID <= 0 {
		return nil, ErrInvalidCityID
	}

	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidDistrictName
	}

	now := time.Now()
	return &District{
		cityID:    cityID,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ---- SETTERS ----

// Rename renames the district
// or returns error if new name is invalid.
func (d *District) Rename(newName string) error {
	if newName == "" {
		return ErrInvalidDistrictName
	}
	d.name = newName
	d.updatedAt = time.Now()
	return nil
}

// SetID is used by repository layer only.
func (d *District) SetID(id int) {
	d.id = id
}

// ---- GETTERS ----

// ID returns id of the district.
func (d *District) ID() int {
	return d.id
}

// CityID returns id of the city in which
// the district is located.
func (d *District) CityID() int {
	return d.cityID
}

// Name returns name of the district.
func (d *District) Name() string {
	return d.name
}

// CreatedAt returns time where the district was create.
func (d *District) CreatedAt() time.Time {
	return d.createdAt
}

// UpdatedAt returns time where the district was update.
func (d *District) UpdatedAt() time.Time {
	return d.updatedAt
}
