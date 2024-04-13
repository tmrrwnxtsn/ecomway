package payment

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

type OperationRepository interface {
	Create(ctx context.Context, op *model.Operation) error
	AcquireOneLocked(ctx context.Context, criteria model.OperationCriteria, script model.ScriptAcquiredFor) (err error)
}

type IntegrationClient interface {
	CreatePayment(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error)
}

type Service struct {
	operationRepository OperationRepository
	integrationClient   IntegrationClient
}

func NewService(operationRepository OperationRepository, integrationClient IntegrationClient) *Service {
	return &Service{
		operationRepository: operationRepository,
		integrationClient:   integrationClient,
	}
}
