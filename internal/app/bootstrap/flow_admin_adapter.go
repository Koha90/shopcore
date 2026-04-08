package bootstrap

import (
	"context"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
	"github.com/koha90/shopcore/internal/flow"
)

// flowCatalogAdminAdapter adapts catalog application service to flow admin ports.
type flowCatalogAdminAdapter struct {
	svc *catalogservice.Service
}

func newFlowCatalogAdminAdapter(svc *catalogservice.Service) *flowCatalogAdminAdapter {
	if svc == nil {
		return nil
	}

	return &flowCatalogAdminAdapter{svc: svc}
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
