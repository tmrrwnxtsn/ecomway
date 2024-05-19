package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
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
	if data.ToolID != "" {
		tool, err := s.toolRepository.GetOne(ctx, data.ToolID, data.UserID, data.ExternalMethod)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return result, perror.NewInternal().WithCode(
					perror.CodeObjectNotFound,
				).WithDescription(
					fmt.Sprintf(
						"tool with id %q, user id %v and external method %q not found",
						data.ToolID, data.UserID, data.ExternalMethod,
					),
				)
			}
			return result, fmt.Errorf("get tool from db: %w", err)
		}

		if tool.Removed() {
			return result, perror.NewInternal().WithCode(
				perror.CodeToolHasBeenRemoved,
			).WithDescription(
				fmt.Sprintf("cannot create payment for removed tool with id %v", tool.ID),
			)
		}
		data.Tool = tool
	}

	err := s.operationRepository.Create(ctx, op)
	if err != nil {
		return result, fmt.Errorf("create operation in db: %w", err)
	}

	// после создания операции на нашей стороне получаем её идентификатор и передаем его в платежную систему
	data.OperationID = op.ID

	result, err = s.integrationClient.CreatePayment(ctx, data)
	if err != nil {
		if saveOpErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &op.ID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = err.Error()
				return nil
			},
		); saveOpErr != nil {
			err = multierror.Append(err, saveOpErr)
		}
		return result, err
	}

	switch result.ExternalStatus {
	case model.OperationExternalStatusSuccess:
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
