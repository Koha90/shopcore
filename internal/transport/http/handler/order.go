// Package handler contains HTTP handlers.
//
// Handlers translate HTTP requests into service calls
// and map service responses to HTTP responses.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/koha90/shopcore/internal/domain"
	"github.com/koha90/shopcore/internal/service"
	"github.com/koha90/shopcore/internal/transport/http/dto"
)

// OrderHandler handles HTTP requests related to orders.
type OrderHandler struct {
	service *service.OrderService
}

// NewOrderHandler creates new OrderHandler instance.
func NewOrderHandler(s *service.OrderService) *OrderHandler {
	return &OrderHandler{service: s}
}

// Create creates a new order for selected product variant.
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	order, err := h.service.CreateForVariant(
		r.Context(),
		req.CustomerID,
		req.ProductID,
		req.VariantID,
	)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.OrderReponse{
		ID:     order.ID(),
		UserID: order.UserID(),
		Status: string(order.Status()),
		Total:  order.Total(),
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ConfirmPayment confirms external payment for an existing order.
func (h *OrderHandler) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	orderID, err := parseOrderID(r)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	err = h.service.ConfirmPayment(r.Context(), orderID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrOrderNotFound):
			http.Error(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, domain.ErrOrderAlreadyPaid),
			errors.Is(err, domain.ErrOrderAlreadyCancelled),
			errors.Is(err, domain.ErrOrderNotPending):
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Cancel handles order cancellation.
func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.Cancel(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseOrderID(r *http.Request) (int, error) {
	orderIDStr := chi.URLParam(r, "id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil {
		return 0, err
	}

	return orderID, err
}
