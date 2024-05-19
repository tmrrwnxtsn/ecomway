package tool

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	AvailableTools(ctx context.Context, userID int64) ([]*model.Tool, error)
	RecoverTool(ctx context.Context, id string, userID int64, externalMethod string) error
	RemoveTool(ctx context.Context, id string, userID int64, externalMethod string) error
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

func (s *Service) RecoverTool(ctx context.Context, id string, userID int64, externalMethod string) error {
	return s.engineClient.RecoverTool(ctx, id, userID, externalMethod)
}

func (s *Service) RemoveTool(ctx context.Context, id string, userID int64, externalMethod string) error {
	return s.engineClient.RemoveTool(ctx, id, userID, externalMethod)
}
