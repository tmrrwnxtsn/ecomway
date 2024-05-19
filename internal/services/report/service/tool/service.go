package tool

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	AvailableTools(ctx context.Context, userID int64) ([]*model.Tool, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) AllTools(ctx context.Context, userID int64) ([]*model.Tool, error) {
	return s.engineClient.AvailableTools(ctx, userID)
}
