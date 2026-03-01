package services

import (
	"context"
	"log/slog"

	"botmanager/internal/domain"
)

// ProductVariantRepository defines persistence operations
// required by ProductVariantService.
type ProductVariantRepository interface {
	Save(ctx context.Context, v *domain.ProductVariant) error
	ByID(ctx context.Context, id int) (*domain.ProductVariant, error)
}

// ProductVariantService orchestrates product variant use cases.
//
// It coordinates:
//   - aggregate loading
//   - domain mutations
//   - persistence
//   - event publishing
//   - loading
type ProductVariantService struct {
	repo   ProductVariantRepository
	events EventPublisher

	logger *slog.Logger
}

// NewProductVariantService creates a new instance
// of ProductVariantRepository.
func NewProductVariantService(
	repo ProductVariantRepository,
	events EventPublisher,
	logger *slog.Logger,
) *ProductVariantService {
	return &ProductVariantService{
		repo:   repo,
		logger: logger,
		events: events,
	}
}

// CreateVariant creates a new product variant
// and persists it.
func (s *ProductVariantService) CreateVariant(
	ctx context.Context,
	packSize string,
	districtID int,
	price int64,
) error {
	variant, err := domain.NewProductVariant(packSize, districtID, price)
	if err != nil {
		s.logger.Error("failed to create product variant", "err", err)
		return err
	}

	if err := s.repo.Save(ctx, variant); err != nil {
		s.logger.Error("failed to persist product variant", "err", err)
		return err
	}

	s.logger.Info("product variant created", "variant_id", variant.ID())
	return nil
}

// ChangePrice updates variant price.
func (s *ProductVariantService) ChangePrice(
	ctx context.Context,
	id int,
	newPrice int64,
) error {
	variant, err := s.repo.ByID(ctx, id)
	if err != nil {
		s.logger.Error("failed to find variant by id",
			"variant_id", id,
			"err", err)
		return err
	}

	if err := variant.ChangePrice(newPrice); err != nil {
		s.logger.Error("failed to change price", "err", err)
		return err
	}

	if err := s.repo.Save(ctx, variant); err != nil {
		s.logger.Error("failed to save product variant",
			"variant_id", id,
			"err", err)
		return err
	}

	s.logger.Info("product variant price changed", "variant_id", id, "new_price", newPrice)
	return nil
}
