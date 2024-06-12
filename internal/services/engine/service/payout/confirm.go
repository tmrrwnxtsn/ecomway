package payout

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	perror "github.com/tmrrwnxtsn/ecomway/internal/pkg/error"
	"github.com/tmrrwnxtsn/ecomway/internal/pkg/model"
)

func (s *Service) Confirm(ctx context.Context, data model.ConfirmPayoutData) error {
	if err := s.operationRepository.AcquireOneLocked(ctx, model.OperationCriteria{ID: &data.OperationID},
		func(ctx context.Context, op *model.Operation) error {
			if op.Type != model.OperationTypePayout {
				return fmt.Errorf("payout confirm called for unresolved type: %q", op.Type)
			}

			if op.UserID != data.UserID {
				return fmt.Errorf("invalid user id: %v", data.UserID)
			}

			if op.Status != model.OperationStatusNew {
				if op.Status == model.OperationStatusConfirmed {
					return nil
				}

				return perror.NewInternal().WithCode(
					perror.CodeUnresolvedStatusConflict,
				).WithDescription(
					fmt.Sprintf("cannot confirm payout with status %q", op.Status),
				)
			}

			if data.ConfirmationCode != op.ConfirmationCode {
				op.ConfirmationAttempts++

				if op.ConfirmationAttempts >= s.wrongCodeLimit {
					op.Status = model.OperationStatusFailed
					op.FailReason = model.OperationFailReasonWrongCodeLimitExceeded

					return perror.NewInternal().WithCode(
						perror.CodeConfirmationAttemptsExceeded,
					).WithDescription(
						fmt.Sprintf("code confirm attempts has been exceeded: %v", s.wrongCodeLimit),
					)
				}

				return perror.NewInternal().WithCode(
					perror.CodeWrongConfirmationCode,
				).WithDescription(
					fmt.Sprintf("invalid confirmation code: %q", data.ConfirmationCode),
				)
			}

			op.Status = model.OperationStatusConfirmed

			return nil
		},
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return perror.NewInternal().WithCode(
				perror.CodeObjectNotFound,
			).WithDescription(
				fmt.Sprintf("payout with id %v not found", data.OperationID),
			)
		}
		return err
	}
	return nil
}
