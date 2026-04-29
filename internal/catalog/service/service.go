package service

// Service provides catalog application use cases.
type Service struct {
	categories                  CategoryWriter
	cities                      CityWriter
	districts                   DistrictWriter
	products                    ProductWriter
	variants                    VariantWriter
	districtVariants            DistrictVariantWriter
	districtVariantPriceUpdater DistrictVariantPriceUpdater
	productImageUpdater         ProductImageUpdater
	variantImageUpdater         VariantImageUpdater
}

// New constructs catalog application service.
func New(
	categories CategoryWriter,
	cities CityWriter,
	districts DistrictWriter,
	products ProductWriter,
	variants VariantWriter,
	districtVariants DistrictVariantWriter,
	districtVariantPriceUpdater DistrictVariantPriceUpdater,
	productImageUpdater ProductImageUpdater,
	variantImageUpdater VariantImageUpdater,
) *Service {
	return &Service{
		categories:                  categories,
		cities:                      cities,
		districts:                   districts,
		products:                    products,
		variants:                    variants,
		districtVariants:            districtVariants,
		districtVariantPriceUpdater: districtVariantPriceUpdater,
		productImageUpdater:         productImageUpdater,
		variantImageUpdater:         variantImageUpdater,
	}
}
