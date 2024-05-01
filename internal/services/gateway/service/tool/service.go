package tool

import (
	"context"
	"slices"

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

func (s *Service) AvailableTools(ctx context.Context, userID int64) ([]*model.Tool, error) {
	tools, err := s.engineClient.AvailableTools(ctx, userID)
	if err != nil {
		return nil, err
	}

	tools = slices.DeleteFunc(tools, func(t *model.Tool) bool {
		return t.Fake
	})

	return tools, nil
}

func (s *Service) AvailableToolsGroupedByMethod(ctx context.Context, userID int64) (map[string][]*model.Tool, error) {
	tools, err := s.engineClient.AvailableTools(ctx, userID)
	if err != nil {
		return nil, err
	}

	toolsGroupedByExternalMethod := make(map[string][]*model.Tool)

	for _, tool := range tools {
		if tool.Fake {
			continue
		}

		toolsGroupedByExternalMethod[tool.ExternalMethod] = append(toolsGroupedByExternalMethod[tool.ExternalMethod], tool)
	}

	return toolsGroupedByExternalMethod, nil
}
