package payout

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error) {
	return s.engineClient.CreatePayout(ctx, data)
}
