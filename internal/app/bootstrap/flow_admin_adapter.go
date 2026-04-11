package bootstrap

import (
	"context"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
	"github.com/koha90/shopcore/internal/flow"
)

// flowCatalogAdminAdapter adapts catalog application service to flow admin ports.
type flowCatalogAdminAdapter struct {
	svc        *catalogservice.Service
	cities     flow.CityLister
	categories flow.CategoryLister
	districts  flow.DistrictLister
	products   flow.ProductLister
}

func newFlowCatalogAdminAdapter(
	svc *catalogservice.Service,
	cities flow.CityLister,
	categories flow.CategoryLister,
	districts flow.DistrictLister,
	products flow.ProductLister,
) *flowCatalogAdminAdapter {
	if svc == nil && cities == nil && categories == nil {
		return nil
	}

	return &flowCatalogAdminAdapter{
		svc:        svc,
		cities:     cities,
		categories: categories,
		districts:  districts,
		products:   products,
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

func (a *flowCatalogAdminAdapter) CreateProduct(ctx context.Context, params flow.CreateProductParams) error {
	return a.svc.CreateProduct(ctx, catalogservice.CreateProductParams{
		CategoryID: params.CategoryID,
		Code:       params.Code,
		Name:       params.Name,
	})
}

func (a *flowCatalogAdminAdapter) CreateDistrict(ctx context.Context, params flow.CreateDistrictParams) error {
	return a.svc.CreateDistrict(ctx, catalogservice.CreateDistrictParams{
		CityID: params.CityID,
		Code:   params.Code,
		Name:   params.Name,
	})
}

func (a *flowCatalogAdminAdapter) ListCities(ctx context.Context) ([]flow.CityListItem, error) {
	if a == nil || a.cities == nil {
		return nil, nil
	}

	return a.cities.ListCities(ctx)
}

func (a *flowCatalogAdminAdapter) ListCategories(ctx context.Context) ([]flow.CategoryListItem, error) {
	if a == nil || a.categories == nil {
		return nil, nil
	}

	return a.categories.ListCategories(ctx)
}

func (a *flowCatalogAdminAdapter) ListDistricts(ctx context.Context) ([]flow.DistrictListItem, error) {
	if a == nil || a.products == nil {
		return nil, nil
	}

	return a.districts.ListDistricts(ctx)
}

func (a *flowCatalogAdminAdapter) ListProducts(ctx context.Context) ([]flow.ProductListItem, error) {
	if a == nil || a.products == nil {
		return nil, nil
	}

	return a.products.ListProducts(ctx)
}

func (a *flowCatalogAdminAdapter) CreateVariant(ctx context.Context, params flow.CreateVariantParams) error {
	return a.svc.CreateVariant(ctx, catalogservice.CreateVariantParams{
		ProductID: params.ProductID,
		Code:      params.Code,
		Name:      params.Name,
	})
}
