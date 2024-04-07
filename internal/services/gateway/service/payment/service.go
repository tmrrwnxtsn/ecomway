package payment

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) Create(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error) {
	return s.engineClient.CreatePayment(ctx, data)
}
