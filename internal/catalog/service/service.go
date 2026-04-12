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
) *Service {
	return &Service{
		categories:                  categories,
		cities:                      cities,
		districts:                   districts,
		products:                    products,
		variants:                    variants,
		districtVariants:            districtVariants,
		districtVariantPriceUpdater: districtVariantPriceUpdater,
	}
}
