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
		ToolID:         data.ToolID,
		Additional:     data.AdditionalData,
	}

	// если для операции используется сохраненное платежное средство, передаём его тоже
	if data.ToolID != 0 {
		tool, err := s.toolRepository.GetOne(ctx, data.ToolID, data.UserID, data.ExternalMethod)
		if err != nil {
			return result, fmt.Errorf("get tool from db: %w", err)
		}
		data.Tool = tool
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

	switch result.ExternalStatus {
	case model.OperationExternalStatusSuccess:
		// новые инструменты создаём, а уже созданные - обновляем
		if result.Tool != nil {
			result.Tool.ID = data.ToolID
		}

		data := model.SuccessPaymentData{
			ProcessedAt:    result.ProcessedAt,
			ExternalID:     result.ExternalID,
			ExternalStatus: result.ExternalStatus,
			OperationID:    data.OperationID,
			NewAmount:      result.NewAmount,
			Tool:           result.Tool,
		}

		if err = s.Success(ctx, data); err != nil {
			err = fmt.Errorf("failed to success payment: %w", err)
		}

		result.Status = model.OperationStatusSuccess
	case model.OperationExternalStatusFailed:
		data := model.FailPaymentData{
			ExternalID:     result.ExternalID,
			ExternalStatus: result.ExternalStatus,
			FailReason:     result.FailReason,
			OperationID:    data.OperationID,
		}

		if err = s.Fail(ctx, data); err != nil {
			err = fmt.Errorf("failed to fail payment: %w", err)
		}

		result.Status = model.OperationStatusFailed
	default:
		err = s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
			func(ctx context.Context, op *model.Operation) error {
				op.ExternalID = result.ExternalID
				op.ExternalStatus = result.ExternalStatus
				return nil
			},
		)

		result.Status = model.OperationStatusNew
	}
	if err != nil {
		return result, err
	}

	// если операция успешно создана, возвращаем идентификатор созданной операции
	result.OperationID = op.ID

	return result, nil
}
