package payout

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/payout/confirmation"
)

type OperationRepository interface {
	Create(ctx context.Context, op *model.Operation) error
	AcquireOneLocked(ctx context.Context, criteria model.OperationCriteria, script model.ScriptAcquiredFor) (err error)
}

type IntegrationClient interface {
	CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
}

type ToolRepository interface {
	Update(ctx context.Context, tool *model.Tool) error
	GetOne(ctx context.Context, id string, userID int64, externalMethod string) (*model.Tool, error)
}

type ConfirmationCodeManager interface {
	GenerateCode() string
	SendCode(ctx context.Context, operationID int64, email, code, langCode string) error
}

type Service struct {
	operationRepository OperationRepository
	integrationClient   IntegrationClient
	toolRepository      ToolRepository
	codeManager         ConfirmationCodeManager
}

func NewService(
	operationRepository OperationRepository,
	integrationClient IntegrationClient,
	toolRepository ToolRepository,
	isTest bool,
) *Service {
	return &Service{
		operationRepository: operationRepository,
		integrationClient:   integrationClient,
		toolRepository:      toolRepository,
		codeManager:         confirmation.NewCodeManager(isTest),
	}
}
