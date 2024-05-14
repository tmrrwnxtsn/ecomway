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

type ToolRepository interface {
	Update(ctx context.Context, tool *model.Tool) error
	Create(ctx context.Context, tool *model.Tool) error
	GetOne(ctx context.Context, id string, userID int64, externalMethod string) (*model.Tool, error)
}

type Service struct {
	operationRepository OperationRepository
	integrationClient   IntegrationClient
	toolRepository      ToolRepository
}

func NewService(
	operationRepository OperationRepository,
	integrationClient IntegrationClient,
	toolRepository ToolRepository,
) *Service {
	return &Service{
		operationRepository: operationRepository,
		integrationClient:   integrationClient,
		toolRepository:      toolRepository,
	}
}
