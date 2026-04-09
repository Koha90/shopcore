package service

// Service provides catalog application use cases.
type Service struct {
	categories CategoryWriter
	cities     CityWriter
	districts  DistrictWriter
}

// New constructs catalog application service.
func New(
	categories CategoryWriter,
	cities CityWriter,
	districts DistrictWriter,
) *Service {
	return &Service{
		categories: categories,
		cities:     cities,
		districts:  districts,
	}
}
