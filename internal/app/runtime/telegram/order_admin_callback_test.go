package telegram

import (
	"context"

	ordersvc "github.com/koha90/shopcore/internal/order/service"
)

type orderRuntimeServiceStub struct {
	order        ordersvc.Order
	updateID     int64
	updateStatus ordersvc.OrderStatus
	err          error
}

func (s *orderRuntimeServiceStub) Create(ctx context.Context, params ordersvc.CreateOrderParams) (ordersvc.CreateResult, error) {
	return ordersvc.CreateResult{}, s.err
}

func (s *orderRuntimeServiceStub) ByID(ctx context.Context, id int64) (ordersvc.Order, error) {
	return s.order, s.err
}

func (s *orderRuntimeServiceStub) UpdateStatus(ctx context.Context, id int64, status ordersvc.OrderStatus) error {
	s.updateID = id
	s.updateStatus = status
	return s.err
}
