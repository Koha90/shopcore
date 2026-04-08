package service

import (
	"context"
	"fmt"
	"strings"
)

// CityWriter stores city data.
type CityWriter interface {
	CreateCity(ctx context.Context, params CreateCityParams) error
}

// CreateCityParams contains data required to create one city.
type CreateCityParams struct {
	Code      string
	Name      string
	NameLatin string
	SortOrder int
}

// CreateCity validate input and stores a new city.
func (s *Service) CreateCity(ctx context.Context, params CreateCityParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.cities == nil {
		return fmt.Errorf("city writer is nil")
	}

	params.Code = normalizeCode(params.Code)
	params.Name = strings.TrimSpace(params.Name)
	params.NameLatin = strings.TrimSpace(params.NameLatin)

	switch {
	case params.Code == "":
		return ErrCityCodeEmpty
	case params.Name == "":
		return ErrCityNameEmpty
	}

	return s.cities.CreateCity(ctx, params)
}
