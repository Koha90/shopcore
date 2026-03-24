package app

import (
	"context"
	"log/slog"
	"sync"

	"github.com/koha90/shopcore/internal/domain"
	"github.com/koha90/shopcore/internal/infrastructure/eventbus"
	"github.com/koha90/shopcore/internal/service"
	"github.com/koha90/shopcore/internal/storage/memory"
)

// BuildOrderService creates in-memory dependencies for local development
// and returns configured order service instance.
func BuildOrderService(ctx context.Context, logger *slog.Logger) (*service.OrderService, error) {
	mu := &sync.Mutex{}
	// repositories
	orderRepo := memory.NewOrderRepository(mu)
	productRepo := memory.NewProductRepository(mu)
	userRepo := memory.NewUserRepository(mu)
	txManager := memory.NewTxManager(mu)
	bus := eventbus.New(logger)

	product, err := domain.NewProduct("Amnesia", 1, "good stuff", "/tmp/img.png")
	if err != nil {
		return nil, err
	}

	if err := productRepo.Save(ctx, product); err != nil {
		return nil, err
	}

	orderService := service.NewOrderService(
		productRepo,
		orderRepo,
		userRepo,
		bus,
		txManager,
		logger,
	)

	return orderService, nil
}
