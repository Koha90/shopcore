package service

import (
	"context"
	"log/slog"

	"github.com/koha90/shopcore/internal/domain"
)

type ProductService struct {
	repo   ProductRepository
	tx     TxManager
	bus    EventBus
	logger *slog.Logger
}

// NewProductService creates a new ProductService instance.
//
// logger may be nil, in that case slog.Default() is used.
func NewProductService(
	repo ProductRepository,
	tx TxManager,
	bus EventBus,
	logger *slog.Logger,
) *ProductService {
	if repo == nil {
		panic("service: ProductRepository is nil")
	}

	if tx == nil {
		panic("service: TxManager is nil")
	}

	if bus == nil {
		panic("service: EventBus is nil")
	}

	if logger == nil {
		logger = slog.Default()
	}
	return &ProductService{
		repo:   repo,
		tx:     tx,
		bus:    bus,
		logger: logger,
	}
}

// Create creates a new product.
func (s *ProductService) Create(
	ctx context.Context,
	name string,
	categoryID int,
	description string,
	imagePath string,
) (*domain.Product, error) {
	var product *domain.Product
	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		p, err := domain.NewProduct(name, categoryID, description, imagePath)
		if err != nil {
			s.logger.Warn("failed to create product", "err", err)
			return err
		}

		if err := s.repo.Save(ctx, p); err != nil {
			s.logger.Error("failed to save product", "err", err)
			return err
		}

		product = p
		return nil
	})
	if err != nil {
		return nil, err
	}

	return product, nil
}

// AddVariant adds a new variant to an existing product.
func (s *ProductService) AddVariant(
	ctx context.Context,
	productID int,
	packSize string,
	districtID int,
	price int64,
) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		product, err := s.repo.ByID(ctx, productID)
		if err != nil {
			s.logger.Error(
				"failed to load product",
				"id", productID,
				"err", err,
			)
			return err
		}

		if err := product.AddVariant(packSize, districtID, price); err != nil {
			s.logger.Warn(
				"failed to add variant",
				"product_id", productID,
				"pack_size", packSize,
				"district_id", districtID,
				"price", price,
				"err", err)
			return err
		}

		if err := s.repo.Save(ctx, product); err != nil {
			s.logger.Error(
				"failed to save product",
				"product_id", productID,
				"err", err,
			)
			return err
		}

		s.logger.Info(
			"variant added",
			"product_id", productID,
			"packSize", packSize,
			"district_id", districtID,
			"price", price,
		)
		return nil
	})
}
