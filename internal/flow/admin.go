package flow

import "context"

// CategoryCreator defines the admin write use case required by flow.
//
// Flow depends on this narrow port to create catalog categories without
// depending on storage or transport-specific details.
type CategoryCreator interface {
	CreateCategory(ctx context.Context, params CreateCategoryParams) error
}

// CreateCategoryParams contains data required by flow admin category creation.
//
// This is a flow-local model passed through CategoryCreator.
type CreateCategoryParams struct {
	Code string
	Name string
}

// CityCreator defines the admin write use case required by flow.
type CityCreator interface {
	CreateCity(ctx context.Context, params CreateCityParams) error
}

// CreateCityParams contains data required by flow admin city creation.
type CreateCityParams struct {
	Code string
	Name string
}

// CityListItem contains one city option for admin selection flows.
type CityListItem struct {
	ID    int
	Code  string
	Label string
}

// CityLister defines the admin read use case required by flow
// to select an existing city before nested catalog creation steps.
type CityLister interface {
	ListCities(ctx context.Context) ([]CityListItem, error)
}
