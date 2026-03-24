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
	"time"

	"github.com/koha90/shopcore/internal/domain"
)

// OrderService orchestrates order-related use cases.
type OrderService struct {
	products ProductReader
	orders   OrderRepository
	users    UserRepository

	bus    EventBus
	tx     TxManager
	logger *slog.Logger
}

// NewOrderService creates a new OrderService instance.
//
// logger may be nil, in that case slog.Default() is used.
func NewOrderService(
	products ProductReader,
	orders OrderRepository,
	users UserRepository,
	bus EventBus,
	tx TxManager,
	logger *slog.Logger,
) *OrderService {
	if products == nil {
		panic("service: ProductReader is nil")
	}

	if orders == nil {
		panic("service: OrderRepository is nil")
	}

	if users == nil {
		panic("service: UserRepository is nil")
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

	return &OrderService{
		products: products,
		orders:   orders,
		users:    users,
		bus:      bus,
		tx:       tx,
		logger:   logger,
	}
}

// CreateForVariant creates a new order for a selected product variant.
func (s *OrderService) CreateForVariant(
	ctx context.Context,
	userID int,
	productID int,
	variantID int,
) (*domain.Order, error) {
	var created *domain.Order

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info(
			"creating order",
			"user_id", userID,
			"product_id", productID,
			"variant_id", variantID,
		)

		product, err := s.products.ByID(ctx, productID)
		if err != nil {
			s.logger.Error(
				"failed to load product",
				"product_id", productID,
				"variant_id", variantID,
				"err", err,
			)
			return fmt.Errorf("load product: %w", err)
		}

		variant, err := product.VariantByID(variantID)
		if err != nil {
			s.logger.Warn(
				"failed to load product variant",
				"product_id", productID,
				"variant_id", variantID,
				"err", err,
			)
			return fmt.Errorf("load product variant: %w", err)
		}

		items := []domain.OrderItem{
			domain.NewOrderItem(product.ID(), variant.ID(), 1, variant.Price()),
		}

		order, err := domain.NewOrder(userID, items, time.Now())
		if err != nil {
			return err
		}

		if err := s.orders.Save(ctx, order); err != nil {
			s.logger.Error(
				"failed to create order",
				"user_id", userID,
				"product_id", productID,
				"variant_id", variantID,
				"err", err,
			)
			return err
		}

		created = order

		s.logger.Info(
			"order created successfully",
			"order_id", order.ID(),
			"user_id", userID,
		)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return created, nil
}

// ConfirmPayment marks order as paid after external payment confirmation.
//
// Use this method when payment was completed outside of internal balance
// workflow, for example via card, crypto, SBP or another payment gateway.
func (s *OrderService) ConfirmPayment(
	ctx context.Context,
	orderID int,
) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info("confirming external payment", "order_id", orderID)

		order, err := s.orders.ByID(ctx, orderID)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				return domain.ErrOrderNotFound
			}

			s.logger.Error(
				"failed to load order",
				"order_id", orderID,
				"err", err,
			)
			return fmt.Errorf("load order: %w", err)
		}

		if err := order.MarkPaid(time.Now()); err != nil {
			s.logger.Warn(
				"failed to mark order as paid",
				"order_id", orderID,
				"err", err,
			)
			return err
		}

		if err := s.orders.Save(ctx, order); err != nil {
			s.logger.Error(
				"filed to update order",
				"order_id", orderID,
				"err", err,
			)
			return domain.ErrOrderUpdate
		}

		events := order.PullEvents()
		if len(events) > 0 {
			if err := s.bus.Publish(ctx, events...); err != nil {
				s.logger.Error(
					"failed to publish order events",
					"order_id", orderID,
					"err", err,
				)
				return fmt.Errorf("publish events: %w", err)
			}
		}

		s.logger.Info("external payment confirmed successfully",
			"order_id", orderID,
		)

		return nil
	})
}

// PayFromBalance deducts user balance and marks order as paid.
//
// Use this method when payment is performed with internal user balance.
func (s *OrderService) PayFromBalance(
	ctx context.Context,
	orderID int,
) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info("paying order from balance", "order_id", orderID)

		order, err := s.orders.ByID(ctx, orderID)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				return domain.ErrOrderNotFound
			}

			s.logger.Error(
				"failed to load order",
				"order_id", orderID,
				"err", err,
			)
			return fmt.Errorf("load order: %w", err)
		}

		user, err := s.users.ByID(ctx, order.UserID())
		if err != nil {
			s.logger.Error(
				"failed to load user",
				"user_id", order.UserID(),
				"order_id", orderID,
				"err", err,
			)
			return fmt.Errorf("load user: %w", err)
		}

		if err := user.DeductBalance(order.Total()); err != nil {
			s.logger.Warn(
				"failed to deduct user balance",
				"user_id", order.UserID(),
				"order_id", orderID,
				"amount", order.Total(),
				"err", err,
			)
			return err
		}

		if err := order.MarkPaid(time.Now()); err != nil {
			s.logger.Warn(
				"failed to mark order as paid",
				"order_id", orderID,
				"err", err,
			)
			return err
		}

		if err := s.users.Save(ctx, user); err != nil {
			s.logger.Error(
				"failed to save user",
				"user_id", user.ID(),
				"err", err,
			)
			return fmt.Errorf("save user: %w", err)
		}

		if err := s.orders.Save(ctx, order); err != nil {
			s.logger.Error(
				"failed to update order",
				"order_id", orderID,
				"err", err,
			)
			return domain.ErrOrderUpdate
		}

		events := order.PullEvents()
		if len(events) > 0 {
			if err := s.bus.Publish(ctx, events...); err != nil {
				s.logger.Error(
					"failed to publish order events",
					"order_id", orderID,
					"err", err,
				)
				return fmt.Errorf("publish events: %w", err)
			}
		}

		s.logger.Info(
			"order paid from balance successfully",
			"order_id", orderID,
			"user_id", user.ID(),
			"amount", order.Total(),
		)

		return nil
	})
}

// Cancel cancels an existing order and publishes domain events.
func (s *OrderService) Cancel(ctx context.Context, orderID int) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		s.logger.Info("cancelling order", "order_id", orderID)

		order, err := s.orders.ByID(ctx, orderID)
		if err != nil {
			if errors.Is(err, domain.ErrOrderNotFound) {
				return domain.ErrOrderNotFound
			}

			s.logger.Error(
				"failed to load order",
				"order_id", orderID,
				"err", err,
			)
			return fmt.Errorf("load order: %w", err)
		}

		if err := order.Cancel(time.Now()); err != nil {
			s.logger.Warn(
				"failed to cancel order",
				"order_id", orderID,
				"err", err,
			)
			return err
		}

		if err := s.orders.Save(ctx, order); err != nil {
			s.logger.Error(
				"failed to save order",
				"order_id", orderID,
				"err", err,
			)
			return domain.ErrOrderSave
		}

		events := order.PullEvents()
		if len(events) > 0 {
			if err := s.bus.Publish(ctx, events...); err != nil {
				s.logger.Error(
					"failed to publish order events",
					"order_id", orderID,
					"err", err,
				)
				return fmt.Errorf("publish events: %w", err)
			}
		}

		s.logger.Info(
			"order cancelled successfully",
			"order_id", orderID,
		)

		return nil
	})
}
