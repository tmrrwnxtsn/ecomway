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

func (s *Service) Create(ctx context.Context, data model.CreatePayoutData) (model.CreatePayoutResult, error) {
	var result model.CreatePayoutResult

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
			fmt.Sprintf("cannot create payout for removed tool with id %v", tool.ID),
		)
	}

	op := &model.Operation{
		UserID:           data.UserID,
		Type:             model.OperationTypePayout,
		Currency:         data.Currency,
		Amount:           data.Amount,
		Status:           model.OperationStatusNew,
		ExternalSystem:   data.ExternalSystem,
		ExternalMethod:   data.ExternalMethod,
		ToolID:           data.ToolID,
		Additional:       data.AdditionalData,
		ConfirmationCode: s.codeManager.GenerateCode(),
	}

	if err = s.operationRepository.Create(ctx, op); err != nil {
		return result, fmt.Errorf("create operation in db: %w", err)
	}

	email, _ := op.Additional["email"].(string)

	if err = s.codeManager.SendCode(ctx, op.ID, email, op.ConfirmationCode, data.LangCode); err != nil {
		if saveErr := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &op.ID},
			func(ctx context.Context, op *model.Operation) error {
				op.Status = model.OperationStatusFailed
				op.FailReason = fmt.Sprintf("sending confirmation code: %v", err)
				return nil
			},
		); saveErr != nil {
			err = multierror.Append(err, saveErr)
		}
		return result, err
	}

	result.OperationID = op.ID
	result.Status = op.Status

	return result, nil
}
