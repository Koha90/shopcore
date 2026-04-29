package service

import (
	"context"
	"fmt"
	"strings"
)

// ProductImageUpdater updates image source for an existing product.
type ProductImageUpdater interface {
	UpdateProductImage(ctx context.Context, params UpdateProductImageParams) error
}

// VariantImageUpdater updates image source for an existing variant.
type VariantImageUpdater interface {
	UpdateVariantImage(ctx context.Context, params UpdateVariantImageParams) error
}

// UpdateProductImageParams contains input for product image update.
type UpdateProductImageParams struct {
	ProductID int
	ImageURL  string
}

// UpdateVariantImageParams contains input for variant image update.
type UpdateVariantImageParams struct {
	VariantID int
	ImageURL  string
}

// UpdateProductImage updates product image source.
func (s *Service) UpdateProductImage(ctx context.Context, params UpdateProductImageParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.productImageUpdater == nil {
		return fmt.Errorf("catalog product image updater is nil")
	}

	params.ImageURL = strings.TrimSpace(params.ImageURL)

	switch {
	case params.ProductID <= 0:
		return ErrProductIDInvalid
	case params.ImageURL == "":
		return ErrImageURLInvalid
	}

	return s.productImageUpdater.UpdateProductImage(ctx, params)
}

// UpdateVariantImage updates variant image source.
func (s *Service) UpdateVariantImage(ctx context.Context, params UpdateVariantImageParams) error {
	if s == nil {
		return fmt.Errorf("catalog service is nil")
	}
	if s.variantImageUpdater == nil {
		return fmt.Errorf("catalog variant image updater is nil")
	}

	params.ImageURL = strings.TrimSpace(params.ImageURL)

	switch {
	case params.VariantID <= 0:
		return ErrVariantIDInvalid
	case params.ImageURL == "":
		return ErrImageURLInvalid
	}

	return s.variantImageUpdater.UpdateVariantImage(ctx, params)
}
