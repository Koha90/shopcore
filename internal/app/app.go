// Package app composes and wires application dependecies.
//
// It initializes infrastructure, services and transport layer
// and builds application runtime container.
package app

import (
	"context"
	"net/http"

	"github.com/koha90/shopcore/internal/config"
	transporthttp "github.com/koha90/shopcore/internal/transport/http"
	"github.com/koha90/shopcore/internal/transport/http/handler"
	"github.com/koha90/shopcore/pkg/logger"
)

type App struct {
	server *http.Server
}

// New creates application container with configured dependecies.
//
// It returns an error if dependecy wiring fails.
func New(cfg *config.Config) (*App, error) {
	logg, err := logger.Setup(cfg.Env)
	if err != nil {
		return nil, err
	}

	orderService, err := BuildOrderService(context.Background(), logg.Logger)
	if err != nil {
		return nil, err
	}

	orderHandler := handler.NewOrderHandler(orderService)
	router := transporthttp.NewRouter(orderHandler)

	server := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: router,
	}

	return &App{server: server}, nil
}

// Run starts HTTP server.
func (a *App) Run() error {
	return a.server.ListenAndServe()
}
