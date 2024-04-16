package tool

import (
	"context"
	"sort"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	All(ctx context.Context, userID int64) ([]*model.Tool, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, userID int64) ([]*model.Tool, error) {
	tools, err := s.repository.All(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(tools) > 1 {
		sort.SliceStable(tools, func(i, j int) bool {
			return tools[i].UpdatedAt.UTC().After(tools[j].UpdatedAt.UTC())
		})
	}

	return tools, nil
}
