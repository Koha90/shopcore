package service

import (
	"context"
	"fmt"
	"strings"
)

// BotPaymentMethodReader reads payment methods enabled for a bot.
type BotPaymentMethodReader interface {
	ListBotPaymentMethods(ctx context.Context, botID string) ([]BotPaymentMethod, error)
}

// Service contains payment application use cases.
type Service struct {
	methods BotPaymentMethodReader
}

// New constructs payment service.
func New(methods BotPaymentMethodReader) *Service {
	return &Service{
		methods: methods,
	}
}

// ListBotPaymentMethods returns active payment methods configured for a bot.
func (s *Service) ListBotPaymentMethods(ctx context.Context, botID string) ([]BotPaymentMethod, error) {
	botID = strings.TrimSpace(botID)
	if botID == "" {
		return nil, ErrBotIDEmpty
	}
	if s == nil {
		return nil, fmt.Errorf("payment service is nil")
	}
	if s.methods == nil {
		return nil, fmt.Errorf("payment methods reader is nil")
	}

	return s.methods.ListBotPaymentMethods(ctx, botID)
}
