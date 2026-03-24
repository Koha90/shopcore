// Package http provides HTTP transport layer.
//
// It defines routing, middleware and conencts HTTP handlers
// with application services.
package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/koha90/shopcore/internal/transport/http/handler"
)

// NewRouter configurates and returns HTTP router.
//
// It defines API routes, groups and middleware.
// The router is responsible only for HTTP concerns.
func NewRouter(orderHandler *handler.OrderHandler) http.Handler {
	r := chi.NewRouter()

	// ---- Global middleware ----
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ---- API grouping ----
	r.Route("/api", func(r chi.Router) {
		// Versioning group
		r.Route("/v1", func(r chi.Router) {
			// Orders endpoints
			r.Route("/orders", func(r chi.Router) {
				// POST /api/v1/orders
				r.Post("/", orderHandler.Create)
				// POST /api/v1/orders/{id}/confirm
				r.Post("/{id}/confirm", orderHandler.ConfirmPayment)
				// POST /api/v1/orders/{id}/cancel
				r.Post("/{id}/cancel", orderHandler.Cancel)
			})
		})
	})

	return r
}
