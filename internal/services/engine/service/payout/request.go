package payout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) RequestPayout(ctx context.Context, data model.CreatePayoutData) error {
	tool, err := s.toolRepository.GetOne(ctx, data.ToolID, data.UserID, data.ExternalMethod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf(
					"tool with id %q, user id %v and external method %q not found",
					data.ToolID, data.UserID, data.ExternalMethod,
				),
			)
			if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
				func(ctx context.Context, op *model.Operation) error {
					op.Status = model.OperationStatusFailed
					op.FailReason = err.Error()
					return nil
				},
			); saveErr != nil {
				err = multierror.Append(err, saveErr)
			}
			return err
		}
		return fmt.Errorf("get tool from db: %w", err)
	}

	if tool.Removed() {
		err = perror.NewInternal().WithCode(
			perror.CodeToolHasBeenRemoved,
		).WithDescription(
			fmt.Sprintf("cannot create payout for removed tool with id %v", tool.ID),
		)
		if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = err.Error()
				return nil
			},
		); saveErr != nil {
			err = multierror.Append(err, saveErr)
		}
		return err
	}

	data.Tool = tool

	result, err := s.integrationClient.CreatePayout(ctx, data)
	if err != nil {
		if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = err.Error()
				return nil
			},
		); saveErr != nil {
			err = multierror.Append(err, saveErr)
		}
		return err
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
	default:
		err = s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
			func(ctx context.Context, op *model.Operation) error {
				op.ExternalID = result.ExternalID
				op.ExternalStatus = result.ExternalStatus
				op.Status = model.OperationStatusPending
				return nil
			},
		)
	}

	return err
}
