package operation

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type EngineClient interface {
	ReportOperations(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error)
	GetExternalOperationStatus(ctx context.Context, id int64) (model.OperationExternalStatus, error)
}

type Service struct {
	engineClient EngineClient
}

func NewService(engineClient EngineClient) *Service {
	return &Service{
		engineClient: engineClient,
	}
}

func (s *Service) ReportOperations(ctx context.Context, criteria model.OperationCriteria) ([]model.ReportOperation, error) {
	return s.engineClient.ReportOperations(ctx, criteria)
}

func (s *Service) GetExternalOperationStatus(ctx context.Context, id int64) (model.OperationExternalStatus, error) {
	return s.engineClient.GetExternalOperationStatus(ctx, id)
}
