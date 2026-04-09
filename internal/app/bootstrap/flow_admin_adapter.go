package bootstrap

import (
	"context"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
	"github.com/koha90/shopcore/internal/flow"
)

// flowCatalogAdminAdapter adapts catalog application service to flow admin ports.
type flowCatalogAdminAdapter struct {
	svc    *catalogservice.Service
	cities flow.CityLister
}

func newFlowCatalogAdminAdapter(
	svc *catalogservice.Service,
	cities flow.CityLister,
) *flowCatalogAdminAdapter {
	if svc == nil && cities == nil {
		return nil
	}

	return &flowCatalogAdminAdapter{
		svc:    svc,
		cities: cities,
	}
}

func (a *flowCatalogAdminAdapter) CreateCategory(ctx context.Context, params flow.CreateCategoryParams) error {
	return a.svc.CreateCategory(ctx, catalogservice.CreateCategoryParams{
		Code: params.Code,
		Name: params.Name,
	})
}

func (a *flowCatalogAdminAdapter) CreateCity(ctx context.Context, params flow.CreateCityParams) error {
	return a.svc.CreateCity(ctx, catalogservice.CreateCityParams{
		Code: params.Code,
		Name: params.Name,
	})
}

func (a *flowCatalogAdminAdapter) ListCities(ctx context.Context) ([]flow.CityListItem, error) {
	if a == nil || a.cities == nil {
		return nil, nil
	}

	return a.cities.ListCities(ctx)
}
