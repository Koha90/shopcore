package service

import (
	"context"
	"fmt"
	"strings"
)

// VariantWriter stores variant data.
type VariantWriter interface {
	CreateVariant(ctx context.Context, params CreateVariantParams) error
}

// CreateVariantParams contains data required to create one variant.
type CreateVariantParams struct {
	ProductID   int
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
}

// CreateVariant validates input and stores a new variant.
func (s *Service) CreateVariant(ctx context.Context, params CreateVariantParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.variants == nil {
		return fmt.Errorf("variant writer is nil")
	}

	params.Code = normalizeCode(params.Code)
	params.Name = strings.TrimSpace(params.Name)
	params.NameLatin = strings.TrimSpace(params.NameLatin)
	params.Description = strings.TrimSpace(params.Description)

	switch {
	case params.ProductID <= 0:
		return ErrVariantProductIDInvalid
	case params.Code == "":
		return ErrVariantCodeEmpty
	case params.Name == "":
		return ErrVariantNameEmpty
	}

	return s.variants.CreateVariant(ctx, params)
}
