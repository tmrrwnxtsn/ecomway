package operation

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error) {
	return s.repository.All(ctx, criteria)
}
