package service

import (
	"context"
	"fmt"
	"strings"
)

// DistrictWriter stores district data.
type DistrictWriter interface {
	CreateDistrict(ctx context.Context, params CreateDistrictParams) error
}

// CreateDistrictParams contains data required to create one district.
type CreateDistrictParams struct {
	CityID    int
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

// CreateDistrict validates input and stores a new district.
func (s *Service) CreateDistrict(ctx context.Context, params CreateDistrictParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.districts == nil {
		return fmt.Errorf("district writer is nil")
	}

	params.Code = normalizeCode(params.Code)
	params.Name = strings.TrimSpace(params.Name)
	params.NameLatin = strings.TrimSpace(params.NameLatin)

	switch {
	case params.CityID <= 0:
		return ErrDistrictCityIDInvalid
	case params.Code == "":
		return ErrDistrictCodeEmpty
	case params.Name == "":
		return ErrDistrictNameEmpty
	}

	return s.districts.CreateDistrict(ctx, params)
}
