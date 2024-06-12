package favorites

import (
	"context"
	"fmt"
	"slices"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	AddToFavorites(ctx context.Context, data model.FavoritesData) error
	RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error
	GetFavorites(ctx context.Context, userID int64) (model.UserFavorites, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddToFavorites(ctx context.Context, data model.FavoritesData) error {
	return s.repository.AddToFavorites(ctx, data)
}

func (s *Service) RemoveFromFavorites(ctx context.Context, data model.FavoritesData) error {
	return s.repository.RemoveFromFavorites(ctx, data)
}

func (s *Service) FillForMethods(ctx context.Context, opType model.OperationType, userID int64, methods []model.Method) error {
	favorites, err := s.repository.GetFavorites(ctx, userID)
	if err != nil {
		return err
	}

	switch opType {
	case model.OperationTypePayment:
		if len(favorites.Payment) == 0 {
			return nil
		}

		for i, m := range methods {
			if m.IsFavorite {
				continue
			}

			favoriteExternalMethods, found := favorites.Payment[m.ExternalSystem]
			if !found {
				continue
			}

			methods[i].IsFavorite = slices.Contains(favoriteExternalMethods, m.ExternalMethod)
		}
	case model.OperationTypePayout:
		if len(favorites.Payout) == 0 {
			return nil
		}

		for i, m := range methods {
			if m.IsFavorite {
				continue
			}

			favoriteExternalMethods, found := favorites.Payout[m.ExternalSystem]
			if !found {
				continue
			}

			methods[i].IsFavorite = slices.Contains(favoriteExternalMethods, m.ExternalMethod)
		}
	default:
		return fmt.Errorf("unresolved operation type: %q", opType)
	}

	return nil
}
