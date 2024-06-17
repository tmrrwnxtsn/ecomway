package tool

import (
	"context"
	"slices"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	AvailableTools(ctx context.Context, userID string) ([]*model.Tool, error)
	EditTool(ctx context.Context, id string, userID string, externalMethod, name string) (*model.Tool, error)
	RemoveTool(ctx context.Context, id string, userID string, externalMethod string) error
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) AvailableTools(ctx context.Context, userID string) ([]*model.Tool, error) {
	tools, err := s.engineClient.AvailableTools(ctx, userID)
	if err != nil {
		return nil, err
	}

	tools = slices.DeleteFunc(tools, func(t *model.Tool) bool {
		return t.Fake || t.Removed()
	})

	return tools, nil
}

func (s *Service) AvailableToolsGroupedByMethod(ctx context.Context, userID string) (map[string][]*model.Tool, error) {
	tools, err := s.AvailableTools(ctx, userID)
	if err != nil {
		return nil, err
	}

	toolsGrouped := make(map[string][]*model.Tool, len(tools))
	for _, tool := range tools {
		toolsGrouped[tool.ExternalMethod] = append(toolsGrouped[tool.ExternalMethod], tool)
	}
	return toolsGrouped, nil
}

func (s *Service) EditTool(ctx context.Context, id string, userID string, externalMethod, name string) (*model.Tool, error) {
	return s.engineClient.EditTool(ctx, id, userID, externalMethod, name)
}

func (s *Service) RemoveTool(ctx context.Context, id string, userID string, externalMethod string) error {
	return s.engineClient.RemoveTool(ctx, id, userID, externalMethod)
}
