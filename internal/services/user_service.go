package services

import (
	"context"
	"log/slog"

	"botmanager/internal/domain"
)

// UserRepository defines persistence operations
// required by UserService.
type UserRepository interface {
	Save(ctx context.Context, u *domain.User) error
	ByID(ctx context.Context, id int) (*domain.User, error)
}

// UserService orchestrates user-related use cases.
//
// It coordinates:
//   - aggregates loading
//   - domain mutations
//   - persistence
//   - transaction handling
//   - event publishing
//   - structured logging
type UserService struct {
	repo   UserRepository
	tx     TxManager
	bus    EventBus
	logger *slog.Logger
}

// NewUserService creates a new UserService instance.
//
// All dependencies must provided and must not be nil.
func NewUserService(
	repo UserRepository,
	tx TxManager,
	bus EventBus,
	logger *slog.Logger,
) *UserService {
	return &UserService{
		repo:   repo,
		tx:     tx,
		bus:    bus,
		logger: logger,
	}
}

// CreateUser creates a new application user.
func (s *UserService) CreateUser(
	ctx context.Context,
	params domain.NewUserParams,
) (*domain.User, error) {
	var user *domain.User

	err := s.tx.WithinTransaction(ctx, func(ctx context.Context) error {
		u, err := domain.NewUser(params)
		if err != nil {
			s.logger.Warn("failed to create user", "err", err)
			return err
		}

		if err := s.repo.Save(ctx, u); err != nil {
			s.logger.Error("failed to save user", "err", err)
			return err
		}

		user = u
		return nil
	})

	return user, err
}
