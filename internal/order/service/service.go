package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRepositoryNil      = errors.New("order repository is nil")
	ErrBotIDEmpty         = errors.New("order bot id is required")
	ErrChatIDEmpty        = errors.New("order chat id is required")
	ErrUserIDEmpty        = errors.New("order user id is required")
	ErrCityIDEmpty        = errors.New("order city id is required")
	ErrDistrictIDEmpty    = errors.New("oreder district id is required")
	ErrProductIDEmpty     = errors.New("order product id is required")
	ErrVariantIDEmpty     = errors.New("order variant id is required")
	ErrCityNameEmpty      = errors.New("order city name is required")
	ErrDistrictNameEmpty  = errors.New("order district name is required")
	ErrProductNameEmpty   = errors.New("order product name is required")
	ErrVariantNameEmpty   = errors.New("order variant name is required")
	ErrOrderIDEmpty       = errors.New("order id is required")
	ErrOrderStatusInvalid = errors.New("order status is invalid")
)

// Service creates persisted orders through one repository.
type Service struct {
	repo Repository
}

// New creates order application service.
func New(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create validates confirmed input and stores it as a new order record.
func (s *Service) Create(ctx context.Context, params CreateOrderParams) (CreateResult, error) {
	if s == nil || s.repo == nil {
		return CreateResult{}, ErrRepositoryNil
	}

	record, err := buildRecord(params)
	if err != nil {
		return CreateResult{}, err
	}

	result, err := s.repo.Create(ctx, record)
	if err != nil {
		return CreateResult{}, fmt.Errorf("create order record: %w", err)
	}

	return result, nil
}

// ByID returns one order by identifier.
func (s *Service) ByID(ctx context.Context, id int64) (Order, error) {
	if s == nil || s.repo == nil {
		return Order{}, ErrRepositoryNil
	}
	if id <= 0 {
		return Order{}, ErrOrderIDEmpty
	}

	order, err := s.repo.ByID(ctx, id)
	if err != nil {
		return Order{}, fmt.Errorf("get order by id: %w", err)
	}

	return order, nil
}

// UpdateStatus validates status transition input and stores new order status.
func (s *Service) UpdateStatus(ctx context.Context, id int64, status OrderStatus) error {
	if s == nil || s.repo == nil {
		return ErrRepositoryNil
	}
	if id <= 0 {
		return ErrOrderIDEmpty
	}

	status = normalizeStatus(status)
	if !isValidStatus(status) {
		return ErrOrderStatusInvalid
	}

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("update order status: %w", err)
	}

	return nil
}

// buildRecord validates input and converts it into storage-ready record.
func buildRecord(params CreateOrderParams) (OrderRecord, error) {
	record := OrderRecord{
		BotID:        strings.TrimSpace(params.BotID),
		BotName:      strings.TrimSpace(params.BotName),
		ChatID:       params.ChatID,
		UserID:       params.UserID,
		UserName:     strings.TrimSpace(params.UserName),
		UserUsername: strings.TrimSpace(params.UserUsername),
		CityID:       strings.TrimSpace(params.CityID),
		CityName:     strings.TrimSpace(params.CityName),
		DistrictID:   strings.TrimSpace(params.DistrictID),
		DistrictName: strings.TrimSpace(params.DistrictName),
		ProductID:    strings.TrimSpace(params.ProductID),
		ProductName:  strings.TrimSpace(params.ProductName),
		VariantID:    strings.TrimSpace(params.VariantID),
		VariantName:  strings.TrimSpace(params.VariantName),
		PriceText:    strings.TrimSpace(params.PriceText),
		Status:       OrderStatusNew,
	}

	switch {
	case record.BotID == "":
		return OrderRecord{}, ErrBotIDEmpty
	case record.ChatID == 0:
		return OrderRecord{}, ErrChatIDEmpty
	case record.UserID == 0:
		return OrderRecord{}, ErrUserIDEmpty
	case record.CityID == "":
		return OrderRecord{}, ErrCityIDEmpty
	case record.DistrictID == "":
		return OrderRecord{}, ErrDistrictIDEmpty
	case record.ProductID == "":
		return OrderRecord{}, ErrProductIDEmpty
	case record.VariantID == "":
		return OrderRecord{}, ErrVariantIDEmpty
	case record.CityName == "":
		return OrderRecord{}, ErrCityNameEmpty
	case record.DistrictName == "":
		return OrderRecord{}, ErrDistrictNameEmpty
	case record.ProductName == "":
		return OrderRecord{}, ErrProductNameEmpty
	case record.VariantName == "":
		return OrderRecord{}, ErrVariantNameEmpty
	}

	return record, nil
}

func normalizeStatus(status OrderStatus) OrderStatus {
	return OrderStatus(strings.TrimSpace(string(status)))
}

func isValidStatus(status OrderStatus) bool {
	switch status {
	case OrderStatusNew, OrderStatusInProgress, OrderStatusClosed:
		return true
	default:
		return false
	}
}
