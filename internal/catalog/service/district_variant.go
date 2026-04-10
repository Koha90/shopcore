package service

import (
	"context"
	"fmt"
)

// DistrictVariantWriter stores district variant placement data.
type DistrictVariantWriter interface {
	CreateDistrictVariant(ctx context.Context, params CreateDistrictVariantParams) error
}

// CreateDistrictVariantParams contains data required to place one variant in one district.
type CreateDistrictVariantParams struct {
	DistrictID int
	VariantID  int
	Price      int
}

// CreateDistrictVariant validates input and stores one district-variant placement.
func (s *Service) CreateDistrictVariant(ctx context.Context, params CreateDistrictVariantParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.districtVariants == nil {
		return fmt.Errorf("district variant writer is nil")
	}

	switch {
	case params.DistrictID <= 0:
		return ErrDistrictVariantDistrictIDInvalid
	case params.VariantID <= 0:
		return ErrDistrictVariantVariantIDInvalid
	case params.Price <= 0:
		return ErrDistrictVariantPriceInvalid
	}

	return s.districtVariants.CreateDistrictVariant(ctx, params)
}
