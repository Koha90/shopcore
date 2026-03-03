// Package service contain application use cases.
//
// It coordinates domain logic, repositories and transactions.
// It does not contain business rules.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"botmanager/internal/domain"
)

type OrderService struct {
	products productReader
	orders   orderRepository

	tx     TxManager
	bus    EventBus
	logger *slog.Logger
}

func NewOrderService(
	products productReader,
	orders orderRepository,

	tx TxManager,
	bus EventBus,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		products: products,
		orders:   orders,
		idGen:    idGen,
		tx:       tx,
		bus:      bus,
		logger:   logger,
	}
}

func (s *OrderService) Create(
	ctx context.Context,
	customerID int,
	productID int,
	variantID int,
) (*domain.Order, error) {
	var created *domain.Order

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		product, err := s.products.ByID(ctx, productID)
		if err != nil {
			return err
		}

		variant, err := product.VariantByID(variantID)
		if err != nil {
			return err
		}

		order := domain.NewOrder(
			customerID,
			productID,
			variant.ID(),
			variant.Price(),
		)

		if err := s.orders.Create(ctx, order); err != nil {
			return err
		}

		created = order
		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *OrderService) Confirm(ctx context.Context, id int) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info("confirming order", "order_id", id)

		order, err := s.orders.ByID(ctx, id)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				return domain.ErrOrderNotFound
			}
			return fmt.Errorf("get order: %w", err)
		}

		if err := order.Confirm(); err != nil {
			return err
		}

		if err := s.orders.Update(ctx, order); err != nil {
			s.logger.Error("failed to update order", "error", err)
			return domain.ErrOrderUpdate
		}

		events := order.PullEvents()
		if len(events) > 0 {
			if err := s.bus.Publish(ctx, events...); err != nil {
				return fmt.Errorf("publish events: %w", err)
			}
		}

		s.logger.Info("order confirmed successfully", "order_id", id)

		return nil
	})
}

func (s *OrderService) Cancel(ctx context.Context, id int) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info("cancelling order", "order_id", id)

		order, err := s.orders.ByID(ctx, id)
		if err != nil {
			s.logger.Error("failed to load order", "error", err)
			return domain.ErrOrderNotFound
		}

		if err := order.Cancel(); err != nil {
			s.logger.Warn("order cancel rejected", "error", err)
			return err
		}

		if err := s.orders.Update(ctx, order); err != nil {
			s.logger.Error("failed to update order", "error", err)
			return err
		}

		events := order.PullEvents()
		err = s.bus.Publish(ctx, events...)
		if err != nil {
			s.logger.Error("failed publish event cancel", "err", err)
		}

		s.logger.Info("order canceled successfully", "order_id", id)

		return nil
	})
}
