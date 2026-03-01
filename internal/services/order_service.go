// Package services contain application use cases.
//
// It coordinates domain logic, repositories and transactions.
// It does not contain business rules.
package services

import (
	"context"
	"log/slog"
	"time"

	"botmanager/internal/domain"
)

// OrderRepository defines persistence contain for Order aggregate.
type OrderRepository interface {
	Save(ctx context.Context, order *domain.Order) error
	ByID(ctx context.Context, id int) (*domain.Order, error)
}

// OrderService orchestrates order use cases.
type OrderService struct {
	repo     OrderRepository
	userRepo UserRepository
	bus      EventBus
	tx       TxManager
	logger   *slog.Logger
}

// NewOrderService creates OrderService instance.
//
// logger must not be nil.
func NewOrderService(
	repo OrderRepository,
	userRepo UserRepository,
	bus EventBus,
	tx TxManager,
	logger *slog.Logger,
) *OrderService {
	return &OrderService{
		repo:     repo,
		userRepo: userRepo,
		bus:      bus,
		tx:       tx,
		logger:   logger,
	}
}

// ConfirmPayment marks order as paid and publishes domain events.
func (s *OrderService) ConfirmPayment(
	ctx context.Context,
	orderID int,
) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		order, err := s.repo.ByID(ctx, orderID)
		if err != nil {
			s.logger.Error(
				"failed to find order",
				"orderID", orderID,
				"err", err,
			)
			return err
		}

		if err := order.MarkPaid(time.Now()); err != nil {
			s.logger.Error(
				"failes to mark order as paid",
				"orderID", orderID,
				"err", err,
			)
			return err
		}

		if err := s.repo.Save(ctx, order); err != nil {
			s.logger.Error(
				"failed to save order",
				"err", err,
			)
			return err
		}

		events := order.PullEvents()

		if err := s.bus.Publish(ctx, events...); err != nil {
			s.logger.Error(
				"failed to publishe order events",
				"orderID", orderID,
				"err", err,
			)
			return err
		}

		s.logger.Info(
			"order successfully paid",
			"orderID", orderID,
		)
		return nil
	})
}

// ConfirmOrder confirms order and deduct user balance.
// Operation is executed atomically.
func (s *OrderService) ConfirmOrder(
	ctx context.Context,
	orderID int,
) error {
	return s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		order, err := s.repo.ByID(ctx, orderID)
		if err != nil {
			s.logger.Error(
				"faled to load order",
				"orderID", orderID,
				"err", err,
			)
			return err
		}

		user, err := s.userRepo.ByID(ctx, order.UserID())
		if err != nil {
			s.logger.Error(
				"failed to load user",
				"userID", order.UserID(),
				"err", err,
			)
			return err
		}

		if err := user.DeductBalance(order.Total()); err != nil {
			s.logger.Warn(
				"insufficient user balance",
				"userID", user.ID,
				"amount", order.Total(),
				"err", err,
			)
		}

		if err := order.MarkPaid(time.Now()); err != nil {
			s.logger.Warn(
				"failed to mark order as paid",
				"orderID", orderID,
				"err", err,
			)
			return err
		}

		if err := s.userRepo.Save(ctx, user); err != nil {
			return err
		}

		if err := s.repo.Save(ctx, order); err != nil {
			return err
		}

		s.logger.Info(
			"order confirmed successfully",
			"orderID", orderID,
			"userID", user.ID(),
			"amount", order.Total(),
		)

		return nil
	})
}
