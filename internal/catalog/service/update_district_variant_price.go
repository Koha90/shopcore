package service

import (
	"context"
	"fmt"
)

// DistrictVariantPriceUpdater updates price for an existing district variant placement.
type DistrictVariantPriceUpdater interface {
	UpdateDistrictVariantPrice(ctx context.Context, params UpdateDistrictVariantPriceParams) error
}

// UpdateDistrictVariantPriceParams contains input for district variant price update use case.
type UpdateDistrictVariantPriceParams struct {
	DistrictID int
	VariantID  int
	Price      int
}

// UpdateDistrictVariantPrice updates price for an existing variant placement.
func (s *Service) UpdateDistrictVariantPrice(ctx context.Context, params UpdateDistrictVariantPriceParams) error {
	if s == nil {
		return fmt.Errorf("catlog service is nil")
	}
	if s.districtVariantPriceUpdater == nil {
		return fmt.Errorf("catalog district variant price updater is nil")
	}

	switch {
	case params.DistrictID <= 0:
		return ErrDistrictVariantDistrictIDInvalid

	case params.VariantID <= 0:
		return ErrDistrictVariantVariantIDInvalid

	case params.Price <= 0:
		return ErrDistrictVariantPriceInvalid
	}

	return s.districtVariantPriceUpdater.UpdateDistrictVariantPrice(ctx, params)
}
