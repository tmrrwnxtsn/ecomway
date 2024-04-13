package payment

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Create(ctx context.Context, data model.CreatePaymentData) (model.CreatePaymentResult, error) {
	var result model.CreatePaymentResult

	op := &model.Operation{
		UserID:         data.UserID,
		Type:           model.OperationTypePayment,
		Currency:       data.Currency,
		Amount:         data.Amount,
		Status:         model.OperationStatusNew,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		ToolID:         0, // TODO: возможно, здесь стоит передавать инструмент из model.CreatePaymentData
		Additional:     data.AdditionalData,
	}

	err := s.operationRepository.Create(ctx, op)
	if err != nil {
		return result, fmt.Errorf("save operation: %w", err)
	}

	// после создания операции на нашей стороне получаем её идентификатор, и передаем его в платежную систему
	data.OperationID = op.ID

	result, err = s.integrationClient.CreatePayment(ctx, data)
	if err != nil {
		if err := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &op.ID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = err.Error()
				return nil
			},
		); err != nil {
			slog.Error("failed to create payment", "error", err)
		}
		return result, err
	}

	if err = s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &op.ID},
		func(ctx context.Context, op *model.Operation) error {
			op.ExternalID = result.ExternalID
			op.ExternalStatus = result.ExternalStatus
			return nil
		},
	); err != nil {
		return result, err
	}

	// если операция успешно создана, возвращаем идентификатор созданной операции
	result.OperationID = op.ID

	return result, nil
}
