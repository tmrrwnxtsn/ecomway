package operation

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	Operations(ctx context.Context, userID int64) ([]model.Operation, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) Operations(ctx context.Context, userID int64) ([]model.Operation, error) {
	return s.engineClient.Operations(ctx, userID)
}
