package bootstrap

import (
	"context"

	catalogservice "github.com/koha90/shopcore/internal/catalog/service"
	"github.com/koha90/shopcore/internal/flow"
)

// flowCategoryCreator adapts catalog application service to flow admin port.
type flowCategoryCreator struct {
	svc *catalogservice.Service
}

func newFlowCategoryCreator(svc *catalogservice.Service) flow.CategoryCreator {
	if svc == nil {
		return nil
	}

	return &flowCategoryCreator{svc: svc}
}

func (c *flowCategoryCreator) CreateCategory(ctx context.Context, params flow.CreateCategoryParams) error {
	return c.svc.CreateCategory(ctx, catalogservice.CreateCategoryParams{
		Code: params.Code,
		Name: params.Name,
	})
}
