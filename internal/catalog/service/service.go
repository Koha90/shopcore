package service

import (
	"context"
	"fmt"
	"strings"
)

// CategoryWriter stores category data.
type CategoryWriter interface {
	CreateCategory(ctx context.Context, params CreateCategoryParams) error
}

// Service provides catalog application use cases.
type Service struct {
	categories CategoryWriter
}

// New constructs catalog application service.
func New(categories CategoryWriter) *Service {
	return &Service{
		categories: categories,
	}
}

// CreateCategoryParams contains data required to create one catalog category.
type CreateCategoryParams struct {
	Code        string
	Name        string
	NameLatin   string
	Description string
	SortOrder   int
	// IsActive    bool
}

// CreateCategory validates input and stores a new catalog category.
func (s *Service) CreateCategory(ctx context.Context, params CreateCategoryParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.categories == nil {
		return fmt.Errorf("category writer is nil")
	}

	params.Code = normalizeCode(params.Code)
	params.Name = strings.TrimSpace(params.Name)
	params.NameLatin = strings.TrimSpace(params.NameLatin)
	params.Description = strings.TrimSpace(params.Description)

	switch {
	case params.Code == "":
		return ErrCategoryCodeEmpty
	case params.Name == "":
		return ErrCategoryNameEmpty
	}

	return s.categories.CreateCategory(ctx, params)
}
