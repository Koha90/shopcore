package service

import (
	"context"
	"fmt"
	"strings"
)

// ProductWriter stores product data.
type ProductWriter interface {
	CreateProduct(ctx context.Context, params CreateProductParams) error
}

// CreateProductParams contains data required to create one product.
type CreateProductParams struct {
	CategoryID  int
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
}

// CreateProduct validates input and stores a new product.
func (s *Service) CreateProduct(ctx context.Context, params CreateProductParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.products == nil {
		return fmt.Errorf("product writer is nil")
	}

	params.Code = normalizeCode(params.Code)
	params.Name = strings.TrimSpace(params.Name)
	params.NameLatin = strings.TrimSpace(params.NameLatin)
	params.Description = strings.TrimSpace(params.Description)

	switch {
	case params.CategoryID <= 0:
		return ErrProductCategoryIDInvalid
	case params.Code == "":
		return ErrProductCodeEmpty
	case params.Name == "":
		return ErrProductNameEmpty
	}

	return s.products.CreateProduct(ctx, params)
}
