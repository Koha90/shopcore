package bootstrap

import (
	"context"
	"errors"

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
	variants   flow.VariantLister
}

func newFlowCatalogAdminAdapter(
	svc *catalogservice.Service,
	cities flow.CityLister,
	categories flow.CategoryLister,
	districts flow.DistrictLister,
	products flow.ProductLister,
	variants flow.VariantLister,
) *flowCatalogAdminAdapter {
	if svc == nil && cities == nil && categories == nil && districts == nil && products == nil && variants == nil {
		return nil
	}

	return &flowCatalogAdminAdapter{
		svc:        svc,
		cities:     cities,
		categories: categories,
		districts:  districts,
		products:   products,
		variants:   variants,
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
	if a == nil || a.districts == nil {
		return nil, nil
	}

	return a.districts.ListDistricts(ctx)
}

func (a *flowCatalogAdminAdapter) ListProducts(ctx context.Context) ([]flow.ProductListItem, error) {
	if a == nil {
		return nil, errors.New("flow catalog admin adapter is nil")
	}
	if a.products == nil {
		return nil, errors.New("flow product lister is nil inside admin adapter")
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

func (a *flowCatalogAdminAdapter) ListVariants(ctx context.Context) ([]flow.VariantListItem, error) {
	if a == nil {
		return nil, errors.New("flow catalog admin adapter is nil")
	}
	if a.variants == nil {
		return nil, errors.New("flow variant lister is nil inside admin adapter")
	}

	return a.variants.ListVariants(ctx)
}

func (a *flowCatalogAdminAdapter) CreateDistrictVariant(ctx context.Context, params flow.CreateDistrictVariantParams) error {
	return a.svc.CreateDistrictVariant(ctx, catalogservice.CreateDistrictVariantParams{
		DistrictID: params.DistrictID,
		VariantID:  params.VariantID,
		Price:      params.Price,
	})
}

func (a *flowCatalogAdminAdapter) UpdateDistrictVariantPrice(
	ctx context.Context,
	params flow.UpdateDistrictVariantPriceParams,
) error {
	return a.svc.UpdateDistrictVariantPrice(ctx, catalogservice.UpdateDistrictVariantPriceParams{
		DistrictID: params.DistrictID,
		VariantID:  params.VariantID,
		Price:      params.Price,
	})
}
