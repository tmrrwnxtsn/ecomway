package payout

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"

	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error) {
	var result model.CreatePayoutResult

	op := &model.Operation{
		UserID:         data.UserID,
		Type:           model.OperationTypePayout,
		Currency:       data.Currency,
		Amount:         data.Amount,
		Status:         model.OperationStatusNew,
		ExternalSystem: data.ExternalSystem,
		ExternalMethod: data.ExternalMethod,
		ToolID:         data.ToolID,
		Additional:     data.AdditionalData,
	}

	tool, err := s.toolRepository.GetOne(ctx, data.ToolID, data.UserID, data.ExternalMethod)
	if err != nil {
		return result, fmt.Errorf("get tool from db: %w", err)
	}
	data.Tool = tool

	if err = s.operationRepository.Create(ctx, op); err != nil {
		return result, fmt.Errorf("create operation in db: %w", err)
	}

	// после создания операции на нашей стороне получаем её идентификатор и передаем его в платежную систему
	data.OperationID = op.ID

	result, err = s.integrationClient.CreatePayout(ctx, data)
	if err != nil {
		if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &op.ID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = err.Error()
				return nil
			},
		); saveErr != nil {
			err = multierror.Append(err, saveErr)
		}
		return result, err
	}

	switch result.ExternalStatus {
	case model.OperationExternalStatusSuccess:
		data := model.SuccessPayoutData{
			ProcessedAt:    result.ProcessedAt,
			ExternalID:     result.ExternalID,
			ExternalStatus: result.ExternalStatus,
			OperationID:    data.OperationID,
			Tool:           tool,
		}

		if err = s.Success(ctx, data); err != nil {
			err = fmt.Errorf("failed to success payout: %w", err)
		}

		result.Status = model.OperationStatusSuccess
	case model.OperationExternalStatusFailed:
		data := model.FailPayoutData{
			ExternalID:     result.ExternalID,
			ExternalStatus: result.ExternalStatus,
			FailReason:     result.FailReason,
			OperationID:    data.OperationID,
		}

		if err = s.Fail(ctx, data); err != nil {
			err = fmt.Errorf("failed to fail payout: %w", err)
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