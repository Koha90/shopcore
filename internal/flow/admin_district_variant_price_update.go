package flow

import "context"

// DistrictVariantPriceUpdater updates price for an existing district variant placement
type DistrictVariantPriceUpdater interface {
	UpdateDistrictVariantPrice(ctx context.Context, params UpdateDistictVariantPriceParams) error
}

// UpdateDistictVariantPriceParams contains input for district variant price update.
type UpdateDistictVariantPriceParams struct {
	DistrictID int
	VariantID  int
	Price      int
}
