package payout

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
	ConfirmPayout(ctx context.Context, data model.ConfirmPayoutData) error
	ResendCode(ctx context.Context, opID, userID int64, langCode string) error
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

func (s *Service) Confirm(ctx context.Context, data model.ConfirmPayoutData) error {
	return s.engineClient.ConfirmPayout(ctx, data)
}

func (s *Service) ResendCode(ctx context.Context, opID, userID int64, langCode string) error {
	return s.engineClient.ResendCode(ctx, opID, userID, langCode)
}
