package operation

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type Repository interface {
	All(ctx context.Context, criteria model.OperationCriteria) ([]*model.Operation, error)
	AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
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

func (s *Service) AllForReport(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error) {
	return s.repository.AllForReport(ctx, criteria)
}
