package payout

import (
	"context"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
	"github.com/tmrrwnxtsn/ecomway/internal/services/engine/service/payout/confirmation"
)

type OperationRepository interface {
	Create(ctx context.Context, op *model.Operation) error
	AcquireOneLocked(ctx context.Context, criteria model.OperationCriteria, script model.ScriptAcquiredFor) error
}

type IntegrationClient interface {
	CreatePayout(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error)
}

type ToolRepository interface {
	Update(ctx context.Context, tool *model.Tool) error
	GetOne(ctx context.Context, id string, userID string, externalMethod string) (*model.Tool, error)
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
	wrongCodeLimit      int
}

func NewService(
	operationRepository OperationRepository,
	integrationClient IntegrationClient,
	toolRepository ToolRepository,
	smtpClient confirmation.SMTPClient,
	wrongCodeLimit int,
	isTest bool,
) *Service {
	return &Service{
		operationRepository: operationRepository,
		integrationClient:   integrationClient,
		toolRepository:      toolRepository,
		codeManager:         confirmation.NewCodeManager(smtpClient, isTest),
		wrongCodeLimit:      wrongCodeLimit,
	}
}
