package favorites

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	AddToFavorites(ctx context.Context, data model.FavoritesData) error
	RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) Add(ctx context.Context, data model.FavoritesData) error {
	return s.engineClient.AddToFavorites(ctx, data)
}

func (s *Service) Remove(ctx context.Context, data model.FavoritesData) error {
	return s.engineClient.RemoveFromFavorites(ctx, data)
}
