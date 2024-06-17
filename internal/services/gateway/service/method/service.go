package method

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	AvailableMethods(ctx context.Context, opType model.OperationType, userID string, currency string) ([]model.Method, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) AvailableMethods(ctx context.Context, opType model.OperationType, userID string, currency string) ([]model.Method, error) {
	return s.engineClient.AvailableMethods(ctx, opType, userID, currency)
}
